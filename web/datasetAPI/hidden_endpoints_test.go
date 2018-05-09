package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

// All test responses should return 404 not found,
// even if a valid auth header has been set
func TestPublishingEndpointsAreHiddenForWeb(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published dataset, edition and version (a.k.a instance)", t, func() {
		edition := "2018"
		version := "2"

		docs, err := setupResources(datasetID, editionID, edition, instanceID)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		// DATASET

		// POST request to /datasets/{id}
		Convey("When a POST request to create a new dataset resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("POST request on Dataset resource ", log.Data{"endpoint": "/datasets/{id}", "method": "POST"})

				datasetAPI.POST("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTCreateDatasetJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")

				// ^^ check message is empty, if not it could mean endpoint
				// is available on the web instance of the dataset API
			})
		})

		// PUT request to /datasets/{id}
		Convey("When a PUT request to update the dataset resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("PUT request on Dataset resource ", log.Data{"endpoint": "/datasets/{id}", "method": "PUT"})

				datasetAPI.PUT("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateDatasetJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")
			})
		})

		// DELETE request to /datasets/{id}
		Convey("When a DELETE request to remove the dataset resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("DELETE request on Dataset resource ", log.Data{"endpoint": "/datasets/{id}", "method": "DELETE"})

				datasetAPI.DELETE("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")
			})
		})

		// INSTANCE

		// POST request to /instances
		Convey("When a POST request to create a new instance resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("POST request on Instance resource", log.Data{"endpoint": "/instances", "method": "POST"})

				datasetAPI.POST("/instances").
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTCreateInstanceJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")
			})
		})

		// PUT request to /instances/{instance_id}
		Convey("When a PUT request to update instance meta data in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("PUT request on Instance resource", log.Data{"endpoint": "/instances/{instance_id}", "method": "PUT"})

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTFullInstanceJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")
			})
		})

		// VERSION

		// PUT request to /datasets/{id}/editions/{edition}/versions/{version}
		Convey("When a PUT request to update meta data against the version resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("PUT request on Version resource", log.Data{"endpoint": "/datasets/{id}/editions/{edition}/versions/{version}", "method": "PUT"})

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionMetaDataJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}
	})
}

func setupResources(datasetID, editionID, edition, instanceID string) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	publishedDatasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	publishedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	docs = append(docs, publishedDatasetDoc, publishedEditionDoc, publishedInstanceDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		return nil, err
	}

	return docs, nil
}