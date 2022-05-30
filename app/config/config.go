package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigService struct {
	Character   []string `json:"character"`
	Footwear    []string `json:"footwear"`
	Pants       []string `json:"pants"`
	Torso       []string `json:"torso"`
	Eyewear     []string `json:"eyewear"`
	Headwear    []string `json:"headwear"`
	WeaponRight []string `json:"weaponright"`
	WeaponLeft  []string `json:"weaponleft"`
	Type        []string `json:"type"`
	Background  []string `json:"background"`
}
type ConfigBadges struct {
	Naked           []string `json:"naked"`
	BallTeamAway    []string `json:"ball team away"`
	BallTeamHome    []string `json:"ball team home"`
	AmishFarmer     []string `json:"Amish farmer"`
	Astronaut       []string `json:"Astronaut"`
	Rainbow         []string `json:"Rainbow"`
	Golfer          []string `json:"Golfer"`
	Basketball      []string `json:"basketball"`
	BeerLover       []string `json:"beer lover"`
	Marine          []string `json:"marine"`
	GraySuit        []string `json:"gray suit"`
	BowTie          []string `json:"bow tie"`
	BlackSuit       []string `json:"black suit"`
	PlaidSuit       []string `json:"plaid suit"`
	ClownOutfit     []string `json:"clown outfit"`
	Stoner          []string `json:"stoner"`
	Tennis          []string `json:"tennis"`
	SoccerArgentina []string `json:"soccer Argentina"`
	Soccer          []string `json:"soccer"`
	SoccerBrazil    []string `json:"soccer Brazil"`
	SilverKnight    []string `json:"silver knight"`
	GoldKnight      []string `json:"gold knight"`
	Hazmat          []string `json:"hazmat"`
	Pimp            []string `json:"pimp"`
	SushiChef       []string `json:"sushi chef"`
	Hockey          []string `json:"hockey"`
	Ninja           []string `json:"ninja"`
	Spartan         []string `json:"spartan"`
	Samurai         []string `json:"samurai"`
	Tuxedo          []string `json:"tuxedo"`
	Zombie          []string `json:"zombie"`
	Taekwondoe      []string `json:"taekwondoe"`
}

func NewConfigServices(configPath string, badgesPath string) (*ConfigService, *map[string][]string) {
	jsonFileConf, err := os.Open(configPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFileConf.Close()

	jsonFileBadge, err := os.Open(badgesPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFileBadge.Close()

	byteValueConf, _ := ioutil.ReadAll(jsonFileConf)
	byteValueBadge, _ := ioutil.ReadAll(jsonFileBadge)

	jsonBadgeMap := map[string][]string{}

	json.Unmarshal(byteValueBadge, &jsonBadgeMap)

	var service ConfigService

	// we initialize our Users array

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValueConf, &service)

	return &service, &jsonBadgeMap
}
