package metadata

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	"image"
	"image/color"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const IpfsBaseUploadURL = `https://api.pinata.cloud/pinning/pinFileToIPFS`

type pinataBodyResponse struct {
	IPFSHash string `json:"IpfsHash"`
	PinSize  string `json:"PinSize"`
}

const IMG_SIZE = 4000

// New bucket for 3d polymorphs

func imageExists(imageURL string) bool {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Fatalln(err)
	}
	return resp.StatusCode != 404
}
func objectExists(imageURL string) bool {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	iframeHtmlsBucketName := os.Getenv("IFRAME_HTMLS_BUCKET_NAME")

	if err != nil {
		log.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(iframeHtmlsBucketName)
	exists, err := bucket.Object(imageURL).If(storage.Conditions{DoesNotExist: true}).NewReader(ctx)

	return exists != nil

}

//func cidExists(animationURL *string) string {
//	ctx := context.Background()
//	client, err := storage.NewClient(ctx)
//
//	if err != nil {
//		log.Errorf("storage.NewClient: %v", err)
//	}
//	defer client.Close()
//
//	//var builder cid.V0Builder
//
//	f, err := os.Create("iframe-go.html")
//	if err != nil {
//		fmt.Println("os.Create: %v", err)
//	}
//
//	bucketReader, err := client.Bucket(IFRAME_HTMLS_BUCKET_NAME).Object(*animationURL).NewReader(ctx)
//
//	defer bucketReader.Close()
//
//	if _, err := io.Copy(f, bucketReader); err != nil {
//		fmt.Println("io.Copy: %v", err)
//	}
//
//	if err = f.Close(); err != nil {
//		fmt.Errorf("f.Close: %v", err)
//	}
//
//	pref := cid.Prefix{
//		Version:  1,
//		Codec: mc.Raw,
//		MhType:   mh.SHA2_256,
//		MhLength: -1, // default length
//	}
//
//	dat, err := ioutil.ReadFile("iframe-go.html")
//
//	c, err := pref.Sum(dat)
//
//	fmt.Println("Created CID: ", c.Hash().B58String())
//	return c.String()
//
//	//var msg, _ = ioutil.ReadAll(f)
//	//
//	//c, _ := cid.V0Builder.Sum(builder, msg)
//	//
//	//return c.Hash().B58String()
//}

func combineRemoteImages(bucket *storage.BucketHandle, basePath string, overlayPaths ...string) *image.NRGBA {

	ctx := context.Background()

	baseReader, err := bucket.Object(basePath).NewReader(ctx)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer baseReader.Close()
	base, err := imaging.Decode(baseReader)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}
	dst := imaging.New(IMG_SIZE, IMG_SIZE, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, base, image.Pt(0, 0))

	for _, op := range overlayPaths {
		r, err := bucket.Object(op).NewReader(ctx)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
		defer r.Close()
		o, err := imaging.Decode(r)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
		dst = imaging.Overlay(dst, o, image.Pt(0, 0), 1)
	}

	return dst
}

func reverseGenesOrder(genes []string) []string {
	res := make([]string, 0, len(genes))
	for i := len(genes) - 1; i >= 0; i-- {
		res = append(res, genes[i])
	}
	return res
}

func saveToGCloud(i *image.NRGBA, name string, bucketName string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName).Object(name).NewWriter(ctx)
	// f, err := imaging.FormatFromFilename(name)
	// if err != nil {
	// 	log.Errorf("Format from filename: %v", err)
	// }
	err = imaging.Encode(bucket, i, imaging.JPEG, imaging.JPEGQuality(80))

	if err != nil {
		log.Errorf("Upload: %v", err)
	}

	if err = bucket.Close(); err != nil {
		log.Errorf("Writer.Close: %v", err)
	}

}

func generateAndSaveImage(genes []string, gCloudSourceBucket string, gCloudUploadBucket string) {
	// Reverse
	revGenes := reverseGenesOrder(genes)

	f := make([]string, len(genes))

	for i, gene := range revGenes {
		f[i] = fmt.Sprintf("./images/%v/%s.png", i, gene)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(gCloudSourceBucket)

	i := combineRemoteImages(bucket, f[0], f[1:]...)

	b := strings.Builder{}

	for _, gene := range genes {
		b.WriteString(gene)
	}

	b.WriteString(".jpg") // Finish with jpg extension

	saveToGCloud(i, b.String(), gCloudUploadBucket)
}

type ImageURLs struct {
	V1Img string
	V2Img string
}

type Badge struct {
	Name string
	URL  string
}

type TemplateHTML struct {
	ImgUrls ImageURLs
	Badges  []Badge
}

// generateAndSaveToIpfs Generates the polymorph animation url and uploads it to IPFS
func generateAndSaveToIpfs(iframeURL *string, image2DURL *string, image3DURL *string, badges *[]string) (cid string) {

	htmlBadges := make([]Badge, len(*badges))

	baseBadgeUrl := os.Getenv("BADGE_BASE_URL")

	for i, badge := range *badges {
		htmlBadges[i].Name = badge
		htmlBadges[i].URL = fmt.Sprintf("%v%s.svg", baseBadgeUrl, badge)
		(*badges)[i] = fmt.Sprintf("%v%s.svg", baseBadgeUrl, badge)
	}

	gcloudExists := objectExists(*iframeURL)

	tmpl := template.Must(template.ParseFiles("./serverless_function_source_code/index.html"))
	data := TemplateHTML{
		ImgUrls: ImageURLs{*image2DURL, *image3DURL},
		Badges:  htmlBadges,
	}

	if !gcloudExists { // If false, write the animation html to GCloud
		ctx := context.Background()
		client, err := storage.NewClient(ctx)

		iframeHtmlsBucketName := os.Getenv("IFRAME_HTMLS_BUCKET_NAME")

		if err != nil {
			log.Errorf("storage.NewClient: %v", err)
		}
		defer client.Close()

		tpl := &bytes.Buffer{}

		bucket := client.Bucket(iframeHtmlsBucketName).Object(*iframeURL).NewWriter(ctx)

		if err := tmpl.Execute(tpl, data); err != nil {
			fmt.Println(err)
		}

		bucket.Write(tpl.Bytes())
		if err := bucket.Close(); err != nil {
			fmt.Println("createFile: unable to close bucket %q, file %q: %v", err)
		}
	}

	f, err := os.Create("/tmp/iframe-go.html")
	if err != nil {
		log.Println("create file: ", err.Error())
	}
	err = tmpl.Execute(f, data)

	if err != nil {
		log.Print("Error executing html animation template: ", err.Error())
	}
	err = f.Close()
	if err != nil {
		log.Println("Error closing file %v %s. Original error ", f.Name(), err.Error())
	}

	PinataApiKey := os.Getenv("PINATA_API_KEY")
	PinataSecretKey := os.Getenv("PINATA_SECRET_KEY")

	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("/tmp/iframe-go.html")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("file", filepath.Base("/tmp/iframe-go.html"))

	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest(method, IpfsBaseUploadURL, payload)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("pinata_api_key", PinataApiKey)
	req.Header.Add("pinata_secret_api_key", PinataSecretKey)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("Content-Length", string(payload.Len()))

	// Send HTTP Post request to Pinata
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	pinataResponse := pinataBodyResponse{}

	err = decoder.Decode(&pinataResponse)

	httpClient.CloseIdleConnections()

	e := os.Remove("/tmp/iframe-go.html")
	if e != nil {
		log.Fatal(e)
	}

	return pinataResponse.IPFSHash
}
