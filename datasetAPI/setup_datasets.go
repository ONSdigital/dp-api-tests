package datasetAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dataset = &mongo.Doc{
	Database:   database,
	Collection: "datasets",
	Key:        "_id",
}

// setupDataset cleans existing stale test data and creates a new test dataset.
func setupDataset(id string, data bson.M) {

	removeExistingDataset(id)

	dataset.Value = id
	dataset.Update = data

	if err := mongo.Setup(dataset); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
}

// removeExistingDataset clears out any existing dataset.
func removeExistingDataset(id string) {
	dataset.Value = id
	if err := mongo.Teardown(dataset); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}
}

// removeDataset removes the dataset that was created in a test.
func removeDataset(id string) {
	dataset.Value = id
	if err := mongo.Teardown(dataset); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}
