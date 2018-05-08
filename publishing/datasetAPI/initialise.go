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

	downloadServiceAuthTokenName = "X-Download-Service-Token"
	downloadServiceAuthToken     = "QB0108EZ-825D-412C-9B1D-41EF7747F462"

	downloadServiceTokenName = "Authorization"
	downloadServiceToken     = "Bearer c60198e9-1864-4b68-ad0b-1e858e5b46a4"

	florenceTokenName = "X-Florence-Token"
	florenceToken     = "85c718c3-9ba4-4f31-99bb-3e4eaabb2cc1"

	unauthorisedAuthToken = "0dd023bd-9cc0-4c18-9b4f-e030a1f2b71c"
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
