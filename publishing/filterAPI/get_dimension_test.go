package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"fmt"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting an existing dimension from the filter blueprint", func() {
			Convey("Then return status ok (200) and the expected response for `age` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Object()
				validateDimensionResponse(*response, filterBlueprintID, "age")
			})

			Convey("Then return status ok (200) and the expected response for `sex` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Object()
				validateDimensionResponse(*response, filterBlueprintID, "sex")
			})

			Convey("Then return status ok (200) and the expected response for `aggregate` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Object()
				validateDimensionResponse(*response, filterBlueprintID, "aggregate")
			})

			Convey("Then return status ok (200) and the expected response for `time` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Object()
				validateDimensionResponse(*response, filterBlueprintID, "time")
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	Convey("Given filter blueprint does not exist", t, func() {
		Convey("When a request is made to get a dimension for filter blueprint", func() {
			Convey("Then the response returns status 400 - bad request", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given a filter blueprint", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a request is made to get a dimension for filter blueprint and the dimension does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/wage", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(dimensionNotFoundResponse)
			})
		})

		if err := mongo.Teardown(filter); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func validateDimensionResponse(responseObject httpexpect.Object, filterBlueprintID, dimensionID string) {

	filterURL := fmt.Sprintf("http://localhost:22100/filters/%s", filterBlueprintID)
	selfURL := filterURL + "/dimensions/" + dimensionID
	optionsURL := selfURL + "/options"

	links := responseObject.Value("links").Object()

	So(responseObject.Value("name").Raw(), ShouldEqual, dimensionID)
	So(links.Value("filter").Object().Value("href").Raw(), ShouldEqual, filterURL)
	So(links.Value("filter").Object().Value("id").Raw(), ShouldEqual, filterBlueprintID)
	So(links.Value("self").Object().Value("href").Raw(), ShouldEqual, selfURL)
	So(links.Value("self").Object().Value("id").Raw(), ShouldEqual, dimensionID)
	So(links.Value("options").Object().Value("href").Raw(), ShouldEqual, optionsURL)

}
