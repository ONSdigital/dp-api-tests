package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
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

	Convey("Get a list of all versions of a dataset", t, func() {
		Convey("When user is authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)
			checkVersionResponse(response.Value("items").Array().Element(1).Object())

			response.Value("items").Array().Element(0).Object().Value("id").Equal("799")
			response.Value("items").Array().Element(0).Object().Value("state").Equal("associated")
		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(1)
			checkVersionResponse(response.Value("items").Array().Element(0).Object())
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
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
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", "1234", "2018").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, "2018").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest)
			})

			Convey("When there are no versions", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader(internalToken, internalTokenID).
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

func checkVersionResponse(response *httpexpect.Object) {
	response.Value("id").Equal("789")
	response.Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
	response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
	response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("(.+)/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
	response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
	response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
	response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-1-csv$")
	response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
	response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-1-xls$")
	response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
	response.Value("edition").Equal("2017")
	response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
	response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions$")
	response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
	response.Value("links").Object().Value("edition").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
	response.Value("release_date").Equal("2017-12-12") // TODO Should be isodate
	response.Value("spatial").Equal("http://ons.gov.uk/geographylist")
	response.Value("state").Equal("published")
	response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
	response.Value("version").Equal(1)
}
