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

func TestGetInstanceDimensions_ReturnsAllDimensionsFromAnInstance(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	edition := "2017"

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	dimensionOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData(instanceID),
	}
	dimensionTwoDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData(instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

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

	Convey("Get a list of all dimensions from an instance", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/instances/{instance_id}/dimensions", instanceID).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)

			checkInstanceDimensionsResponse(response)

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/instances/{instance_id}/dimensions", instanceID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)

			checkInstanceDimensionsResponse(response)

		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

// TODO Reinstate tests once endpoint is fixed
// func TestFailureToGetInstanceDimensions(t *testing.T) {
// 	datasetID := uuid.NewV4().String()
// 	instanceID := uuid.NewV4().String()
//
// 	edition := "2017"
//
// 	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
//
// 	Convey("Fail to get instance document", t, func() {
// 		Convey("and return status not found", func() {
// 			Convey("when instance document does not exist", func() {
// 				datasetAPI.GET("/instances/{id}/dimensions", "7990").WithHeader(internalToken, internalTokenID).
// 					Expect().Status(http.StatusNotFound)
// 			})
// 		})
// 		Convey("and return status request is forbidden", func() {
// 			Convey("when an unauthorised user sends a GET request", func() {
// 				if err := mongo.Setup(database, "instances", "_id", "799", validEditionConfirmedInstanceData(datasetID, edition, instanceID)); err != nil {
// 					log.ErrorC("Was unable to run test", err, nil)
// 					os.Exit(1)
// 				}
//
// 				datasetAPI.GET("/instances/{id}/dimensions", "789").
// 					Expect().Status(http.StatusForbidden)
// 			})
// 		})
//
// 		Convey("and return status not unauthorised", func() {
// 			Convey("when an invalid token is provided", func() {
// 				datasetAPI.GET("/instances/{id}/dimensions", "789").WithHeader(internalToken, invalidInternalTokenID).
// 					Expect().Status(http.StatusUnauthorized)
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

func checkInstanceDimensionsResponse(response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("dimension_id").Equal("time")

	response.Value("items").Array().Element(0).Object().Value("label").Equal("")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("202.45")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Empty()

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("option").Equal("202.45")

	response.Value("items").Array().Element(0).Object().Value("node_id").Equal("")

	response.Value("items").Array().Element(1).Object().Value("dimension_id").Equal("aggregate")

	response.Value("items").Array().Element(1).Object().Value("label").Equal("CPI (Overall Index)")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code").Object().Value("id").Equal("cpi1dimA19")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/cpi1dimA19$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("version").Object().Empty()

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(1).Object().Value("option").Equal("cpi1dimA19")

	response.Value("items").Array().Element(1).Object().Value("node_id").Equal("")

}
