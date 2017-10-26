package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetListOfDimensionOptions(t *testing.T) {

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}

	if err := setupInstance(); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter job with dimensions and options", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
			os.Exit(1)
		}

		Convey("When requesting a list of options for a dimension", func() {
			Convey("Then return a list of options for `age` dimension", func() {

				response := filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).
					Expect().Status(http.StatusOK).JSON().Array()

				response.Element(0).Object().Value("option").Equal("27")
				response.Element(0).Object().Value("dimension_option_url").NotNull()
			})

			Convey("Then return a list of options for `sex` dimension", func() {

				response := filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options", filterJobID).
					Expect().Status(http.StatusOK).JSON().Array()

				response.Element(0).Object().Value("option").Equal("male")
				response.Element(0).Object().Value("dimension_option_url").NotNull()

				response.Element(1).Object().Value("option").Equal("female")
				response.Element(1).Object().Value("dimension_option_url").NotNull()
			})

			Convey("Then return a list of options for `goods and services` dimension", func() {

				response := filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options", filterJobID).
					Expect().Status(http.StatusOK).JSON().Array()

				response.Element(0).Object().Value("option").Equal("Education")
				response.Element(0).Object().Value("dimension_option_url").NotNull()

				response.Element(1).Object().Value("option").Equal("health")
				response.Element(1).Object().Value("dimension_option_url").NotNull()

				response.Element(2).Object().Value("option").Equal("communication")
				response.Element(2).Object().Value("dimension_option_url").NotNull()
			})

			Convey("Then return a list of options for `time` dimension", func() {

				response := filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options", filterJobID).
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

	if err := teardownInstance(); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}
}

func TestFailureToGetListOfDimensionOptions(t *testing.T) {

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}

	if err := setupInstance(); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter job does not exist", t, func() {
		Convey("When a request to get a dimension option against filter job", func() {
			Convey("Then return status bad request (400)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
			})
		})
	})

	Convey("Given a filter job", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
			os.Exit(1)
		}

		Convey("When a request to get a dimension option against filter job where the dimension does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/wages/options", filterJobID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})
	})

	if err := teardownInstance(); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}
}