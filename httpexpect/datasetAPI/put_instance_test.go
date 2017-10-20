package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPutInstance_UpdatesInstance(t *testing.T) {

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData,
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      "799",
		Update:     validUnpublishedInstanceData,
	}

	docs = append(docs, datasetDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	mongo.TeardownMany(d)
	mongo.SetupMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Update an instance properties", t, func() {

		datasetAPI.PUT("/instances/{instance_id}", "799").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPUTCompletedInstanceJSON)).
			Expect().Status(http.StatusOK)
	})

	mongo.TeardownMany(d)
}
