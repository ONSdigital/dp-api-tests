package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gedge/mgo"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestSuccessfullyGetListOfDatasetEditions(t *testing.T) {

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	unpublishedEditionID := uuid.NewV4().String()

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	unpublishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      unpublishedEditionID,
		Update:     validUnpublishedEditionData(datasetID, unpublishedEditionID, edition),
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, unpublishedEditionDoc)

	Convey("Given a dataset has an edition that is published and one that is unpublished", t, func() {

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

		Convey("When a user is authenticated", func() {
			Convey("Then the response contains both dataset editions", func() {

				response := datasetAPI.GET("/datasets/{id}/editions", datasetID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Length().Equal(2)
				checkEditionsResponse(datasetID, editionID, "2017", response)

				response.Value("items").Array().Element(1).Object().Value("next").Object().Value("edition").Equal(edition)
				response.Value("items").Array().Element(1).Object().Value("id").Equal(unpublishedEditionID)
				response.Value("items").Array().Element(1).Object().Value("next").Object().Value("state").Equal("edition-confirmed")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetListOfDatasetEditions(t *testing.T) {

	datasetID := uuid.NewV4().String()
	unpublishedEditionID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	unpublishedEdition := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      "466",
		Update:     validUnpublishedEditionData(datasetID, unpublishedEditionID, "2018"),
	}

	Convey("Given the dataset does not exist", t, func() {
		Convey("When a request to get editions for a dataset is made", func() {
			Convey("Then return a status unauthorized (401)", func() {

				datasetAPI.GET("/datasets/{id}/editions", datasetID).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	Convey("Given a dataset exists", t, func() {
		if err := mongo.Setup(dataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("but no editions for the dataset exist", func() {
			Convey("When a request to get editions for a dataset is made", func() {
				Convey("Then return a status not found (404)", func() {

					datasetAPI.GET("/datasets/{id}/editions", datasetID).
						WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).Body().Contains("edition not found")
				})
			})
		})

		Convey("and an unpublished edition exists for dataset", func() {

			if err := mongo.Setup(unpublishedEdition); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When a request to get editions for dataset is made by an unauthenticated user", func() {
				Convey("Then return a status unauthorized (401)", func() {

					datasetAPI.GET("/datasets/{id}/editions", datasetID).
						Expect().Status(http.StatusUnauthorized)
				})
			})
		})

		if err := mongo.Teardown(dataset, unpublishedEdition); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func checkEditionsResponse(datasetID, editionID, edition string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("current").Object().Value("edition").Equal(edition)
	response.Value("items").Array().Element(0).Object().Value("current").Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("items").Array().Element(0).Object().Value("current").Object().Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
	response.Value("items").Array().Element(0).Object().Value("current").Object().Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("items").Array().Element(0).Object().Value("current").Object().Value("links").Object().Value("versions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions$")
	response.Value("items").Array().Element(0).Object().Value("current").Object().Value("state").Equal("published")

	response.Value("items").Array().Element(1).Object().Value("next").Object().Value("edition").Equal(edition)
	response.Value("items").Array().Element(1).Object().Value("next").Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("items").Array().Element(1).Object().Value("next").Object().Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
	response.Value("items").Array().Element(1).Object().Value("next").Object().Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("items").Array().Element(1).Object().Value("next").Object().Value("links").Object().Value("versions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions$")
	response.Value("items").Array().Element(1).Object().Value("next").Object().Value("state").Equal("edition-confirmed")
}
