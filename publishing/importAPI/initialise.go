package importAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	collection                   = "imports"
	jobID                        = "42B41AE3-8EA6-4D0F-8526-71D1999B4A7D"
	invalidJobID                 = "42B41AE38EA64D0F852671D1999B4A7D1234"
	instanceID                   = "da814aee-66f5-4020-a260-3b6bc7363170"
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

	test := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "test_data",
		Value:      "true",
	}

	if err = mongo.Teardown(test); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
