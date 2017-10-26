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

func TestSuccessfullyGetFilterJob(t *testing.T) {

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}

	if err := setupInstance(); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
			os.Exit(1)
		}

		Convey("When requesting to get filter job", func() {
			Convey("Then filter job is returned in the response body", func() {

				response := filterAPI.GET("/filters/{filter_job_id}", filterJobID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimension_list_url").String().Match("(.+)/filters/321/dimensions$")
				response.Value("filter_job_id").Equal(filterJobID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/123/editions/2017/versions/1$")
				response.Value("state").Equal("created")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			os.Exit(1)
		}
	})

	if err := teardownInstance(); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetFilterJob(t *testing.T) {

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}

	if err := setupInstance(); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter job does not exist", t, func() {
		Convey("When requesting to get filter job", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filters/{filter_job_id}", "c387798b1rf-0cb623e-43ddf4-bc5df45-78d2284c45fgh").
					Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
			})
		})
	})

	if err := teardownInstance(); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}
}
