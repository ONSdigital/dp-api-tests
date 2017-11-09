package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyUpdateVersion(t *testing.T) {

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData,
	}

	editionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData,
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData,
	}

	instanceTwoDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      "799",
		Update:     validUnpublishedInstanceData,
	}

	docs = append(docs, datasetDoc, editionDoc, instanceOneDoc, instanceTwoDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.SetupMany(d); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	// TODO test name incorrect
	Convey("Fail to update a version for an edition of a dataset", t, func() {

		// TODO Test only checks that a version state is updated, we should test all possible fields that could be updated
		datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/2", datasetID, edition).WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTUpdateVersionJSON)).
			Expect().Status(http.StatusOK)
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

// TODO acceptance tests for failure scenarios
