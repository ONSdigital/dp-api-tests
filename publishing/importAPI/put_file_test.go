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

func TestAddFileToImportJob(t *testing.T) {

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

	// These tests needs to refine when authentication was handled in the code.
	Convey("Given a valid request", t, func() {
		Convey("When add file is called", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusOK)

				job, err := mongo.GetJob(cfg.MongoImportsDB, collection, "id", jobID)
				if err != nil {
					t.Errorf("unable to retrieve job resource [%s] from mongo, error: [%v]", jobID, err)
				}

				files := *job.UploadedFiles
				So(len(files), ShouldEqual, 2)
				So(files[0].AliasName, ShouldEqual, "v4")
				So(files[0].URL, ShouldEqual, "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/CPIGrowth.csv")
				So(files[1].AliasName, ShouldEqual, "v5")
				So(files[1].URL, ShouldEqual, "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/CPIGrowth.csv")
				So(job.UniqueTimestamp, ShouldNotBeEmpty)
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

func TestFailureToAddFileToAnImportJob(t *testing.T) {

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	if err := mongo.Setup(importJob); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request with a job ID that does not exist", t, func() {
		Convey("When add file is called", func() {
			Convey("Then the response returns status not found (404)", func() {

				importAPI.PUT("/jobs/{id}/files", invalidJobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("")
			})
		})
	})

	Convey("Given a request with invalid json", t, func() {
		Convey("When add file is called", func() {
			Convey("Then the response returns status bad request(400)", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("invalid request")
			})
		})
	})

	if err := mongo.Teardown(importJob); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestAddFileToImportJobUnauthorised(t *testing.T) {

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
		Convey("When add file is called", func() {
			Convey("Then the response returns status not found (404)", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).
					WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusUnauthorized).
					Body().Contains("")
			})
		})
	})

	Convey("Given a request with an unauthorised Authorization header", t, func() {
		Convey("When add file is called", func() {
			Convey("Then the response returns status 401 unauthorised", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).
					WithHeader(serviceAuthTokenName, unauthorisedServiceAuthToken).
					WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusUnauthorized).
					Body().Contains("")
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
