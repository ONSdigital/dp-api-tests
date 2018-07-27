package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetDimensionOption(t *testing.T) {

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

		Convey("When checking the dimension options", func() {
			Convey("Then return status no content (204) for dimension `age` options", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "age", "27")
			})

			Convey("Then return status no content (204) for dimension `sex` options", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "sex", "male")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/female", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "sex", "female")
			})

			Convey("Then return status no content (204) for dimension `goods and services` options", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1S10201", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "aggregate", "cpi1dim1S10201")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1S10105", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "aggregate", "cpi1dim1S10105")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1T60000", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "aggregate", "cpi1dim1T60000")
			})

			Convey("Then return status no content (204) for dimension `time` options", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/March 1997", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "time", "March 1997")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/April 1997", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "time", "April 1997")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/June 1997", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "time", "June 1997")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/September 1997", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "time", "September 1997")

				response = filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/December 1997", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				validateOptionResponse(*response, filterBlueprintID, "time", "December 1997")
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetDimensionOption(t *testing.T) {

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

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a request to get a dimension option against filter blueprint", func() {
			Convey("Then return a status 400 bad request", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
					Expect().Status(http.StatusBadRequest).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given a filter blueprint containing dimension options", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a request to get a dimension option where the dimension does not exist", func() {
			Convey("Then return a status bad request (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/ages/options/27", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains(dimensionNotFoundResponse)
			})
		})

		Convey("When a request to get a dimension option that does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/unknown", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains(optionNotFoundResponse)
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func validateOptionResponse(responseObject httpexpect.Object, filterBlueprintID, dimensionID, option string) {

	filterURL := cfg.FilterAPIURL + "/filters/" + filterBlueprintID
	dimensionURL := filterURL + "/dimensions/" + dimensionID
	selfURL := dimensionURL + "/options/" + option

	links := responseObject.Value("links").Object()

	So(responseObject.Value("option").Raw(), ShouldEqual, option)
	So(links.Value("dimension").Object().Value("id").Raw(), ShouldEqual, dimensionID)
	So(links.Value("dimension").Object().Value("href").Raw(), ShouldEqual, dimensionURL)
	So(links.Value("filter").Object().Value("id").Raw(), ShouldEqual, filterBlueprintID)
	So(links.Value("filter").Object().Value("href").Raw(), ShouldEqual, filterURL)
	So(links.Value("self").Object().Value("id").Raw(), ShouldEqual, option)
	So(links.Value("self").Object().Value("href").Raw(), ShouldEqual, selfURL)
}
