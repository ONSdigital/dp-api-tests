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

func TestSuccessfullyGetListOfImportJobs(t *testing.T) {

	var docs []*mongo.Doc

	importCreateJobDoc := &mongo.Doc{
		Database:   "imports",
		Collection: "imports",
		Key:        "_id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	importSubmittedJobDoc := &mongo.Doc{
		Database:   "imports",
		Collection: "imports",
		Key:        "_id",
		Value:      "01C24F0D-24BE-479F-962B-C76BCCD0AD00",
		Update:     validSubmittedImportJobData,
	}

	docs = append(docs, importCreateJobDoc, importSubmittedJobDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	// These tests needs to refine when authentication was handled in the code.
	Convey("Given an import job exists", t, func() {
		Convey("When a request to get a list of all jobs and the user is authenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				response := importAPI.GET("/jobs").WithHeader(internalToken, internalTokenID).Expect().Status(http.StatusOK).JSON().Array()
				checkImportJobsResponse(response)
			})
		})

		Convey("When a request to get a list of all jobs and the user is authenticated", func() {
			Convey("Then the response returns status OK (200)", func() {

				response := importAPI.GET("/jobs").Expect().Status(http.StatusOK).JSON().Array()
				checkImportJobsResponse(response)
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
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
	response.Element(0).Object().NotContainsKey("last_updated")

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
	response.Element(1).Object().NotContainsKey("last_updated")
}
