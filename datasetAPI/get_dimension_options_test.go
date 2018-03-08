package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestGetDimensionOptions_ReturnsAllDimensionOptionsFromADataset(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	dimensionOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData(instanceID),
	}

	dimensionTwoDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData(instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of time dimension options of a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/time/options", datasetID, edition, 1).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			checkTimeDimensionResponse(instanceID, response)

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/time/options", datasetID, edition, 1).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(1)

			checkTimeDimensionResponse(instanceID, response)

		})
	})

	Convey("Get a list of aggregate dimension options of a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/aggregate/options", datasetID, edition, 1).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()
			response.Value("items").Array().Length().Equal(1)

			checkAggregateDimensionResponse(instanceID, response)

		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/aggregate/options", datasetID, edition, 1).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(1)

			checkAggregateDimensionResponse(instanceID, response)

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
// 2 tests skipped
func TestGetDimensionOptions_Failed(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	dimensionOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsDataWithOutOptions(instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, instanceOneDoc, dimensionOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given Fail to get a list of time dimension options for a dataset", t, func() {
		SkipConvey("When user is authenticated and the dataset does not exist", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", "1234", edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dataset not found\n")
			})
		})

		SkipConvey("When user is authenticated and the edition does not exist", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", datasetID, "2018").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound).Body().Contains("Edition not found\n")
			})
		})

		Convey("When user is authenticated and there are no versions", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/5/dimensions/time/options", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")
			})
		})

		SkipConvey("When user is unauthenticated and the dataset does not exist", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", "1234", edition).
					Expect().Status(http.StatusNotFound).Body().Contains("Dataset not found\n")
			})
		})

		SkipConvey("When user is unauthenticated and the edition does not exist", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", datasetID, "2018").
					Expect().Status(http.StatusNotFound).Body().Contains("Edition not found\n")
			})
		})

		Convey("When user is unauthenticated and there are no published versions", func() {
			Convey("Then return status not found (404)", func() {
				// Create an unpublished instance document
				unpublishedInstance := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: "instances",
					Key:        "_id",
					Value:      "799",
					Update:     validEditionConfirmedInstanceData(datasetID, edition, instanceID),
				}

				if err := mongo.Setup(unpublishedInstance); err != nil {
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/5/dimensions/time/options", datasetID, edition).
					Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")

				if err := mongo.Teardown(unpublishedInstance); err != nil {
					if err != mgo.ErrNotFound {
						log.ErrorC("Failed to tear down test data", err, nil)
						os.Exit(1)
					}
				}
			})
		})
	})

	SkipConvey("Given a valid dataset id, edition and version with no dimensions", t, func() {
		Convey("When user is authenticated", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found\n")
			})
		})

		Convey("When user is unauthenticated", func() {
			Convey("Then return status not found (404)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", datasetID, edition).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found\n")
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

func checkTimeDimensionResponse(instanceID string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("time")

	response.Value("items").Array().Element(0).Object().Value("label").Equal("")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("202.45")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("option").Equal("202.45")

}

func checkAggregateDimensionResponse(instanceID string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("aggregate")

	response.Value("items").Array().Element(0).Object().Value("label").Equal("CPI (Overall Index)")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("cpi1dimA19")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/cpi1dimA19$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("option").Equal("cpi1dimA19")

}
