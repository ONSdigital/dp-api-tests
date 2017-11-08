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

func TestUpdateImportJobState(t *testing.T) {

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup("imports", "imports", "id", jobID, validCreatedImportJobData); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	// This test fails.
	// Bug raised in Trello.
	// These tests needs to refine when authentication was handled in the code.
	Convey("Given an import job exists", t, func() {
		Convey("When a request to update the jobs state with a specific id and the user is authenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.PUT("/jobs/{id}", jobID).WithHeader(internalToken, internalTokenID).
					WithBytes([]byte(validPUTJobJSON)).Expect().Status(http.StatusOK)
			})
		})

		Convey("When a request to update the jobs state with a specific id and the user is unauthenticated", func() {
			Convey("When the user is unauthenticated", func() {

				importAPI.PUT("/jobs/{id}", jobID).WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusOK)
			})
		})
	})

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToUpdateAnImportJob(t *testing.T) {

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup("imports", "imports", "id", jobID, validCreatedImportJobData); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	// This test fails.
	// Bug raised.
	Convey("Given an import job exists", t, func() {
		Convey("When a request to change job state with job id that does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.PUT("/jobs/{id}", invalidJobID).WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	Convey("Given an import job exists", t, func() {
		Convey("When a request to change job state with job id that does not exist", func() {
			Convey("Then the response returns status bad request (400)", func() {
				importAPI.PUT("/jobs/{id}", jobID).WithHeader(internalToken, internalTokenID).WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}
}
