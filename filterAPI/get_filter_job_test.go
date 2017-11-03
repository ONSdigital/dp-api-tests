package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetFilterJob(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		update := GetValidFilterJobWithMultipleDimensions(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get filter job", func() {
			Convey("Then filter job is returned in the response body", func() {

				response := filterAPI.GET("/filters/{filter_job_id}", filterJobID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimension_list_url").String().Match("(.+)/filters/" + filterJobID + "/dimensions$")
				response.Value("filter_job_id").Equal(filterJobID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/123/editions/2017/versions/1$")
				response.Value("state").Equal("created")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToGetFilterJob(t *testing.T) {

	filterID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter job does not exist", t, func() {
		Convey("When requesting to get filter job", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filters/{filter_job_id}", filterID).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
			})
		})
	})
}
