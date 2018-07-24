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

func TestGetDimensions_ReturnsAllDimensionsFromADataset(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      ids.DatasetPublished,
		Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionPublished,
		Update:     ValidPublishedEditionData(ids.DatasetPublished, ids.EditionPublished, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      ids.InstancePublished,
		Update:     validPublishedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
	}

	dimensionOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData("9811", ids.InstancePublished),
	}
	dimensionTwoDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData("9812", ids.InstancePublished),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of all dimensions of a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", ids.DatasetPublished, edition, 1).
				WithHeader(florenceTokenName, florenceToken).
				Expect().Status(http.StatusOK).JSON().Object()

			checkDimensionsResponse(ids.DatasetPublished, edition, ids.InstancePublished, response)
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

// TODO Unskip skipped tests when code has been refactored (and hence fixed)
func TestGetDimensions_Failed(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      ids.DatasetPublished,
		Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionPublished,
		Update:     ValidPublishedEditionData(ids.DatasetPublished, ids.EditionPublished, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      ids.InstancePublished,
		Update:     validPublishedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
	}

	docs = append(docs, datasetDoc, editionDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Fail to get a list of Dimensions for a dataset", t, func() {
		// TODO Remove skip on test once endpoint fixed
		SkipConvey("When user is authenticated and the dataset does not exist", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", "1234", "2018").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("dataset not found")
			})
		})

		// TODO Remove skip on test once endpoint fixed
		SkipConvey("When user is authenticated and the edition does not exist", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", ids.DatasetPublished, "2018").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("edition not found")
			})
		})

		Convey("When user is authenticated and there are no versions", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/3/dimensions", ids.DatasetPublished, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("version not found")
			})
		})

		Convey("When user is unauthenticated and the dataset does not exist", func() {
			Convey("Then return status unauthorized(401)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", "1234", "2018").
					Expect().Status(http.StatusUnauthorized).Body()
			})
		})

		Convey("When user is unauthenticated and the edition does not exist", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", ids.DatasetPublished, "2018").
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When there are no published versions", func() {
			// Create an unpublished instance document
			unpublishedInstance := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      ids.InstanceEditionConfirmed,
				Update:     validEditionConfirmedInstanceData(ids.DatasetPublished, edition, ids.InstanceEditionConfirmed, ids.UniqueTimestamp),
			}

			mongo.Setup(unpublishedInstance)
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/3/dimensions", ids.DatasetPublished, edition).
					Expect().Status(http.StatusUnauthorized)
			})

			mongo.Teardown(unpublishedInstance)
		})
	})

	Convey("Given a valid dataset id, edition and version with no dimensions", t, func() {
		Convey("When authenticated", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", ids.DatasetPublished, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("dimensions not found")
			})
		})

		Convey("When user is unauthenticated", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", ids.DatasetPublished, edition).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func checkDimensionsResponse(datasetID, edition, instanceID string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("id").Equal("aggregate")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/aggregate/options$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("aggregate")
	response.Value("items").Array().Element(0).Object().Value("description").Equal("An aggregate of the data")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("options").Object().Value("id").Equal("time")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/time/options$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")

	response.Value("items").Array().Element(1).Object().Value("dimension").Equal("time")
	response.Value("items").Array().Element(1).Object().Value("description").Equal("The time in which this dataset spans")
}
