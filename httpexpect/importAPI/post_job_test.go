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

// Create an import job
// To import a dataset a job must be created first. To do this a data baker recipe is needed and the number
// of instances which the recipe creates. Once a job is created files can be added to the job and the state
// of the job can be changed.

// 201 - An import job was successfully created

func TestSuccessfullyPostImportJob(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a valid json input to create a job", t, func() {

		Convey("The jobs endpoint returns 201 created", func() {

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
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}
			}
		})
	})
}

// 400 - Invalid json message was sent to the API
func TestFailureToPostImportJob(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Fail to create a import job due to an invalid json body", t, func() {

		importAPI.POST("/jobs").WithBytes([]byte("{")).
			Expect().Status(http.StatusBadRequest)
	})
}
