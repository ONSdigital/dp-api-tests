package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestGetInstanceDimensionOptions_ReturnsAllDimensionOptionsFromAnInstance(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a list of all unique time dimension options for an instance exists", t, func() {
		docs, err := getInstanceDimensionOptionsSetup(ids.DatasetPublished, ids.EditionPublished, edition, ids.InstancePublished, ids.UniqueTimestamp)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated user sends a GET request for a list of time options for instance", func() {
			Convey("Then a list of time options is returned with a status of OK (200)", func() {

				response := datasetAPI.GET("/instances/{instance_id}/dimensions/time/options", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimension").Equal("time")
				response.Value("options").Array().Element(0).Equal("202.45")
			})
		})

		Convey("When an authenticated user sends a GET request for a list of aggregate options for instance", func() {
			Convey("Then a list of aggregate options is returned with a status of OK (200)", func() {

				response := datasetAPI.GET("/instances/{instance_id}/dimensions/aggregate/options", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimension").Equal("aggregate")
				response.Value("options").Array().Element(0).Equal("cpi1dimA19")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetInstanceDimensionOptions(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance document does not exist", t, func() {
		Convey("When an unauthenticated request to get an instances dimension options", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions/time/options", ids.InstancePublished).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a user sends a GET request for an instances dimension options with an invalid Authentication header", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions/time/options", ids.InstancePublished).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an authenticated user sends a GET request for an instances dimension options", func() {
			Convey("Then return status not found (404) with a message `instance not found`", func() {

				datasetAPI.GET("/instances/{id}/dimensions/time/options", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("instance not found")
			})
		})
	})

	Convey("Given an instance document does exist", t, func() {

		instanceDoc := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstancePublished,
			Update:     validEditionConfirmedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
		}

		if err := mongo.Setup(instanceDoc); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When a user sends a GET request for an instances dimension options without sending a token", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions/time/options", ids.InstancePublished).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a user sends a GET request for an instances dimension options with an invalid token", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions/time/options", ids.InstancePublished).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an authenticated user sends a GET request for an instances dimension options", func() {
			Convey("Then return status not found (404) with a message `dimension node not found`", func() {

				datasetAPI.GET("/instances/{id}/dimensions/time/options", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("dimension node not found\n")
			})
		})

		if err := mongo.Teardown(instanceDoc); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Failed to tear down test data", err, nil)
				os.Exit(1)
			}
		}
	})
}

func getInstanceDimensionOptionsSetup(datasetID, editionID, edition, instanceID string, uniqueTimestamp bson.MongoTimestamp) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	dimensionOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData("9811", instanceID),
	}

	dimensionTwoDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData("9812", instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
