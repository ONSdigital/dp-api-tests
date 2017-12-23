package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestGetInstanceDimensionOptions_ReturnsAllDimensionOptionsFromAnInstance(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database: cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database: cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database: cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	dimensionOneDoc := &mongo.Doc{
		Database: cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData(instanceID),
	}
	dimensionTwoDoc := &mongo.Doc{
		Database: cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData(instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of all unique time  dimension options from an instance", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/instances/{instance_id}/dimensions/time/options", instanceID).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("dimension").Equal("time")
			response.Value("values").Array().Element(0).Equal("202.45")

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/instances/{instance_id}/dimensions/time/options", instanceID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("dimension").Equal("time")
			response.Value("values").Array().Element(0).Equal("202.45")

		})
	})

	Convey("Get a list of all unique aggregate dimension options from an instance", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/instances/{instance_id}/dimensions/aggregate/options", instanceID).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("dimension").Equal("aggregate")
			response.Value("values").Array().Element(0).Equal("cpi1dimA19")

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/instances/{instance_id}/dimensions/aggregate/options", instanceID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("dimension").Equal("aggregate")
			response.Value("values").Array().Element(0).Equal("cpi1dimA19")

		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

// TODO Reinstate tests once code has been fixed
// func TestFailureToGetInstanceDimensionOptions(t *testing.T) {
//
// 	datasetID := uuid.NewV4().String()
// 	instanceID := uuid.NewV4().String()
//
// 	edition := "2017"
//
// 	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
//
// 	Convey("Fail to get instance document", t, func() {
// 		Convey("and return status bad request", func() {
// 			Convey("when instance document does not exist", func() {
// 				datasetAPI.GET("/instances/{id}/dimensions/time/options", "7990").WithHeader(internalToken, internalTokenID).
// 					Expect().Status(http.StatusBadRequest)
// 			})
// 		})
// 		Convey("and return status request is forbidden", func() {
// 			Convey("when an unauthorised user sends a GET request", func() {
// 				if err := mongo.Setup(database, "instances", "_id", "799", validEditionConfirmedInstanceData(datasetID, edition, instanceID)); err != nil {
// 					log.ErrorC("Was unable to run test", err, nil)
// 					os.Exit(1)
// 				}
//
// 				datasetAPI.GET("/instances/{id}/dimensions/time/options", "789").
// 					Expect().Status(http.StatusForbidden)
// 			})
// 		})
//
// 		Convey("and return status not unauthorised", func() {
// 			Convey("when an invalid token is provided", func() {
// 				datasetAPI.GET("/instances/{id}/dimensions/time/options", "789").WithHeader(internalToken, invalidInternalTokenID).
// 					Expect().Status(http.StatusUnauthorized)
// 			})
// 		})
//
// 		Convey("and return status not found", func() {
// 			Convey("dimension_id does not match any dimensions within the instance", func() {
// 				datasetAPI.GET("/instances/{id}/dimensions/timeeww2342/options", "789").WithHeader(internalToken, internalTokenID).
// 					Expect().Status(http.StatusNotFound)
// 			})
// 		})
// 	})
//
// 	if err := mongo.Teardown(database, "instances", "_id", "799"); err != nil {
// 		if err != mgo.ErrNotFound {
// 			log.ErrorC("Failed to tear down test data", err, nil)
// 			os.Exit(1)
// 		}
// 	}
// }
