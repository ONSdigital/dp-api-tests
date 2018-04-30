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

func TestSuccessfullyGetAnImportJob(t *testing.T) {

	importCreateJobDoc := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	if err := mongo.Setup(importCreateJobDoc); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request for a job that exists", t, func() {
		Convey("When get job is called", func() {
			Convey("Then the response returns status OK (200)", func() {

				response := importAPI.GET("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).
					JSON().Object()
				checkImportJobResponse(response)
			})
		})
	})

	if err := mongo.Teardown(importCreateJobDoc); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToGetAnImportJob(t *testing.T) {

	importCreateJobDoc := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	if err := mongo.Setup(importCreateJobDoc); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request for a job that does not exist", t, func() {
		Convey("When get job is called", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.GET("/jobs/{id}", invalidJobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	if err := mongo.Teardown(importCreateJobDoc); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestGetAnImportJobUnauthorised(t *testing.T) {

	importCreateJobDoc := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	if err := mongo.Setup(importCreateJobDoc); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request with no Authorization header", t, func() {
		Convey("When get job is called", func() {
			Convey("Then the response returns status unauthorized (401)", func() {
				importAPI.GET("/jobs/{id}", jobID).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	Convey("Given a request with an unauthorised Authorization header", t, func() {
		Convey("When get job is called ", func() {
			Convey("Then the response returns status 401 unauthorised", func() {
				importAPI.GET("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, unauthorisedServiceAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	if err := mongo.Teardown(importCreateJobDoc); err != nil {
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

	response.ContainsKey("last_updated")
}
