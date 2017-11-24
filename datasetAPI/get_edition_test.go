package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetDatasetEdition(t *testing.T) {

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	unpublishedEditionID := uuid.NewV4().String()
	unpublishedEdition := "2018"

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	publishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	unpublishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      unpublishedEditionID,
		Update:     validUnpublishedEditionData(datasetID, unpublishedEditionID, unpublishedEdition),
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

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(unpublishedEditionID)
			response.Value("edition").Equal(unpublishedEdition)
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + unpublishedEdition + "$")
			response.Value("links").Object().Value("versions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + unpublishedEdition + "/versions$")
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

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	unpublishedEditionID := uuid.NewV4().String()
	unpublishedEdition := "2018"

	var docs []mongo.Doc

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	publishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	unpublishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      unpublishedEditionID,
		Update:     validUnpublishedEditionData(datasetID, unpublishedEditionID, unpublishedEdition),
	}

	docs = append(docs, datasetDoc, unpublishedEditionDoc, publishedEditionDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	Convey("When the dataset does not exist", t, func() {
		Convey("Given a request to get an edition of the dataset", func() {
			Convey("Then the response returns a bad request (400)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})

	Convey("When a dataset exists", t, func() {
		if err := mongo.Setup(database, collection, "_id", datasetID, validPublishedDatasetData(datasetID)); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("but no editions exist against the dataset", func() {
			Convey("Given a request to get an edition of the dataset", func() {
				Convey("Then the response returns a not found (404)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).WithHeader(internalToken, internalTokenID).
						Expect().Status(http.StatusNotFound)
				})
			})
		})

		Convey("and an unpublished edition exists for dataset", func() {

			if err := mongo.Setup(database, "editions", "_id", unpublishedEditionID, validUnpublishedEditionData(datasetID, unpublishedEditionID, unpublishedEdition)); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("Given an unauthenticated request to get an edition of the dataset", func() {
				Convey("Then the response returns a not found (404)", func() {

					datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).
						Expect().Status(http.StatusNotFound)
				})
			})
		})

		if err := mongo.TeardownMany(d); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}
