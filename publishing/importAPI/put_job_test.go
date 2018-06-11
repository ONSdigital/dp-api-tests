package importAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestSuccessfullyUpdateImportJobState(t *testing.T) {

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "id",
		Value:      instanceID,
		Update:     validCreatedInstanceData,
	}

	if err := mongo.Setup(importJob, instance); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a valid authenticated request", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.PUT("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusOK)

				job, err := mongo.GetJob(cfg.MongoImportsDB, collection, "id", jobID)
				if err != nil {
					t.Errorf("unable to retrieve job resource [%s] from mongo, error: [%v]", jobID, err)
				}

				So(job.State, ShouldEqual, "submitted")
				So(job.UniqueTimestamp, ShouldNotBeEmpty)

				files := *job.UploadedFiles
				So(len(files), ShouldEqual, 1)
				So(files[0].AliasName, ShouldEqual, "v4")
				So(files[0].URL, ShouldEqual, "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv")
			})
		})
	})

	if err := mongo.Teardown(importJob, instance); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToUpdateImportJobState(t *testing.T) {

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "id",
		Value:      instanceID,
		Update:     validCreatedInstanceData,
	}

	if err := mongo.Setup(importJob, instance); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request with no Authorization header", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status unauthorized (401)", func() {
				importAPI.PUT("/jobs/{id}", jobID).
					WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusUnauthorized).
					Body().Contains("")
			})
		})
	})

	Convey("Given a request with an unauthorised Authorization header", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status unauthorized (401)", func() {

				importAPI.PUT("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, unauthorisedServiceAuthToken).
					WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusUnauthorized).
					Body().Contains("")
			})
		})
	})

	Convey("Given a request for a job that does not exist", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.PUT("/jobs/{id}", invalidJobID).
					WithBytes([]byte(validPUTJobJSON)).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("job not found")
			})
		})
	})

	Convey("Given a request with an invalid job JSON body", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status bad request (400)", func() {
				importAPI.PUT("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("failed to parse json body")
			})
		})
	})

	if err := mongo.Teardown(importJob, instance); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}
