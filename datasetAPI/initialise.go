package datasetAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	collection = "datasets"

	downloadServiceAuthToken   = "X-Download-Service-Token"
	downloadServiceAuthTokenID = "QB0108EZ-825D-412C-9B1D-41EF7747F462"

	serviceAuthTokenName         = "Authorization"
	serviceAuthToken             = "939616dc-7599-4ded-9a86-a9c66fbf98e0"
	unauthorisedServiceAuthToken = "0dd023bd-9cc0-4c18-9b4f-e030a1f2b71c"
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

	var docs []*mongo.Doc
	for _, c := range []string{collection, "editions", "instances"} {
		t := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: c,
			Key:        "test_data",
			Value:      "true",
		}

		docs = append(docs, t)
	}

	if err = mongo.Teardown(docs...); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
