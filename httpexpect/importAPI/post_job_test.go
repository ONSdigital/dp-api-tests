package importAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestSuccessfullyPostImportJob(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a requirement to create an import job exists", t, func() {
		Convey("When a post request with a valid json", func() {
			Convey("Then the response returns import job created (201)", func() {

				response := importAPI.POST("/jobs").WithBytes([]byte(validPOSTCreateJobJSON)).
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

				if err := mongo.Teardown(database, collection, "id", importJobID); err != nil {
					if err != mgo.ErrNotFound {
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

				importAPI.POST("/jobs").WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})
}
