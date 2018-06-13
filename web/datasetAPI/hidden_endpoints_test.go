package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gedge/mgo"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
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
					Body().Contains("404 page not found")

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
					Body().Contains("404 page not found")
			})
		})

		// DELETE request to /datasets/{id}
		Convey("When a DELETE request to remove the dataset resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("DELETE request on Dataset resource ", log.Data{"endpoint": "/datasets/{id}", "method": "DELETE"})

				datasetAPI.DELETE("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
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
					Body().Contains("404 page not found")
			})
		})

		// INSTANCE

		// GET request to /instance
		Convey("When a GET request to retrieve an instance resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("GET request on Instance resource", log.Data{"endpoint": "/instance", "method": "GET"})

				datasetAPI.GET("/instance").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
			})
		})

		// GET request to /instances
		Convey("When a GET request is made to retrieve all instance resources in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("GET request on Instance resources", log.Data{"endpoint": "/instances", "method": "GET"})

				datasetAPI.GET("/instances").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
			})
		})

		// POST request to /instances
		Convey("When a POST request to create a new instance resource in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("POST request on Instance resource", log.Data{"endpoint": "/instances", "method": "POST"})

				datasetAPI.POST("/instances").
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTCreateInstanceJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
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
					Body().Contains("404 page not found")
			})
		})

		// INSTANCE DIMENSION

		// GET request to /instances/{instance_id}/dimensions
		Convey("When a GET request to retrieve a list of dimension resources for an instance in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("GET request on Dimension resources for an instance", log.Data{"endpoint": "/instances/{instance_id}/dimensions", "method": "GET"})

				datasetAPI.GET("/instances/{instance_id}/dimensions", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
			})
		})

		// INSTANCE DIMENSION OPTION

		// GET request to /instances/{instance_id}/dimensions/{name}/options
		Convey("When a GET request to retrieve a list of dimension option resources for an instance in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("GET request on Dimension Option resources for an instance", log.Data{"endpoint": "/instances/{instance_id}/dimensions/{name}/options", "method": "GET"})

				datasetAPI.GET("/instances/{instance_id}/dimensions/time/options", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
			})
		})

		// POST request to /instances/{instance_id}/dimensions
		Convey("When a POST request to create a dimension option resource for an instance in web", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("POST request on Dimension Option resource for an instance", log.Data{"endpoint": "/instances/{instance_id}/dimensions", "method": "POST"})

				datasetAPI.POST("/instances/{instance_id}/dimensions", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
			})
		})

		// PUT request to /instances/{instance_id}/dimensions/{dimension}/options/{value}/node_id/{node_id}
		Convey("When a PUT request to update a dimension option resource for an instance in web with a node_id", func() {
			Convey("Then response returns a status not found (404)", func() {

				log.Debug("PUT request on Dimension Option resource for an instance", log.Data{"endpoint": "/instances/{instance_id}/dimensions/{dimension}/options/{value}/node_id/{node_id}", "method": "PUT"})

				datasetAPI.PUT("/instances/{instance_id}/dimensions/age/options/23/node_id/123456789", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("404 page not found")
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

	dimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData("9811", instanceID),
	}

	docs = append(docs, publishedDatasetDoc, publishedEditionDoc, publishedInstanceDoc, dimensionDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		return nil, err
	}

	return docs, nil
}
