package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetListOfDatasetEditions(t *testing.T) {

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData,
	}

	publishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData,
	}

	unpublishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      "466",
		Update:     validUnpublishedEditionData,
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, unpublishedEditionDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	mongo.TeardownMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of editions for a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)
			checkEditionsResponse(response)

			response.Value("items").Array().Element(1).Object().Value("edition").Equal("2018")
			response.Value("items").Array().Element(1).Object().Value("id").Equal("466")
			response.Value("items").Array().Element(1).Object().Value("state").Equal("edition-confirmed")
		})

		Convey("When user is unauthenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions", datasetID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(1)
			checkEditionsResponse(response)
		})
	})
	mongo.TeardownMany(d)
}

func TestFailureToGetListOfDatasetEditions(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	mongo.Teardown(database, collection, "_id", datasetID)

	Convey("Fail to get a list of editions for a dataset", t, func() {
		Convey("When the dataset does not exist", func() {

			datasetAPI.GET("/datasets/{id}/editions", datasetID).
				Expect().Status(http.StatusBadRequest)
		})

		mongo.Teardown(database, "editions", "links.dataset.id", datasetID)

		if err := mongo.Setup(database, collection, "_id", datasetID, validPublishedDatasetData); err != nil {
			os.Exit(1)
		}

		Convey("When there are no editions for a dataset", func() {
			datasetAPI.GET("/datasets/{id}/editions", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusNotFound)
		})

		Convey("When user is unauthenticated and there are no published editions", func() {
			if err := mongo.Setup(database, "editions", "_id", "466", validUnpublishedEditionData); err != nil {
				os.Exit(1)
			}

			datasetAPI.GET("/datasets/{id}/editions", datasetID).
				Expect().Status(http.StatusNotFound)

			mongo.Teardown(database, "editions", "_id", "466")
		})

		mongo.Teardown(database, collection, "_id", datasetID)
	})
}

func checkEditionsResponse(response *httpexpect.Object) {
	response.Value("items").Array().Element(0).Object().Value("edition").Equal("2017")
	response.Value("items").Array().Element(0).Object().Value("id").Equal(editionID)
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("versions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions$")
	response.Value("items").Array().Element(0).Object().Value("state").Equal("published")
}
