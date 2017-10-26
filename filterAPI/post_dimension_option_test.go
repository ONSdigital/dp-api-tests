package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulPostDimensionOptions(t *testing.T) {

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

		filterAPI.POST("/filters/{filter_job_id}/dimensions/age/options/28", filterJobID).
			Expect().Status(http.StatusCreated)
		filterAPI.POST("/filters/{filter_job_id}/dimensions/sex/options/unknown", filterJobID).
			Expect().Status(http.StatusCreated)
		filterAPI.POST("/filters/{filter_job_id}/dimensions/Goods and services/options/welfare", filterJobID).
			Expect().Status(http.StatusCreated)
		filterAPI.POST("/filters/{filter_job_id}/dimensions/time/options/February 2007", filterJobID).
			Expect().Status(http.StatusCreated)

		filterJob, err := mongo.GetFilterJob(database, collection, "filter_job_id", filterJobID)
		if err != nil {
			log.ErrorC("Unable to retrieve updated document", err, nil)
		}

		// Set downloads empty object to nil to be able to compare other fields
		filterJob.Downloads = nil

		expectedFilterJob := expectedTestData.ExpectedFilterJobUpdated
		expectedFilterJob.InstanceID = instanceID
		expectedFilterJob.FilterID = filterJobID

		So(filterJob, ShouldResemble, expectedFilterJob)
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostDimensionOptions(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterJobID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter job does not exist", t, func() {
		invalidFilterJobID := "12345678"

		Convey("When a post request to add an option to a dimension for that filter job", func() {
			Convey("Then return status bad request (400)", func() {

				filterAPI.POST("/filters/{filter_job_id}/dimensions/age/options/30", invalidFilterJobID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
			})
		})
	})

	Convey("Given a filter job with a state of created exists", t, func() {

		update := GetValidCreatedFilterJob(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a post request to add an option for a dimension that does not exist against that filter job", func() {
			Convey("Then return status not found (404)", func() {

				filterAPI.POST("/filters/{filter_job_id}/dimensions/sex/options/male", filterJobID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a filter job with a state of submitted exists", t, func() {

		update := GetValidSubmittedFilterJob(filterID, instanceID, filterJobID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a post request to add an option for a dimension against that filter job", func() {
			Convey("Then return status forbidden (403)", func() {

				filterAPI.POST("/filters/{filter_job_id}/dimensions/sex/options/male", filterJobID).
					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden, the filter job has been locked as it has been submitted to be processed\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
