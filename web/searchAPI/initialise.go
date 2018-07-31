package searchAPI

import (
	"os"
	"time"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	collection = "datasets"

	instanceID        = "123789"
	serviceToken      = "Bearer a507f722-f25a-4889-9653-23a2655b925c"
	florenceTokenName = "X-Florence-Token"
	florenceToken     = "85c718c3-9ba4-4f31-99bb-3e4eaabb2cc1"

	skipTeardown          = false
	timeout               = 5 * time.Second
	retryPause            = 750 * time.Millisecond
	dimensionKeyAggregate = "aggregate"
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
	for _, c := range []string{collection, "datasets", "editions", "instances"} {
		t := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: c,
			Key:        "test_data",
			Value:      "true",
		}

		docs = append(docs, t)
	}

	if !skipTeardown {
		if err = mongo.Teardown(docs...); err != nil {
			log.ErrorC("Unable to remove all test data from mongo db", err, nil)
			os.Exit(1)
		}
	}
	log.Debug("config is:", log.Data{"config": cfg})
}
