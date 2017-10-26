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

func TestSuccessfullyDeleteDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter job", t, func() {
		update := GetValidFilterJobWithMultipleDimensions(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When sending a delete request to remove an existing dimension on the filter job", func() {
			Convey("Then the filter job should not contain that dimension", func() {
				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/Goods and services", filterJobID).
					Expect().Status(http.StatusOK)

				var expectedDimensions []mongo.Dimension

				dimensionAge := mongo.Dimension{
					DimensionURL: "",
					Name:         "age",
					Options:      []string{"27"},
				}

				dimensionSex := mongo.Dimension{
					DimensionURL: "",
					Name:         "sex",
					Options:      []string{"male", "female"},
				}

				dimensionTime := mongo.Dimension{
					DimensionURL: "",
					Name:         "time",
					Options:      []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
				}

				expectedDimensions = append(expectedDimensions, dimensionAge, dimensionSex, dimensionTime)

				// Check dimension has been removed from filter job
				filterJob, err := mongo.GetFilterJob(database, collection, "filter_job_id", filterJobID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
					os.Exit(1)
				}

				So(filterJob.Dimensions, ShouldResemble, expectedDimensions)
			})
		})
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToDeleteDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter job does not exist", t, func() {
		Convey("When requesting to delete a dimension from filter job", func() {
			Convey("Then response returns status bad request (400)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age", filterJobID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
			})
		})
	})

	Convey("Given an existing filter job with submitted state", t, func() {

		update := GetValidSubmittedFilterJob(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to delete a dimension from filter job", func() {
			Convey("Then response returns status forbidden (403)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age", filterJobID).
					Expect().Status(http.StatusForbidden).Body().
					Contains("Forbidden, the filter job has been locked as it has been submitted to be processed\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given an existing filter", t, func() {

		update := GetValidFilterJobWithMultipleDimensions(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to delete a dimension from filter job where the dimension does not exist", func() {
			Convey("Then response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_job_id}/dimensions/wage", filterJobID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
