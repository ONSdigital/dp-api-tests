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

	docs = append(docs, datasetDoc, publishedEditionDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published edition of a dataset", t, func() {
		Convey("When a GET request is made to retrieve the edition", func() {
			Convey("Then user succeeds and the response returns the edition current sub document", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, edition).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("edition").Equal(edition)
				response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
				response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
				response.Value("links").Object().Value("versions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions$")
				response.Value("state").Equal("published")
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

func TestFailureToGetDatasetEdition(t *testing.T) {

	datasetID := uuid.NewV4().String()
	unpublishedEditionID := uuid.NewV4().String()
	unpublishedEdition := "2018"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	unpublishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     validAssociatedDatasetData(datasetID),
	}

	unpublishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      unpublishedEditionID,
		Update:     validUnpublishedEditionData(datasetID, unpublishedEditionID, unpublishedEdition),
	}

	Convey("Given a request to get an edition of the dataset", t, func() {
		Convey("When the dataset does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Dataset not found")
			})
		})

		Convey("When an unpublished dataset exists", func() {
			if err := mongo.Setup(unpublishedDataset); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("Then the response returns status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Dataset not found")

			})

			Convey("and the request has a valid auth header", func() {
				Convey("Then the response returns status not found (404)", func() {

					datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).
						WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).
						Body().Contains("Dataset not found")

				})
			})

			if err := mongo.Teardown(unpublishedDataset); err != nil {
				if err != mgo.ErrNotFound {
					os.Exit(1)
				}
			}
		})

		Convey("When a published dataset exists but the edition is unpublished for the same dataset", func() {

			if err := mongo.Setup(publishedDataset, unpublishedEditionDoc); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("Then the response returns status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Edition not found")

			})

			Convey("and the request has a valid auth header", func() {
				Convey("Then the response returns status not found (404)", func() {

					datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, unpublishedEdition).
						WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).
						Body().Contains("Edition not found")

				})
			})

			if err := mongo.Teardown(publishedDataset, unpublishedEditionDoc); err != nil {
				if err != mgo.ErrNotFound {
					os.Exit(1)
				}
			}
		})
	})
}
