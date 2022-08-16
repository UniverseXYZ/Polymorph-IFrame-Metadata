package main

import (
	"os"

	"github.com/polymorph-metadata/app/interface/dlt/ethereum"
	log "github.com/sirupsen/logrus"
)

func connectToNodes() (*ethereumclient.EthereumClient, *ethereumclient.EthereumClient) {

	nodeURL := os.Getenv("NODE_URL")

	nodeURLPolygon := os.Getenv("NODE_URL_POLYGON")

	client, err := ethereumclient.NewEthereumClient(nodeURL)

	if err != nil {
		log.Errorln("Error creating new Ethereum client: ", err)
	}

	log.Infoln("Successfully connected to ethereum client")

	clientPolygon, err := ethereumclient.NewEthereumClient(nodeURLPolygon)

	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Successfully connected to polygon client")

	return client, clientPolygon
}
