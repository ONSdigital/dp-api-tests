package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestSuccessfullyGetDatasetEdition(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"
	unpublishedEdition := "2018"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      ids.DatasetPublished,
		Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionPublished,
		Update:     ValidPublishedEditionData(ids.DatasetPublished, ids.EditionPublished, edition),
	}

	unpublishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionUnpublished,
		Update:     ValidUnpublishedEditionData(ids.DatasetPublished, ids.EditionUnpublished, unpublishedEdition),
	}

	docs = append(docs, datasetDoc, unpublishedEditionDoc, publishedEditionDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get an edition of a dataset", t, func() {
		Convey("When user is authenticated and edition is not published", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}", ids.DatasetPublished, unpublishedEdition).
				WithHeader(florenceTokenName, florenceToken).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(ids.EditionUnpublished)
			response.Value("next").Object().Value("edition").Equal(unpublishedEdition)
			response.Value("next").Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(ids.DatasetPublished)
			response.Value("next").Object().Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "$")
			response.Value("next").Object().Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + unpublishedEdition + "$")
			response.Value("next").Object().Value("links").Object().Value("versions").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + unpublishedEdition + "/versions$")
			response.Value("next").Object().Value("state").Equal("edition-confirmed")
		})

		Convey("When user is authenticated and edition is published", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}", ids.DatasetPublished, edition).
				WithHeader(florenceTokenName, florenceToken).
				Expect().Status(http.StatusOK).JSON().Object().Value("current").Object()

			response.Value("edition").Equal(edition)
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(ids.DatasetPublished)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + edition + "$")
			response.Value("links").Object().Value("versions").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + edition + "/versions$")
			response.Value("state").Equal("published")
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

// func TestFailureToGetDatasetEdition(t *testing.T) {
// 	ids, err := helpers.GetIDsAndTimestamps()
// 	if err != nil {
// 		log.ErrorC("unable to generate mongo timestamp", err, nil)
// 		t.FailNow()
// 	}
//
// 	unpublishedEdition := "2018"
//
// 	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
//
// 	dataset := &mongo.Doc{
// 		Database:   cfg.MongoDB,
// 		Collection: collection,
// 		Key:        "_id",
// 		Value:      ids.DatasetPublished,
// 		Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
// 	}
//
// 	unpublishedEditionDoc := &mongo.Doc{
// 		Database:   cfg.MongoDB,
// 		Collection: "editions",
// 		Key:        "_id",
// 		Value:      ids.EditionUnpublished,
// 		Update:     validUnpublishedEditionData(ids.DatasetPublished, ids.EditionUnpublished, unpublishedEdition),
// 	}
//
// 	Convey("When the dataset does not exist", t, func() {
// 		Convey("Given a request to get an edition of the dataset", func() {
// 			Convey("Then the response returns status not found (404)", func() {
//
// 				datasetAPI.GET("/datasets/{id}/editions/{edition}", ids.DatasetPublished, ids.EditionUnpublished).
// 					WithHeader(florenceTokenName, florenceToken).
// 					Expect().Status(http.StatusNotFound).Body().Contains("dataset not found")
// 			})
// 		})
// 	})
//
// 	Convey("When a dataset exists", t, func() {
// 		if err := mongo.Setup(dataset); err != nil {
// 			log.ErrorC("Was unable to run test", err, nil)
// 			os.Exit(1)
// 		}
//
// 		Convey("but no editions exist against the dataset", func() {
// 			Convey("Given a request to get an edition of the dataset", func() {
// 				Convey("Then the response returns status not found (404)", func() {
//
// 					datasetAPI.GET("/datasets/{id}/editions/{edition}", ids.DatasetPublished, ids.EditionUnpublished).
// 						WithHeader(florenceTokenName, florenceToken).
// 						Expect().Status(http.StatusNotFound).Body().Contains("edition not found")
// 				})
// 			})
// 		})
//
// 		Convey("and an unpublished edition exists for dataset", func() {
//
// 			if err := mongo.Setup(unpublishedEditionDoc); err != nil {
// 				log.ErrorC("Was unable to run test", err, nil)
// 				os.Exit(1)
// 			}
//
// 			Convey("Given an unauthenticated request to get an edition of the dataset", func() {
// 				Convey("Then the response returns status unauthorized (401)", func() {
//
// 					datasetAPI.GET("/datasets/{id}/editions/{edition}", ids.DatasetPublished, ids.EditionUnpublished).
// 						Expect().Status(http.StatusUnauthorized)
// 				})
// 			})
// 		})
//
// 		if err := mongo.Teardown(dataset, unpublishedEditionDoc); err != nil {
// 			if err != mgo.ErrNotFound {
// 				os.Exit(1)
// 			}
// 		}
// 	})
// }
