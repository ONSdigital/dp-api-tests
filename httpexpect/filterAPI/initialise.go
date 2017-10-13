package filterAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.ErrorC("Unable to access configurations", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": *cfg})
}
