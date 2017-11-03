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

// Get information about a single job
// 200 - Return a single job information
func TestSuccessfullyGetAnImportJob(t *testing.T) {

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
	Convey("Get information about a single job", t, func() {

		Convey("When the user is authenticated", func() {

			response := importAPI.GET("/jobs/{id}", jobID).WithHeader(internalToken, internalTokenID).Expect().Status(http.StatusOK).JSON().Object()

			checkImportJobResponse(response)

		})

		Convey("When the user is unauthenticated", func() {

			response := importAPI.GET("/jobs/{id}", jobID).Expect().Status(http.StatusOK).JSON().Object()

			checkImportJobResponse(response)

		})

	})

	if err := mongo.Teardown("imports", "imports", "id", jobID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}
}

// 404 - JobId does not match any import jobs
func TestFailureToGetAnImportJob(t *testing.T) {

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

	Convey("Fail to get an import job", t, func() {
		Convey("and return status not found", func() {
			Convey("When the job id does not exist", func() {
				importAPI.GET("/jobs/{id}", invalidJobID).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound)
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

func checkImportJobResponse(response *httpexpect.Object) {

	response.Value("id").Equal(jobID)
	response.Value("recipe").Equal("2080CACA-1A82-411E-AA46-F00804968E78")
	response.Value("state").Equal("Created")

	//Raised bug for this
	response.Value("files").Array().Element(0).Object().Value("alias_name").Equal("v4")

	response.Value("files").Array().Element(0).Object().Value("url").Equal("https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/CPIGrowth.csv")

	response.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").Equal(instanceID)
	response.Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("links").Object().Value("self").Object().Value("id").Equal(jobID)
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/jobs/" + jobID + "$")

	// Raised a bug for this
	response.NotContainsKey("last_updated")
}
