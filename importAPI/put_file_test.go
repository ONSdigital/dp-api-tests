package importAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestAddFileToImportJob(t *testing.T) {

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup("imports", "imports", "id", jobID, validCreatedImportJobData); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	// These tests needs to refine when authentication was handled in the code.
	Convey("Given an import job exists", t, func() {
		Convey("When a request to add a file into a job and the user is authenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).WithHeader(internalToken, internalTokenID).
					WithBytes([]byte(validPUTAddFilesJSON)).Expect().Status(http.StatusOK)
			})
		})

		Convey("When a request to add a file into a job and the user is unauthenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusOK)
			})
		})
	})

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToAddFileToAnImportJob(t *testing.T) {

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup("imports", "imports", "id", jobID, validCreatedImportJobData); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	// This test fails.
	// Bug raised.
	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given an import job exists", t, func() {
		Convey("When a request to add a file into a job with job id that does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.PUT("/jobs/{id}/files", invalidJobID).WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	Convey("Given an import job exists", t, func() {
		Convey("When a request to add a file into a job with invalid json", func() {
			Convey("Then the response returns status bad request(400)", func() {
				importAPI.PUT("/jobs/{id}/files", jobID).WithHeader(internalToken, internalTokenID).WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}
