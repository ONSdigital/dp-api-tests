package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersions_ReturnsListOfVersions(t *testing.T) {
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

	instanceTwoDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      "799",
		Update:     validUnpublishedInstanceData,
	}

	docs = append(docs, datasetDoc, editionDoc, instanceTwoDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	mongo.TeardownMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of all versions of a dataset", t, func() {
		Convey("When user is authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)
			checkVersionResponse(response, 1)

			response.Value("items").Array().Element(0).Object().Value("id").Equal("799")
			response.Value("items").Array().Element(0).Object().Value("state").Equal("associated")
		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(1)
			checkVersionResponse(response, 0)
		})
	})

	mongo.TeardownMany(d)
}

func TestGetVersions_Failed(t *testing.T) {
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

	mongo.TeardownMany(d)
	mongo.SetupMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Fail to get a list of versions for a dataset", t, func() {
		Convey("When authenticated", func() {
			// TODO Uncomment tests once code is fixed
			// Convey("When the dataset does not exist", func() {
			// 	datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", "1234", "2018").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
			// 		Expect().Status(http.StatusBadRequest)
			// })

			// Convey("When the edition does not exist", func() {
			// 	datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, "2018").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
			// 		Expect().Status(http.StatusBadRequest)
			// })

			Convey("When there are no versions", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusNotFound)
			})
		})
		Convey("When unauthenticated", func() {
			// TODO Uncomment tests once code is fixed
			// Convey("When the dataset does not exist", func() {
			// 	datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", "1234", "2018").
			// 		Expect().Status(http.StatusBadRequest)
			// })

			// Convey("When the edition does not exist", func() {
			// 	datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, "2018").
			// 		Expect().Status(http.StatusBadRequest)
			// })

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
	mongo.TeardownMany(d)
}

func checkVersionResponse(response *httpexpect.Object, item int) {
	response.Value("items").Array().Element(item).Object().Value("id").Equal("789")
	response.Value("items").Array().Element(item).Object().Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
	response.Value("items").Array().Element(item).Object().Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-1-csv$")
	response.Value("items").Array().Element(item).Object().Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
	response.Value("items").Array().Element(item).Object().Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-1-xls$")
	response.Value("items").Array().Element(item).Object().Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
	response.Value("items").Array().Element(item).Object().Value("edition").Equal("2017")
	response.Value("items").Array().Element(item).Object().Value("license").Equal("ONS License")
	response.Value("items").Array().Element(item).Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("items").Array().Element(item).Object().Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
	response.Value("items").Array().Element(item).Object().Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions$")
	response.Value("items").Array().Element(item).Object().Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
	response.Value("items").Array().Element(item).Object().Value("links").Object().Value("edition").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
	// TODO uncomment out line below whenapi has been fixed
	// response.Value("items").Array().Element(item).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
	response.Value("items").Array().Element(item).Object().Value("release_date").Equal("2017-12-12") // TODO Should be isodate
	response.Value("items").Array().Element(item).Object().Value("state").Equal("published")
	response.Value("items").Array().Element(item).Object().Value("version").Equal(1)
}
