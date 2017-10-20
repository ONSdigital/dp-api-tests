package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPUTVersion_UpdatesVersion(t *testing.T) {

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

	mongo.TeardownMany(d)
	mongo.SetupMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Update a version for an edition of a dataset", t, func() {

		datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/2", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPUTUpdateVersionJSON)).
			Expect().Status(http.StatusOK)
	})

	mongo.TeardownMany(d)

}
