package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyPostDimension(t *testing.T) {

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

		Convey("Add a dimension to the filter job", func() {

			filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).
				WithBytes([]byte(ValidPOSTAddDimensionToFilterJobJSON)).
				Expect().Status(http.StatusCreated)

			// Check data has been updated as expected
			filterJob, err := mongo.GetFilterJob(database, collection, "filter_job_id", filterJobID)
			if err != nil {
				log.ErrorC("Unable to retrieve updated document", err, nil)
			}

			// Set these empty objects to nil to be able to compare other fields
			filterJob.Downloads = nil
			filterJob.Events = nil

			So(filterJob, ShouldResemble, expectedTestData.ExpectedFilterJob)
		})
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}

	if err := teardownInstance(); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}
}

// func TestFailureToPostDimension(t *testing.T) {
//
// 	if err := setupInstance(); err != nil {
// 		log.ErrorC("Unable to setup instance", err, nil)
// 		os.Exit(1)
// 	}
//
// 	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)
//
// 	Convey("Given an existing filter", t, func() {
//
// 		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
// 			os.Exit(1)
// 		}
//
// 		Convey("Fail to add a dimension to the filter job", func() {
// 			Convey("When the request body is invalid return status bad request (400)", func() {
//
// 				filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).
// 					WithBytes([]byte(InvalidPOSTAddDimensionToFilterJobJSON)).
// 					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body")
// 			})
//
// 			if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
// 				os.Exit(1)
// 			}
//
// 			Convey("When the filter job has a state of `submitted` return status forbidden (403)", func() {
//
// 				// Add submitted filter job to filters collection
// 				if err := mongo.Setup(database, collection, "_id", filterID, ValidSubmittedFilterJob); err != nil {
// 					os.Exit(1)
// 				}
//
// 				filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).
// 					WithBytes([]byte(ValidPOSTAddDimensionToFilterJobJSON)).
// 					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden, the filter job has been locked as it has been submitted to be processed\n")
//
// 			})
//
// 			if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
// 				os.Exit(1)
// 			}
//
// 			Convey("When filter job does not exist returns status not found (404)", func() {
//
// 				filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).
// 					WithBytes([]byte(ValidPOSTAddDimensionToFilterJobJSON)).
// 					Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
// 			})
// 		})
//
// 		if err := teardownInstance(); err != nil {
// 			log.ErrorC("Unable to teardown instance", err, nil)
// 			os.Exit(1)
// 		}
// 	})
// }
