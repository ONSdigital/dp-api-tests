package importAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get a job
// Get information about a single job
// 200 - Return a single jobs information
func TestGetJob_ReturnsSingleJob(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given an existing job", t, func() {

		expected := importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		jobID := expected.Value("id").String().Raw()
		expectedJobID := expected.Value("id").String().Raw()
		expectedRecipe := expected.Value("recipe").String().Raw()
		expectedState := expected.Value("state").String().Raw()
		expectedFileAliasName := expected.Value("files").Array().Element(0).Object().Value("alias_name").String().Raw()
		expectedFileURL := expected.Value("files").Array().Element(0).Object().Value("url").String().Raw()
		expectedIntanceID := expected.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").String().Raw()
		expectedIntanceLink := expected.Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").String().Raw()

		Convey("Get a specific job", func() {

			actual := importAPI.GET("/jobs/{id}", jobID).Expect().Status(http.StatusOK).JSON().Object()

			actualJobID := actual.Value("id").String().Raw()
			actualRecipe := actual.Value("recipe").String().Raw()
			actualState := actual.Value("state").String().Raw()
			actualFileAliasName := actual.Value("files").Array().Element(0).Object().Value("alias_name").String().Raw()
			actualFileURL := actual.Value("files").Array().Element(0).Object().Value("url").String().Raw()
			actualIntanceID := actual.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").String().Raw()
			actualIntanceLink := actual.Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").String().Raw()

			So(actualJobID, ShouldEqual, expectedJobID)
			So(actualRecipe, ShouldEqual, expectedRecipe)
			So(actualState, ShouldEqual, expectedState)
			So(actualFileAliasName, ShouldEqual, expectedFileAliasName)
			So(actualFileURL, ShouldEqual, expectedFileURL)

			So(actualIntanceID, ShouldEqual, expectedIntanceID)
			So(actualIntanceLink, ShouldEqual, expectedIntanceLink)
		})

	})
}

// 404 - JobId does not match any import jobs
func TestGetJob_JobIDDoesNotExists(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("A get request for a job that does not exist returns 404 not found", t, func() {

		importAPI.GET("/jobs/{id}", "c387798b12-0cb623-43d93-bc564-78d2284c684").
			Expect().Status(http.StatusNotFound)
	})
}
