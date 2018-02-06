package searchAPI

import (
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	collection = "datasets"

	instanceID             = "123789"
	internalTokenHeader    = "Internal-Token"
	internalTokenID        = "SD0108EA-825D-411C-45J3-41EF7727F123"
	invalidInternalTokenID = "SD0108EA-825D-411C-45J3-41EF7727F123A"
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

	if err = mongo.Teardown(docs...); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	if err = createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, "aggregate"); err != nil {
		log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}

func createSearchIndex(url, instanceID, dimension string) error {
	log.Info("nath got here", nil)

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	directory := ".."
	if filepath.Base(currentPath) == "dp-api-tests" {
		directory = "."
	}

	index := elasticsearch.Index{
		InstanceID:   instanceID,
		Dimension:    dimension,
		TestDataFile: directory + "/testDataSetup/elasticsearch/testData.json",
		URL:          url,
		MappingsFile: directory + "/testDataSetup/elasticsearch/mappings.json",
	}

	if err := index.CreateSearchIndex(); err != nil {
		return err
	}

	return nil
}
