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

func TestSuccessfullyGetListOfDimensionOptions(t *testing.T) {

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

	Convey("Given an existing filter blueprint with dimensions and options", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting a list of options for a dimension", func() {
			Convey("Then return a list of options for `age` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()
				validateOptionsResponse(*response, filterBlueprintID, "age", []string{"27"})

			})

			Convey("Then return a list of options for `sex` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()
				validateOptionsResponse(*response, filterBlueprintID, "sex", []string{"male", "female"})
			})

			Convey("Then return a list of options for `goods and services` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()
				validateOptionsResponse(*response, filterBlueprintID, "aggregate", []string{"cpi1dim1T60000", "cpi1dim1S10201", "cpi1dim1S10105"})
			})

			Convey("Then return a list of options for `time` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()
				validateOptionsResponse(*response, filterBlueprintID, "time", []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"})
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetListOfDimensionOptions(t *testing.T) {

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
			Convey("Then return status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given a filter blueprint", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a request to get a dimension option against filter blueprint where the dimension does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/wages/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(dimensionNotFoundResponse)
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func validateOptionsResponse(responseArray httpexpect.Array, filterBlueprintID, dimensionID string, expectedOptions []string) {

	for i, option := range expectedOptions {

		filterURL := cfg.FilterAPIURL + "/filters/" + filterBlueprintID
		dimensionURL := filterURL + "/dimensions/" + dimensionID
		selfURL := dimensionURL + "/options/" + option

		links := responseArray.Element(i).Object().Value("links").Object()

		So(responseArray.Element(i).Object().Value("option").Raw(), ShouldEqual, option)
		So(links.Value("dimension").Object().Value("id").Raw(), ShouldEqual, dimensionID)
		So(links.Value("dimension").Object().Value("href").Raw(), ShouldEqual, dimensionURL)
		So(links.Value("filter").Object().Value("id").Raw(), ShouldEqual, filterBlueprintID)
		So(links.Value("filter").Object().Value("href").Raw(), ShouldEqual, filterURL)
		So(links.Value("self").Object().Value("id").Raw(), ShouldEqual, option)
		So(links.Value("self").Object().Value("href").Raw(), ShouldEqual, selfURL)
	}
}
