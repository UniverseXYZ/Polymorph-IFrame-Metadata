package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
	"github.com/polymorph-metadata/app/config"
	"github.com/polymorph-metadata/app/interface/api/handlers"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		godotenv.Load(args[0])
	} else {
		godotenv.Load()
	}

	setupLogger()

	ethClient, polygonClient := connectToNodes()

	port := os.Getenv("API_PORT")
	if envport := os.Getenv("PORT"); envport == "" {
		port = os.Getenv("API_PORT")
	}

	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	contractAddressPolygon := os.Getenv("CONTRACT_ADDRESS_POLYGON")

	configService, badgesJsonMap := config.NewConfigServices("./config.json", "./badges-config.json")

	funcframework.RegisterHTTPFunction("/token", handlers.HandleMetadataRequest(ethClient, polygonClient, contractAddress, contractAddressPolygon, configService, badgesJsonMap))

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}

}
