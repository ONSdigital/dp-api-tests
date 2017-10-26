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

func TestSuccessfullyGetDimension(t *testing.T) {

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

		Convey("When requesting an existing dimension from the filter job", func() {
			Convey("Then return status no content (204) for `age` dimension", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/age", filterJobID).
					Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for `sex` dimension", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/sex", filterJobID).Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for `goods and services` dimension", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services", filterJobID).Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for `time` dimension", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/time", filterJobID).Expect().Status(http.StatusNoContent)
			})
		})
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter job does not exist", t, func() {
		Convey("When a request is made to get a dimension for filter job", func() {
			Convey("Then the response returns status bad request (400)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/age", filterJobID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
			})
		})
	})

	Convey("Given a filter job", t, func() {

		update := GetValidFilterJobWithMultipleDimensions(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}
		Convey("When a request is made to get a dimension for filter job and the dimension does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/wage", filterJobID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
