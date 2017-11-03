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

// Update the state of the job. If this is set to submitted, this shall trigger the import process.
// 200 - The job is in a queue
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
	Convey("Update the state of the import job", t, func() {

		Convey("When the user is authenticated", func() {

			importAPI.PUT("/jobs/{id}", jobID).WithHeader(internalToken, internalTokenID).
				WithBytes([]byte(validPUTJobJSON)).Expect().Status(http.StatusOK)

		})

		Convey("When the user is unauthenticated", func() {

			importAPI.PUT("/jobs/{id}", jobID).WithBytes([]byte(validPUTJobJSON)).
				Expect().Status(http.StatusOK)

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
	Convey("Fail to update an import job", t, func() {
		Convey("and return status not found", func() {
			Convey("When the job id does not exist", func() {
				importAPI.PUT("/jobs/{id}", invalidJobID).WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusNotFound)
			})

		})

	})

	Convey("Fail to update an import job", t, func() {
		Convey("and return status bad request", func() {
			Convey("When the invalid json was sent to API", func() {
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
