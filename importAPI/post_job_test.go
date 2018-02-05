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

func TestSuccessfullyPostImportJob(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a requirement to create an import job exists", t, func() {
		Convey("When a post request with a valid json", func() {
			Convey("Then the response returns import job created (201)", func() {

				response := importAPI.POST("/jobs").WithHeader(tokenName, tokenSecret).WithBytes([]byte(validPOSTCreateJobJSON)).
					Expect().Status(http.StatusCreated).JSON().Object()

				importJobID := response.Value("id").String().Raw()

				response.Value("id").NotNull()
				response.Value("recipe").Equal("b944be78-f56d-409b-9ebd-ab2b77ffe187")
				response.Value("state").Equal("created")

				response.Value("files").Array().Element(0).Object().Value("alias_name").Equal("v4")
				response.Value("files").Array().Element(0).Object().Value("url").Equal("https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv")

				response.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").NotNull()
				response.Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").NotNull()

				response.Value("links").Object().Value("self").Object().Value("id").Equal("")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/jobs/" + importJobID + "$")

				response.Value("last_updated").NotNull()

				importJob := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: collection,
					Key:        "id",
					Value:      importJobID,
				}

				if err := mongo.Teardown(importJob); err != nil {
					if err != mgo.ErrNotFound {
						log.ErrorC("Failed to tear down test data", err, nil)
						os.Exit(1)
					}
				}
			})
		})
	})
}

func TestFailureToPostImportJob(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a requirement to create an import job exists", t, func() {
		Convey("When a post request with an invalid json", func() {
			Convey("Then the response returns bad request (400)", func() {

				importAPI.POST("/jobs").WithHeader(tokenName, tokenSecret).WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})
}

func TestPostImportJobWithNoAuthentication(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a requirement to create an import job exists", t, func() {
		Convey("When a post request with an valid json body but no authentication header", func() {
			Convey("Then the response returns not found (404)", func() {

				importAPI.POST("/jobs").WithBytes([]byte(validPOSTCreateJobJSON)).
					Expect().Status(http.StatusNotFound)

			})
		})
	})
}
