package filterAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	database   = "filters"
	collection = "filters"

	internalTokenHeader    = "Internal-Token"
	internalTokenID        = "FD0108EA-825D-411C-9B1D-41EF7727F465"
	invalidInternalTokenID = "FD0108EA-825D-411C-9B1D-41EF7727F465A"
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

	if err = mongo.Teardown(database, "filterOutputs", "test_data", "true"); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	if err = mongo.Teardown("datasets", "instances", "test_data", "true"); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	if err = mongo.Teardown("datasets", "dimension.options", "test_data", "true"); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
