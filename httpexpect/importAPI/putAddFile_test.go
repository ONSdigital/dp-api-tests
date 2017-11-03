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

// Add a file into a job, for each file added an alias name needs to be given.
// This name needs to link to the recipe
// 200 - The file was added to the import job
func TestAddFileToImportJob(t *testing.T) {

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

	// These tests needs to refine when authentication was handled in the code.
	Convey("Add file to an import job", t, func() {

		Convey("When the user is authenticated", func() {

			importAPI.PUT("/jobs/{id}/files", jobID).WithHeader(internalToken, internalTokenID).
				WithBytes([]byte(validPUTAddFilesJSON)).Expect().Status(http.StatusOK)

		})

		Convey("When the user is unauthenticated", func() {

			importAPI.PUT("/jobs/{id}/files", jobID).WithBytes([]byte(validPUTAddFilesJSON)).
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

func TestFailureToAddFileToAnImportJob(t *testing.T) {

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
	Convey("Fail to add file to an import job", t, func() {
		Convey("and return status not found", func() {
			Convey("When the job id does not exist", func() {
				importAPI.PUT("/jobs/{id}/files", invalidJobID).WithBytes([]byte(validPUTAddFilesJSON)).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	Convey("Fail to add file to an import job", t, func() {
		Convey("and return status bad request", func() {
			Convey("When the invalid json was sent to API", func() {
				importAPI.PUT("/jobs/{id}/files", jobID).WithHeader(internalToken, internalTokenID).WithBytes([]byte("{")).
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
