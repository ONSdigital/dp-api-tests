package downloadService

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	collection  = "datasets"
	publicLink  = "https://s3-eu-west-1.amazonaws.com/dp-frontend-florence-file-uploads/2470609-cpicoicoptestcsv"
	privateLink = "s3://csv-exported/v4TestFile.csv"
	region      = "eu-west-1"
	bucketName  = "csv-exported"
	fileName    = "v4TestFile.csv"

	authHeader   = "Authorization"
	serviceToken = "Bearer c60198e9-1864-4b68-ad0b-1e858e5b46a4"

	publishedTrue = true
	publishedFalse = false
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
