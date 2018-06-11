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

func TestSuccessfullyGetListOfImportJobs(t *testing.T) {

	var docs []*mongo.Doc

	importCreateJobDoc := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "_id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	importSubmittedJobDoc := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "_id",
		Value:      "01C24F0D-24BE-479F-962B-C76BCCD0AD00",
		Update:     validSubmittedImportJobData,
	}

	docs = append(docs, importCreateJobDoc, importSubmittedJobDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given an import job exists", t, func() {
		Convey("When get jobs is called with an authenticated request", func() {
			Convey("Then the response returns status OK (200)", func() {

				response := importAPI.GET("/jobs").
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).
					JSON().Array()
				checkImportJobsResponse(response)
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToGetListOfImportJobs(t *testing.T) {

	var docs []*mongo.Doc

	importCreateJobDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	docs = append(docs, importCreateJobDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given an import job exists", t, func() {
		Convey("When get jobs is called with no Authorization header", func() {
			Convey("Then the response returns status Unauthorized (401)", func() {
				importAPI.GET("/jobs").
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	Convey("Given an import job exists", t, func() {
		Convey("when get jobs is called with an unauthorised Authorization header", func() {
			Convey("Then the response returns status 401 unauthorised", func() {
				importAPI.GET("/jobs").
					WithHeader(serviceAuthTokenName, unauthorisedServiceAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	Convey("Given no import job are in a state of submitted", t, func() {
		Convey("When get jobs is called with an authenticated request", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.GET("/jobs?state=submitted").
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("requested resource not found")
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func checkImportJobsResponse(response *httpexpect.Array) {

	response.Element(0).Object().Value("id").Equal(jobID)
	response.Element(0).Object().Value("recipe").Equal("2080CACA-1A82-411E-AA46-F00804968E78")
	response.Element(0).Object().Value("state").Equal("Created")

	//Raised bug for this
	response.Element(0).Object().Value("files").Array().Element(0).Object().Value("alias_name").Equal("v4")

	response.Element(0).Object().Value("files").Array().Element(0).Object().Value("url").Equal("https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/CPIGrowth.csv")

	response.Element(0).Object().Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").Equal(instanceID)
	response.Element(0).Object().Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Element(0).Object().Value("links").Object().Value("self").Object().Value("id").Equal(jobID)
	response.Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/jobs/" + jobID + "$")

	// Raised a bug for this
	response.Element(0).Object().ContainsKey("last_updated")

	response.Element(1).Object().Value("id").Equal("01C24F0D-24BE-479F-962B-C76BCCD0AD00")
	response.Element(1).Object().Value("recipe").Equal("6C9D2696-131F-40C3-B598-12200C90415C")
	response.Element(1).Object().Value("state").Equal("Submitted")

	//Raised bug for this
	response.Element(1).Object().Value("files").Array().Element(0).Object().Value("alias_name").Equal("v4")

	response.Element(1).Object().Value("files").Array().Element(0).Object().Value("url").Equal("https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/CPIGrowth.csv")

	response.Element(1).Object().Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").Equal(instanceID)
	response.Element(1).Object().Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Element(1).Object().Value("links").Object().Value("self").Object().Value("id").Equal("01C24F0D-24BE-479F-962B-C76BCCD0AD00")
	response.Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/jobs/01C24F0D-24BE-479F-962B-C76BCCD0AD00$")

	// Raised a bug for this
	response.Element(1).Object().ContainsKey("last_updated")
}
