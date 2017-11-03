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

func TestSuccessfullyDeleteRemoveDimensionOptions(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	Convey("Given an existing filter", t, func() {

		update := GetValidFilterJobWithMultipleDimensions(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("Remove an option to a dimension to filter on and Verify options are removed", func() {

			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/sex/options/male", filterJobID).Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/Goods and services/options/communication", filterJobID).Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/time/options/April 1997", filterJobID).Expect().Status(http.StatusOK)

			filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).Expect().Status(http.StatusOK)
			sexDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			goodsAndServicesDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			timeDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			sexDimResponse.Element(0).Object().Value("option").NotEqual("male").Equal("female")

			goodsAndServicesDimResponse.Element(0).Object().Value("option").NotEqual("communication").Equal("Education")
			goodsAndServicesDimResponse.Element(1).Object().Value("option").NotEqual("communication").Equal("health")
			timeDimResponse.Element(0).Object().Value("option").NotEqual("April 1997").Equal("March 1997")
			timeDimResponse.Element(1).Object().Value("option").NotEqual("April 1997").Equal("June 1997")
			timeDimResponse.Element(2).Object().Value("option").NotEqual("April 1997").Equal("September 1997")
			timeDimResponse.Element(3).Object().Value("option").NotEqual("April 1997").Equal("December 1997")
		})
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToDeleteRemoveDimensionOptions(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	Convey("Given filter job does not exist", t, func() {
		Convey("When requesting to delete an option from the filter job", func() {

			Convey("Then the response returns status bad request (400)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/wages/options/27000", filterJobID).
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

		// TODO Reinstate commented out code below once API has been updated
		Convey("When requesting to delete an option from a dimension that does not exist against the filter job", func() {
			Convey("Then the response returns status bad request (400)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/wages/options/27000", filterJobID).
					Expect().Status(http.StatusBadRequest) //.Body().Contains("Bad request - dimension not found")
			})
		})

		Convey("When requesting to delete an option that does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/44", filterJobID).
					Expect().Status(http.StatusNotFound) //.Body().Contains("Dimension option not found")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a submitted filter job", t, func() {

		update := GetValidSubmittedFilterJob(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to delete an option", func() {
			Convey("Then the response returns status forbidden (403)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).
					Expect().Status(http.StatusForbidden).Body().
					Contains("Forbidden, the filter job has been locked as it has been submitted to be processed\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
