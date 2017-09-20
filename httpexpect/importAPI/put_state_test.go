package importAPI

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUpdateJobState(t *testing.T) {

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

		Convey("Update job state", func() {

			importAPI.PUT("/jobs/{id}", jobID).WithBytes([]byte(validPUTJobJSON)).Expect().Status(http.StatusOK)

			Convey("Verify job state is updated", func() {

				updated := importAPI.GET("/jobs/{id}", jobID).Expect().Status(http.StatusOK).JSON().Object()

				updatedJobID := updated.Value("id").String().Raw()
				updatedRecipe := updated.Value("recipe").String().Raw()
				updatedState := updated.Value("state").String().Raw()
				updatedFileAliasName := updated.Value("files").Array().Element(0).Object().Value("alias_name").String().Raw()
				updatedFileURL := updated.Value("files").Array().Element(0).Object().Value("url").String().Raw()
				updatedIntanceID := updated.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").String().Raw()
				updatedIntanceLink := updated.Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").String().Raw()

				So(updatedJobID, ShouldEqual, expectedJobID)
				So(updatedRecipe, ShouldEqual, expectedRecipe)
				So(updatedState, ShouldNotEqual, expectedState)
				So(updatedFileAliasName, ShouldEqual, expectedFileAliasName)
				So(updatedFileURL, ShouldEqual, expectedFileURL)

				So(updatedIntanceID, ShouldEqual, expectedIntanceID)
				So(updatedIntanceLink, ShouldEqual, expectedIntanceLink)
			})
		})

	})
}
func TestPUTUpdateJobState_InvalidInput(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given invalid json input to update job state", t, func() {

		Convey("Given an existing job", func() {

			expected := importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()
			jobID := expected.Value("id").String().Raw()

			fmt.Println(jobID)

			Convey("The jobs endpoint returns 400 invalid json message ", func() {

				importAPI.PUT("/jobs/{id}", jobID).WithBytes([]byte(invalidSyntaxJSON)).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})
}

func TestPutUpdateJobState_JobIDDoesNotExists(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("A get request for a job that does not exist returns 404 not found", t, func() {

		importAPI.PUT("/jobs/{id}", "99cc5ba6wd-1827fg-407e23-849f36-7fb7ac8a422f5tg").WithBytes([]byte(validJSON)).
			Expect().Status(http.StatusNotFound)
	})
}
