package importAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	database   = "imports"
	collection = "imports"
	jobID      = "42B41AE3-8EA6-4D0F-8526-71D1999B4A7D"
	instanceID = "da814aee-66f5-4020-a260-3b6bc7363170"
)

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.ErrorC("Unable to access configurations", err, nil)
		os.Exit(1)
	}

	if err = mongo.NewDatastore(cfg.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		os.Exit(1)
	}

	if err = mongo.Teardown(database, collection, "test_data", "true"); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
