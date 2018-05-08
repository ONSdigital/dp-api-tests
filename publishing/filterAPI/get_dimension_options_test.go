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

				response.Element(0).Object().Value("option").Equal("27")
				response.Element(0).Object().Value("dimension_option_url").NotNull()
			})

			Convey("Then return a list of options for `sex` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				response.Element(0).Object().Value("option").Equal("male")
				response.Element(0).Object().Value("dimension_option_url").NotNull()

				response.Element(1).Object().Value("option").Equal("female")
				response.Element(1).Object().Value("dimension_option_url").NotNull()
			})

			Convey("Then return a list of options for `goods and services` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				response.Element(0).Object().Value("option").Equal("cpi1dim1T60000")
				response.Element(0).Object().Value("dimension_option_url").NotNull()

				response.Element(1).Object().Value("option").Equal("cpi1dim1S10201")
				response.Element(1).Object().Value("dimension_option_url").NotNull()

				response.Element(2).Object().Value("option").Equal("cpi1dim1S10105")
				response.Element(2).Object().Value("dimension_option_url").NotNull()
			})

			Convey("Then return a list of options for `time` dimension", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				response.Element(0).Object().Value("option").Equal("March 1997")
				response.Element(0).Object().Value("dimension_option_url").NotNull()

				response.Element(1).Object().Value("option").Equal("April 1997")
				response.Element(1).Object().Value("dimension_option_url").NotNull()

				response.Element(2).Object().Value("option").Equal("June 1997")
				response.Element(2).Object().Value("dimension_option_url").NotNull()

				response.Element(3).Object().Value("option").Equal("September 1997")
				response.Element(3).Object().Value("dimension_option_url").NotNull()

				response.Element(4).Object().Value("option").Equal("December 1997")
				response.Element(4).Object().Value("dimension_option_url").NotNull()
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
					Expect().Status(http.StatusNotFound).Body().Contains("Filter blueprint not found")
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
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}
