package handlers

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/polymorph-metadata/app/interface/dlt/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/render"
	"github.com/polymorph-metadata/app/config"
	"github.com/polymorph-metadata/app/contracts"
	"github.com/polymorph-metadata/app/domain/metadata"
	log "github.com/sirupsen/logrus"
)

func HandleMetadataRequest(ethClient *ethereumclient.EthereumClient, address string, configService *config.ConfigService, badgesJsonMap *map[string][]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		instance, err := contracts.NewPolymorph(common.HexToAddress(address), ethClient.Client)
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, err)
			log.Errorln(err)
			return
		}

		tokenId := r.URL.Query().Get("id")

		iTokenId, err := strconv.Atoi(tokenId)
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, err)
			log.Errorln(err)
			return
		}

		genomeInt, err := instance.GeneOf(nil, big.NewInt(int64(iTokenId)))
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, err)
			log.Errorln(err)
			return
		}

		// Disable it as we fetch this info from images function
		// rarityResponse := GetRarityById(iTokenId)

		g := metadata.Genome(genomeInt.String())

		render.JSON(w, r, (&g).Metadata(ethClient, address, tokenId, configService, badgesJsonMap))
	}
}
