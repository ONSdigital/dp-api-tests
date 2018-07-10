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

func TestGetDimensionOptions_ReturnsAllDimensionOptionsFromADataset(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	secondInstanceID := uuid.NewV4().String()

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	publishedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	associatedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      secondInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, secondInstanceID),
	}

	publishedTimeDimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData("9811", instanceID),
	}

	publishedAggregateDimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData("9812", instanceID),
	}

	unpublishedTimeDimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9813",
		Update:     validTimeDimensionsData("9813", secondInstanceID),
	}

	unpublishedAggregateDimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9814",
		Update:     validAggregateDimensionsData("9814", secondInstanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, publishedInstanceDoc, publishedTimeDimensionDoc, publishedAggregateDimensionDoc, associatedInstanceDoc, unpublishedTimeDimensionDoc, unpublishedAggregateDimensionDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a list of time dimension options for a published version", t, func() {
		Convey("When an authenticated request is made to get a list of time dimension options", func() {
			Convey("Then return status OK and response body containing dimension options", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/time/options", datasetID, edition, 1).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Length().Equal(1)

				checkTimeDimensionResponse(datasetID, edition, "1", response)
			})
		})
	})

	Convey("Given a list of time dimension options for a unpublished version", t, func() {
		Convey("When an authenticated request is made to get a list of time dimension options", func() {
			Convey("Then return with status OK and response body containing dimension ", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/time/options", datasetID, edition, 2).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				checkTimeDimensionResponse(datasetID, edition, "2", response)
			})
		})
	})

	Convey("Given a list of aggregate dimension options for a published version", t, func() {
		Convey("When an authenticated request is made to get a list of aggregate dimension options", func() {
			Convey("Then return status OK and response body containing dimension options", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/aggregate/options", datasetID, edition, 1).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Length().Equal(1)

				checkAggregateDimensionResponse(datasetID, edition, "1", response)
			})
		})
	})

	Convey("Given a list of aggregate dimension options for a unpublished version", t, func() {
		Convey("When an authenticated request is made to get a list of aggregate dimension options", func() {
			Convey("Then return with status OK and response body containing dimension ", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions/aggregate/options", datasetID, edition, 2).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Length().Equal(1)

				checkAggregateDimensionResponse(datasetID, edition, "2", response)
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

// TODO Unskip skipped tests when code has been refactored (and hence fixed)
// 4 tests skipped
func TestGetDimensionOptions_Failed(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()

	edition := "2017"

	docs := setUpDimensionOptionTestData(datasetID, editionID, instanceID,
		unpublishedInstanceID, edition)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a list of time dimension options for a dataset that does not exist", t, func() {
		SkipConvey("When an authenticated request is made to get a list of time dimension options", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", "1234", edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("dataset not found")
			})
		})
	})

	Convey("Given a list of time dimension options for an edition that does not exist", t, func() {
		SkipConvey("When an authenticated request is made to get a list of time dimension options", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/time/options", datasetID, "2018").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("edition not found")
			})
		})
	})

	Convey("Given a list of time dimension options for a version that does not exist", t, func() {
		Convey("When an authenticated request is made to get a list of time dimension options", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/5/dimensions/time/options", datasetID, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("version not found")
			})
		})
	})

	Convey("Given a list of time dimension options for an unpublished version", t, func() {
		Convey("When an unauthorised request is made to get a list of time dimension options", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/dimensions/time/options", datasetID, edition).
					Expect().Status(http.StatusUnauthorized)

			})
		})
	})

	SkipConvey("Given aggregate dimension does not exist for a version", t, func() {
		Convey("When a request is made to get a list of aggregate dimension options", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions/aggregate/options", datasetID, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Dimension not found")
			})
		})
	})

	SkipConvey("Given aggregate dimension does not exist for a unpublished version", t, func() {
		Convey("When an authenticated request is made to get a list of aggregate dimension options", func() {
			Convey("Then return status not found (404)", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/dimensions/time/options", datasetID, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Dimension not found")
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

func checkTimeDimensionResponse(datasetID, edition, version string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("time")

	response.Value("items").Array().Element(0).Object().Value("label").Equal("")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("202.45")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/" + version + "$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("option").Equal("202.45")

}

func checkAggregateDimensionResponse(datasetID, edition, version string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("aggregate")

	response.Value("items").Array().Element(0).Object().Value("label").Equal("CPI (Overall Index)")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("cpi1dimA19")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/cpi1dimA19$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/" + version + "$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("option").Equal("cpi1dimA19")

}

func setUpDimensionOptionTestData(datasetID, editionID, instanceID, unpublishedInstanceID, edition string) []*mongo.Doc {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	publishedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	publishedTimeDimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsDataWithOutOptions("9811", instanceID),
	}

	unpublishedInstance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      "799",
		Update:     validEditionConfirmedInstanceData(datasetID, edition, unpublishedInstanceID),
	}

	unpublishedTimeDimensionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validTimeDimensionsDataWithOutOptions("9812", instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, publishedInstanceDoc, publishedTimeDimensionDoc, unpublishedInstance, unpublishedTimeDimensionDoc)

	return docs
}
