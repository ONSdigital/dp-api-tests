package datasetAPI

import (
	"gopkg.in/mgo.v2/bson"
	"os"
	"gopkg.in/mgo.v2"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

// setupDataset cleans existing stale test data and creates a new test dataset.
func setupDataset(datasetID string, datasetData bson.M) {

	removeExistingDataset(datasetID)

	if err := mongo.Setup(database, "datasets", "_id", datasetID, datasetData); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
}

// removeExistingDataset clears out any existing dataset.
func removeExistingDataset(datasetID string) {
	if err := mongo.Teardown(database, collection, "_id", datasetID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}
}

// removeDataset removes the dataset that was created in a test.
func removeDataset(datasetID string) {
	if err := mongo.Teardown(database, "datasets", "_id", datasetID); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}


