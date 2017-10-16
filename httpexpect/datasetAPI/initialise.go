package datasetAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const database = "datasets"

var (
	collection = "datasets"
	datasetID  = "123"
	editionID  = "456"
	edition    = "2017"
	instanceID = "789" // This maybe known as the version Id
)

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.ErrorC("Unable to access configurations", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
