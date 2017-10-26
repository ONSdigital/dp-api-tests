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

func TestSuccessfullyGetDimensionOption(t *testing.T) {

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

		Convey("When checking the dimension options", func() {
			Convey("Then return status no content (204) for dimension `age` options", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).
					Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for dimension `sex` options", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options/male", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options/female", filterJobID).
					Expect().Status(http.StatusNoContent)

			})

			Convey("Then return status no content (204) for dimension `goods and services` options", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options/Education", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options/health", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options/communication", filterJobID).
					Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for dimension `time` options", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/March 1997", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/April 1997", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/June 1997", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/September 1997", filterJobID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/December 1997", filterJobID).
					Expect().Status(http.StatusNoContent)
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

func TestFailureToGetDimensionOption(t *testing.T) {

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
			Convey("Then return a status bad request (400)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).
					Expect().Status(http.StatusBadRequest).Body().Contains("filter or dimension not found")
			})
		})
	})

	Convey("Given a filter job containing dimension options", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
			os.Exit(1)
		}

		Convey("When a request to get a dimension option where the dimension does not exist", func() {
			Convey("Then return a status bad request (400)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/ages/options/27", filterJobID).
					Expect().Status(http.StatusBadRequest).Body().Contains("filter or dimension not found")
			})
		})

		Convey("When a request to get a dimension option that does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options/unknown", filterJobID).
					Expect().Status(http.StatusNotFound).Body().Contains("Option not found")
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
