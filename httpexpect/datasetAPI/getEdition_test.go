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

func TestSuccessfullyGetDatasetEdition(t *testing.T) {

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

	docs = append(docs, datasetDoc, unpublishedEditionDoc, publishedEditionDoc)

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

	Convey("Get an edition of a dataset", t, func() {
		Convey("When user is authenticated and edition is not published", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, "2018").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal("466")
			response.Value("edition").Equal("2018")
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/2018$")
			response.Value("links").Object().Value("versions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/2018/versions$")
			response.Value("state").Equal("edition-confirmed")
		})

		Convey("When user is unauthenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, edition).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(editionID)
			response.Value("edition").Equal(edition)
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
			response.Value("links").Object().Value("versions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions$")
			response.Value("state").Equal("published")
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

func TestFailureToGetDatasetEdition(t *testing.T) {
	var docs []mongo.Doc

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

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

	docs = append(docs, datasetDoc, unpublishedEditionDoc, publishedEditionDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	Convey("Fail to get an edition of a dataset", t, func() {
		Convey("When dataset does not exist", func() {
			datasetAPI.GET("/datasets/{id}/editions/{edition}", "133", "2018").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusBadRequest)
		})

		if err := mongo.Setup(database, collection, "_id", datasetID, validPublishedDatasetData); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When the edition does not exist against dataset", func() {
			datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, "2018").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusNotFound)
		})

		Convey("When the user is unauthenticated and the edition state is NOT set to `published`", func() {
			if err := mongo.Setup(database, "editions", "_id", "466", validUnpublishedEditionData); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, "2018").
				Expect().Status(http.StatusNotFound)
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}
