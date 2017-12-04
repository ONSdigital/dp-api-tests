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

func TestSuccessfullyGetAnImportJob(t *testing.T) {

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

	// TODO Dont skip test once endpoint has been refactored
	// These tests needs to refine when authentication was handled in the code.
	SkipConvey("Given an import job exists", t, func() {
		Convey("When a request to get the job with a specific id and the user is authenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				response := importAPI.GET("/jobs/{id}", jobID).WithHeader(internalToken, internalTokenID).Expect().Status(http.StatusOK).JSON().Object()
				checkImportJobResponse(response)
			})
		})

		Convey("When a request to get the job with a specific id and the user is unauthenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				response := importAPI.GET("/jobs/{id}", jobID).Expect().Status(http.StatusOK).JSON().Object()
				checkImportJobResponse(response)

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

func TestFailureToGetAnImportJob(t *testing.T) {

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

	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given an import job exists", t, func() {
		Convey("When a request to get the job with id does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.GET("/jobs/{id}", invalidJobID).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound)
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
