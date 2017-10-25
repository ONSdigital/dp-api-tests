package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestGetDimensions_ReturnsAllDimensionsFromADataset(t *testing.T) {
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

	dimensionOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData,
	}
	dimensionTwoDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData,
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

	Convey("Get a list of all dimensions of a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", datasetID, edition, 1).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			checkDimensionsResponse(response)

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", datasetID, edition, 1).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)

			checkDimensionsResponse(response)

		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

func TestGetDimensions_Failed(t *testing.T) {
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

	docs = append(docs, datasetDoc, editionDoc)

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

	Convey("Fail to get a list of Dimensions for a dataset", t, func() {
		Convey("When authenticated", func() {
			Convey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", "1234", "2018").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusNotFound)
			})

			Convey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, "2018").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusNotFound)
			})

			Convey("When there are no versions", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusNotFound)
			})
		})
		Convey("When unauthenticated", func() {
			Convey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", "1234", "2018").
					Expect().Status(http.StatusNotFound)
			})

			Convey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, "2018").
					Expect().Status(http.StatusNotFound)
			})

			Convey("When there are no published versions", func() {
				// Create an unpublished instance document
				mongo.Teardown(database, "instances", "_id", "799")
				mongo.Setup(database, "instances", "_id", "799", validUnpublishedInstanceData)
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, edition).
					Expect().Status(http.StatusNotFound)

				mongo.Teardown(database, "instances", "_id", "799")
			})
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

func checkDimensionsResponse(response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("id").Equal("aggregate")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/aggregate/options$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("items").Array().Element(0).Object().Value("dimension_id").Equal("aggregate")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("options").Object().Value("id").Equal("time")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/time/options$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("items").Array().Element(1).Object().Value("dimension_id").Equal("time")
}
