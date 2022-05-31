package metadata

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"strings"
	"text/template"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
)

const IMG_SIZE = 4000

// New bucket for 3d polymorphs
const GCLOUD_UPLOAD_BUCKET_NAME = "polymorphs-v1-test"
const GCLOUD_UPLOAD_3D_BUCKET_NAME = "polymorph-images_test"
const IFRAME_HTMLS_BUCKET_NAME = "iframe-htmls"
const BADGE_BASE_URL = "https://storage.googleapis.com/iframe-source/img/badge/"

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

	if err != nil {
		log.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(IFRAME_HTMLS_BUCKET_NAME)
	stats, err := bucket.Object(imageURL).Attrs(ctx)

	return stats != nil
}

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

//func uploadHtml() {
//	func uploadHandler(w http.ResponseWriter, r *http.Request) {
//		if r.Method != "POST" {
//			http.Error(w, "", http.StatusMethodNotAllowed)
//			return
//		}
//
//		ctx := appengine.NewContext(r)
//
//		f, fh, err := r.FormFile("file")
//		if err != nil {
//			msg := fmt.Sprintf("Could not get file: %v", err)
//			http.Error(w, msg, http.StatusBadRequest)
//			return
//		}
//		defer f.Close()
//
//		sw := storageClient.Bucket(bucket).Object(fh.Filename).NewWriter(ctx)
//		if _, err := io.Copy(sw, f); err != nil {
//			msg := fmt.Sprintf("Could not write file: %v", err)
//			http.Error(w, msg, http.StatusInternalServerError)
//			return
//		}
//
//		if err := sw.Close(); err != nil {
//			msg := fmt.Sprintf("Could not put file: %v", err)
//			http.Error(w, msg, http.StatusInternalServerError)
//			return
//		}
//
//		u, _ := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
//
//		fmt.Fprintf(w, "Successful! URL: https://storage.googleapis.com%s", u.EscapedPath())
//	}
//
//}

func generateAndSaveIFrame(iframeURL *string, image2DURL *string, image3DURL *string, badges *[]string) {

	htmlBadges := make([]Badge, len(*badges))

	for i, badge := range *badges {
		htmlBadges[i].Name = badge
		htmlBadges[i].URL = fmt.Sprintf("%v%s.svg", BADGE_BASE_URL, badge)
		(*badges)[i] = fmt.Sprintf("%v%s.svg", BADGE_BASE_URL, badge)
	}

	tmpl := template.Must(template.ParseFiles("./serverless_function_source_code/index.html"))
	data := TemplateHTML{
		ImgUrls: ImageURLs{*image2DURL, *image3DURL},
		Badges:  htmlBadges,
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(IFRAME_HTMLS_BUCKET_NAME).Object(*iframeURL).NewWriter(ctx)

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	bucket.Write(tpl.Bytes())
	if err := bucket.Close(); err != nil {
		fmt.Println("createFile: unable to close bucket %q, file %q: %v", err)
		return
	}
}
