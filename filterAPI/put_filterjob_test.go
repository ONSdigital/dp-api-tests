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

func TestSuccessfulPutFilterJob(t *testing.T) {

	if err := setupInstance(); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter with a state of created", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
			os.Exit(1)
		}

		Convey("When filter job is updated with new properties and a change of state to submitted", func() {

			filterAPI.PUT("/filters/{filter_job_id}", filterJobID).
				WithBytes([]byte(ValidPUTUpdateFilterJobJSON)).
				Expect().Status(http.StatusOK)

			Convey("Then filter job state is updated and new dimension options are added", func() {

				filterJob, err := mongo.GetFilterJob(database, collection, "filter_job_id", filterJobID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterJob.State, ShouldNotEqual, "created") // It could be submitted or completed
				So(filterJob.FilterID, ShouldEqual, filterJobID)
				So(len(filterJob.Dimensions), ShouldEqual, 1)
				So(filterJob.Dimensions[0].Name, ShouldEqual, "sex")
				So(filterJob.Dimensions[0].Options, ShouldResemble, []string{"intersex", "other"})
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

func TestFailureToPutFilterJob(t *testing.T) {

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		os.Exit(1)
	}

	if err := setupInstance(); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter job does not exist", t, func() {
		Convey("When a post request is made to update filter job", func() {
			Convey("Then the request fails and returns status not found (404)", func() {

				filterAPI.PUT("/filters/{filter_job_id}", filterJobID).WithBytes([]byte(ValidPUTUpdateFilterJobJSON)).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
			})
		})
	})

	Convey("Given an existing filter job", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidFilterJobWithMultipleDimensions); err != nil {
			os.Exit(1)
		}

		Convey("When an invalid json body is sent to update filter job", func() {
			Convey("Then fail to update filter job and return status bad request (400)", func() {

				filterAPI.PUT("/filters/{filter_job_id}", filterJobID).WithBytes([]byte(InvalidSyntaxJSON)).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			os.Exit(1)
		}
	})

	Convey("Given an existing filter with submitted state", t, func() {

		if err := mongo.Setup(database, collection, "_id", filterID, ValidSubmittedFilterJob); err != nil {
			os.Exit(1)
		}

		Convey("When attempting to update filter job", func() {
			Convey("Then fail to update filter job and return status forbidden (403)", func() {

				filterAPI.PUT("/filters/{filter_job_id}", filterJobID).WithBytes([]byte(ValidPUTUpdateFilterJobJSON)).
					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden, the filter job has been locked as it has been submitted to be processed\n")
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