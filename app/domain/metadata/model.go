package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/polymorph-metadata/app/config"
	"github.com/polymorph-metadata/app/interface/dlt/ethereum"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const POLYMORPH_IMAGE_URL_V1 string = "https://storage.googleapis.com/polymorphs-v1-test/"
const POLYMORPH_IMAGE_URL_V2 string = "https://storage.googleapis.com/polymorph-images_test/"
const GCLOUD_SOURCE_V1_BUCKET_NAME string = "polymorph-source-images"
const GCLOUD_SOURCE_V2_BUCKET_NAME string = "polymorph-source-imgs_test"
const POLYMORPH_V2_THE_GRAPH_HTTP = "https://api.thegraph.com/subgraphs/name/universexyz/polymorph-v2-rinkeby"

const IFRAME_UPLOADED_BASE_URL string = "https://storage.googleapis.com/iframe-htmls/"
const EXTERNAL_URL string = "https://universe.xyz/polymorphs/"
const GENES_COUNT = 9
const BACKGROUND_GENE_COUNT int = 12
const BASE_GENES_COUNT int = 11
const SHOES_GENES_COUNT int = 25
const PANTS_GENES_COUNT int = 33
const TORSO_GENES_COUNT int = 34
const EYEWEAR_GENES_COUNT int = 13
const HEAD_GENES_COUNT int = 31
const WEAPON_RIGHT_GENES_COUNT int = 32
const WEAPON_LEFT_GENES_COUNT int = 32

type Genome string

type Gene int
type StringAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type IntegerAttribute struct {
	TraitType string `json:"trait_type"`
	Value     int    `json:"value"`
}

type FloatAttribute struct {
	TraitType   string  `json:"trait_type"`
	Value       float64 `json:"value"`
	DisplayType string  `json:"display_type"`
}

func (g Gene) toPath() string {
	if g < 10 {
		return fmt.Sprintf("0%s", strconv.Itoa(int(g)))
	}

	return strconv.Itoa(int(g))
}

func getGeneInt(g string, start, end, count int) int {
	genomeLen := len(g)
	geneStr := g[genomeLen+start : genomeLen+end]
	gene, _ := strconv.Atoi(geneStr)
	return gene % count
}

func getWeaponLeftGene(g string) int {
	return getGeneInt(g, -18, -16, WEAPON_LEFT_GENES_COUNT)
}

func getWeaponLeftGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getWeaponLeftGene(g)
	return StringAttribute{
		TraitType: "Left Hand",
		Value:     configService.WeaponLeft[gene],
	}
}

func getWeaponLeftGenePath(g string) string {
	gene := getWeaponLeftGene(g)
	return Gene(gene).toPath()
}

func getWeaponRightGene(g string) int {
	return getGeneInt(g, -16, -14, WEAPON_RIGHT_GENES_COUNT)
}

func getWeaponRightGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getWeaponRightGene(g)
	return StringAttribute{
		TraitType: "Right Hand",
		Value:     configService.WeaponRight[gene],
	}
}

func getWeaponRightGenePath(g string) string {
	gene := getWeaponRightGene(g)
	return Gene(gene).toPath()
}

func getHeadGene(g string) int {
	return getGeneInt(g, -14, -12, HEAD_GENES_COUNT)
}

func getHeadGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getHeadGene(g)
	return StringAttribute{
		TraitType: "Headwear",
		Value:     configService.Headwear[gene],
	}
}

func getHeadGenePath(g string) string {
	gene := getHeadGene(g)
	return Gene(gene).toPath()
}

func getEyewearGene(g string) int {
	return getGeneInt(g, -12, -10, EYEWEAR_GENES_COUNT)
}

func getEyewearGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getEyewearGene(g)
	return StringAttribute{
		TraitType: "Eyewear",
		Value:     configService.Eyewear[gene],
	}
}

func getEyewearGenePath(g string) string {
	gene := getEyewearGene(g)
	return Gene(gene).toPath()
}

func getShoesGene(g string) int {
	return getGeneInt(g, -10, -8, SHOES_GENES_COUNT)
}

func getShoesGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getShoesGene(g)
	return StringAttribute{
		TraitType: "Footwear",
		Value:     configService.Footwear[gene],
	}
}

func getShoesGenePath(g string) string {
	gene := getShoesGene(g)
	return Gene(gene).toPath()
}

func getTorsoGene(g string) int {
	return getGeneInt(g, -8, -6, TORSO_GENES_COUNT)
}

func getTorsoGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getTorsoGene(g)
	return StringAttribute{
		TraitType: "Torso",
		Value:     configService.Torso[gene],
	}
}

func getTorsoGenePath(g string) string {
	gene := getTorsoGene(g)
	return Gene(gene).toPath()
}

func getPantsGene(g string) int {
	return getGeneInt(g, -6, -4, PANTS_GENES_COUNT)
}

func getPantsGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getPantsGene(g)
	return StringAttribute{
		TraitType: "Pants",
		Value:     configService.Pants[gene],
	}
}

func getPantsGenePath(g string) string {
	gene := getPantsGene(g)
	return Gene(gene).toPath()
}

func getBackgroundGene(g string) int {
	return getGeneInt(g, -4, -2, BACKGROUND_GENE_COUNT)
}

func getBackgroundGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getBackgroundGene(g)
	return StringAttribute{
		TraitType: "Background",
		Value:     configService.Background[gene],
	}
}

func getBackgroundGenePath(g string) string {
	gene := getBackgroundGene(g)
	return Gene(gene).toPath()
}

func getBaseGene(g string) int {
	return getGeneInt(g, -2, 0, BASE_GENES_COUNT)
}

func getBaseGeneAttribute(g string, configService *config.ConfigService) StringAttribute {
	gene := getBaseGene(g)
	return StringAttribute{
		TraitType: "Character",
		Value:     configService.Character[gene],
	}
}

func getBaseGenePath(g string) string {
	gene := getBaseGene(g)
	return Gene(gene).toPath()
}

func (g *Genome) name(configService *config.ConfigService, tokenId string) string {
	gStr := string(*g)
	gene := getBaseGene(gStr)
	return fmt.Sprintf("%v #%v", configService.Character[gene], tokenId)
}

func (g *Genome) description(configService *config.ConfigService, tokenId string) string {
	gStr := string(*g)
	gene := getBaseGene(gStr)
	return fmt.Sprintf("The %v named %v #%v is a citizen of the Polymorph Universe and has a unique genetic code! You can scramble your Polymorph at anytime.", configService.Type[gene], configService.Character[gene], tokenId)
}

func (g *Genome) genes() []string {
	gStr := string(*g)

	res := make([]string, 0, GENES_COUNT)

	res = append(res, getWeaponRightGenePath(gStr))
	res = append(res, getWeaponLeftGenePath(gStr))
	res = append(res, getHeadGenePath(gStr))
	res = append(res, getEyewearGenePath(gStr))
	res = append(res, getTorsoGenePath(gStr))
	res = append(res, getPantsGenePath(gStr))
	res = append(res, getShoesGenePath(gStr))
	res = append(res, getBaseGenePath(gStr))
	res = append(res, getBackgroundGenePath(gStr))

	return res
}

func getRarityScoreAttribute(rarity float64) FloatAttribute {
	return FloatAttribute{
		TraitType:   "Rarity Score",
		DisplayType: "number",
		Value:       math.Round(rarity*100) / 100,
	}
}

func getRankAttribute(rank int) IntegerAttribute {
	return IntegerAttribute{
		TraitType: "Rank",
		Value:     rank,
	}
}

func (g *Genome) attributes(configService *config.ConfigService) []interface{} {
	gStr := string(*g)

	res := []interface{}{}
	res = append(res, getBaseGeneAttribute(gStr, configService))
	res = append(res, getShoesGeneAttribute(gStr, configService))
	res = append(res, getPantsGeneAttribute(gStr, configService))
	res = append(res, getTorsoGeneAttribute(gStr, configService))
	res = append(res, getEyewearGeneAttribute(gStr, configService))
	res = append(res, getHeadGeneAttribute(gStr, configService))
	res = append(res, getWeaponLeftGeneAttribute(gStr, configService))
	res = append(res, getWeaponRightGeneAttribute(gStr, configService))
	res = append(res, getBackgroundGeneAttribute(gStr, configService))
	//res = append(res, getRarityScoreAttribute(rarityResponse.RarityScore))
	//res = append(res, getRankAttribute(rarityResponse.Rank))

	return res
}

func badgeGeneContains(s string, list []string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}
	return false
}

// Counts all morph events for a particular polymorph id but doesn't work for our purposes because 8 morphs for a single trait
// is still single trait scrambled
func singleTraitScrambled(id int, ethClient *ethereumclient.EthereumClient) bool {
	header, err := ethClient.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	hashedTopic := common.HexToHash("0x8c0bdd7bca83c4e0c810cbecf44bc544a9dc0b9f265664e31ce0ce85f07a052b")

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   header.Number,
		Addresses: []common.Address{
			contractAddress,
		},
		Topics: [][]common.Hash{{hashedTopic}},
	}
	logs, err := ethClient.Client.FilterLogs(context.Background(), query)
	if err != nil {
		fmt.Println(err.Error())
	}

	morphCount := 0

	id64 := big.NewInt(int64(id))

	for i := range logs {
		currentMorphedId := logs[i].Topics[1].Big()
		if currentMorphedId.Cmp(id64) == 0 {
			morphCount++
		}
		if morphCount > 1 {
			fmt.Println("MorphCount", morphCount)
			return false
		}
	}
	fmt.Println("MorphCount", morphCount)
	return morphCount == 1
}

func getVirginScrambledBadges(id *big.Int) (isVirgin bool, hasSingleTraitScrambledBadge bool) {

	jsonData := map[string]string{
		"query": `
	           {
	 				tokenMorphedEntities(where:{tokenId:` + id.String() + `, eventType_not: 2}, orderBy: timestamp, orderDirection: asc) {
	   				tokenId
					priceForGenomeChange
	 			}
				}
	       `,
	}
	jsonValue, _ := json.Marshal(jsonData)

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", POLYMORPH_V2_THE_GRAPH_HTTP, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonMap map[string]map[string][]map[string]string
	json.Unmarshal([]byte(body), &jsonMap)

	entities := jsonMap["data"]["tokenMorphedEntities"]

	virginBadge := true
	singleTrait := false

	lenEntities := len(entities)

	if len(entities) > 0 && entities[lenEntities-1]["priceForGenomeChange"] == "10000000000000000" {
		virginBadge = false
		singleTrait = false
		return virginBadge, singleTrait
	} else if len(entities) > 0 {
		virginBadge = false
		singleTrait = true
		return virginBadge, singleTrait
	}

	return virginBadge, singleTrait
}

func assignBadges(id string, polymorphGenesList *[]string, badgesJsonMap *map[string][]string) *[]string {

	iTokenId, _ := strconv.Atoi(id)

	var badges []string
	var hasBadge bool

	isVirgin, singleTrait := getVirginScrambledBadges(big.NewInt(int64(iTokenId)))

	if isVirgin {
		badges = append(badges, "never-scrambled")
	}
	if singleTrait {
		badges = append(badges, "single-trait-scrambled")
	}

	leftHandGene := (*polymorphGenesList)[7]
	rightHandGene := (*polymorphGenesList)[8]

	if leftHandGene != "00" && rightHandGene != "00" && leftHandGene == rightHandGene {
		badges = append(badges, "akimbo")
	}

	for badge, geneReqs := range *badgesJsonMap {
		hasBadge = true
		for i, genes := range geneReqs {
			if genes != "**" {
				requirements := strings.Split(genes, "/")
				geneHasBadgeRequirement := badgeGeneContains((*polymorphGenesList)[i], requirements)
				if !geneHasBadgeRequirement {
					hasBadge = false
					break
				}
			}
		}
		if hasBadge {
			badges = append(badges, badge)
		}
	}
	return &badges
}

func (g *Genome) Metadata(ethClient *ethereumclient.EthereumClient, address string, tokenId string, configService *config.ConfigService, badgesJsonMap *map[string][]string) Metadata {
	var m Metadata
	genes := g.genes()

	revGenes := reverseGenesOrder(genes)

	m.Attributes = g.attributes(configService)
	m.Name = g.name(configService, tokenId)
	m.Description = g.description(configService, tokenId)
	m.ExternalUrl = fmt.Sprintf("%s%s", EXTERNAL_URL, tokenId)
	m.Badges = assignBadges(tokenId, &revGenes, badgesJsonMap)

	b := strings.Builder{}
	t := strings.Builder{}
	d := strings.Builder{}

	b.WriteString(POLYMORPH_IMAGE_URL_V1) // Start with base url
	t.WriteString(POLYMORPH_IMAGE_URL_V2)
	// d.WriteString(POLYMORPH_HTML_URL) // Start with HTML BaseURL

	for _, gene := range genes {
		b.WriteString(gene)
		t.WriteString(gene)
		d.WriteString(gene)
	}

	b.WriteString(".jpg") // Finish with jpg extension
	t.WriteString(".jpg")
	// d.WriteString(".html")

	image2DURL := b.String()
	image3DURL := t.String()
	animationURL := d.String()

	image2DExists := imageExists(image2DURL)
	image3DExists := imageExists(image3DURL)

	if !image2DExists {
		generateAndSaveImage(genes, GCLOUD_SOURCE_V1_BUCKET_NAME, GCLOUD_UPLOAD_BUCKET_NAME)
	}
	if !image3DExists {
		generateAndSaveImage(genes, GCLOUD_SOURCE_V2_BUCKET_NAME, GCLOUD_UPLOAD_3D_BUCKET_NAME)
	}

	m.Image2D = image2DURL
	m.Image3D = image3DURL

	cid := generateAndSaveToIpfs(&animationURL, &image2DURL, &image3DURL, m.Badges)

	m.AnimateUrl = "ipfs://" + cid
	return m
}

type Metadata struct {
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Image2D     string      `json:"image2D"`
	Image3D     string      `json:"image3D"`
	Badges      *[]string   `json:"badges_urls"`
	AnimateUrl  string      `json:"animation_url"`
	Attributes  interface{} `json:"attributes"`
	ExternalUrl string      `json:"external_url"`
}
