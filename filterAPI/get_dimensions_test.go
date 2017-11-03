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

func TestSuccessfullyGetListOfDimensions(t *testing.T) {

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

		Convey("When requesting a list of dimensions in a filter job", func() {
			Convey("Then return a list of all dimensions for filter job", func() {

				actual := filterAPI.GET("/filters/{filter_job_id}/dimensions", filterJobID).
					Expect().Status(http.StatusOK).JSON().Array()

				actual.Element(0).Object().Value("dimension_url").NotNull()
				actual.Element(0).Object().Value("name").Equal("age")
				actual.Element(1).Object().Value("dimension_url").NotNull()
				actual.Element(1).Object().Value("name").Equal("sex")
				actual.Element(2).Object().Value("dimension_url").NotNull()
				actual.Element(2).Object().Value("name").Equal("Goods and services")
				actual.Element(3).Object().Value("dimension_url").NotNull()
				actual.Element(3).Object().Value("name").Equal("time")
			})
		})
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetListOfDimensions(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter job does not exist", t, func() {
		Convey("When requesting a list of dimensions for filter job", func() {
			Convey("Then return status not found (404)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions", filterJobID).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
			})
		})
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}
