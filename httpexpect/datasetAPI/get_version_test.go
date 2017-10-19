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

func TestSuccessfullyGetVersionOfADatasetEdition(t *testing.T) {
	d, err := teardownVersion()
	if err != nil {
		log.ErrorC("Failed to tear down test data", err, nil)
		os.Exit(1)
	}

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get an existing version for an edition of a dataset", t, func() {
		Convey("When user is authenticated and version is not published", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal("799")
			response.Value("collection_id").Equal("208064B3-A808-449B-9041-EA3A2F72CFAB")
			response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-2-csv$")
			response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
			response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-2-xls$")
			response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
			response.Value("edition").Equal(edition)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2/dimensions$")
			response.Value("links").Object().Value("edition").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
			response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
			response.Value("release_date").Equal("2017-12-12")
			response.Value("state").Equal("associated")
			response.Value("version").Equal(2)
		})

		Convey("When user is unauthenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(instanceID)
			response.Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
			response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-1-csv$")
			response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
			response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-1-xls$")
			response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
			response.Value("edition").Equal(edition)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions$")
			response.Value("links").Object().Value("edition").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
			response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
			response.Value("release_date").Equal("2017-12-12")
			response.Value("state").Equal("published")
			response.Value("version").Equal(1)
		})
	})

	mongo.TeardownMany(d)
}

func TestFailureToGetVersionOfADatasetEdition(t *testing.T) {
	d, err := teardownVersion()
	if err != nil {
		log.ErrorC("Failed to tear down test data", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Fail to get version document", t, func() {
		Convey("and return status bad request", func() {
			Convey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When the edition does not exist", func() {
				mongo.Setup(database, collection, "_id", datasetID, validPublishedDatasetData)
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusBadRequest)
			})
		})
		Convey("and return status not found", func() {
			mongo.Setup(database, "editions", "_id", editionID, validPublishedEditionData)

			Convey("When the version does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusNotFound)
			})
			Convey("When user is unauthenticated and version is not published", func() {
				mongo.Setup(database, "instances", "_id", "799", validUnpublishedInstanceData)

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	mongo.TeardownMany(d)
}

func teardownVersion() (*mongo.ManyDocs, error) {
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

	publishedVersionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData,
	}

	unpublishedVersionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      "799",
		Update:     validUnpublishedInstanceData,
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, publishedVersionDoc, unpublishedVersionDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			return nil, err
		}
	}

	return d, nil
}
