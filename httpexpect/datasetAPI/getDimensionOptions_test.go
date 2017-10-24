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

func TestGetDimensionOptions_ReturnsAllOptionsOfAnDimension(t *testing.T) {
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
		Collection: "dimensions",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData,
	}
	// dimensionTwoDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimensions",
	// 	Key:        "_id",
	// 	Value:      "9812",
	// 	Update:     validSexDimensionsData,
	// }

	// dimensionTimeOptionsDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimension.options",
	// 	Key:        "_id",
	// 	Value:      dimensionOptionID,
	// 	Update:     validTimeDimensionsOptionsData,
	// }

	// dimensionSexOptionsDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimension.options",
	// 	Key:        "_id",
	// 	Value:      dimensionOptionID,
	// 	Update:     validSexDimensionsOptionsData,
	// }

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	mongo.TeardownMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of all options of a dimension", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/sex/options", datasetID, edition, 1).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(1)
			response.Value("items").Array().Element(0).Object().Value("dimension_id").Equal("sex")
			response.Value("items").Array().Element(0).Object().Value("label").Equal("male")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("2050.56")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/2050.56$")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

			// Why it is empty?
			// response.Value("items").Array().Element(0).Object().Value("options").Equal("")

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/sex/options", datasetID, edition, 1).
				Expect().Status(http.StatusOK).JSON().Object()
			response.Value("items").Array().Length().Equal(1)
			response.Value("items").Array().Element(0).Object().Value("dimension_id").Equal("sex")
			response.Value("items").Array().Element(0).Object().Value("label").Equal("male")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("2050.56")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/2050.56$")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

			// Why it is empty?
			// response.Value("items").Array().Element(0).Object().Value("options").Equal("")
		})
	})

	mongo.TeardownMany(d)
}

func TestGetDimensionOptions_Failed(t *testing.T) {
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
		Collection: "dimensions",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData,
	}
	// dimensionTwoDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimensions",
	// 	Key:        "_id",
	// 	Value:      "9812",
	// 	Update:     validSexDimensionsData,
	// }

	docs = append(docs, datasetDoc, editionDoc, instanceOneDoc, dimensionOneDoc)

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

	Convey("Fail to get a list of versions for a dataset", t, func() {
		Convey("When authenticated", func() {
			Convey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", "1234", "2018").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, "2018").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When there are no versions", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusNotFound)
			})
		})
		Convey("When unauthenticated", func() {
			Convey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", "1234", "2018").
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, "2018").
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When there are no published versions", func() {
				// Create an unpublished instance document
				mongo.Teardown(database, "instances", "_id", "799")
				mongo.Setup(database, "instances", "_id", "799", validUnpublishedInstanceData)
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
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
