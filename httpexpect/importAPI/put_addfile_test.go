package importAPI

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPutAddFilesToJob(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given an existing job", t, func() {

		expected := importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		jobID := expected.Value("id").String().Raw()
		expectedFileAliasName := expected.Value("files").Array().Element(0).Object().Value("alias_name").String().Raw()
		expectedFileURL := expected.Value("files").Array().Element(0).Object().Value("url").String().Raw()

		fmt.Println(jobID)
		Convey("Add files to job", func() {

			importAPI.PUT("/jobs/{id}/files", jobID).WithBytes([]byte(validPUTAddFilesJSON)).Expect().Status(http.StatusOK)

			Convey("Verify files added to job", func() {

				updated := importAPI.GET("/jobs/{id}", jobID).Expect().Status(http.StatusOK).JSON().Object()
				updatedFileAliasName := updated.Value("files").Array().Element(0).Object().Value("alias_name").String().Raw()
				updatedFileURL := updated.Value("files").Array().Element(0).Object().Value("url").String().Raw()

				So(updatedFileAliasName, ShouldEqual, expectedFileAliasName)
				So(updatedFileURL, ShouldEqual, expectedFileURL)

			})
		})

	})
}
func TestPUTAddFilesToJob_InvalidInput(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given invalid json input to add files to a job", t, func() {

		Convey("Given an existing job", func() {

			expected := importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()
			jobID := expected.Value("id").String().Raw()

			fmt.Println(jobID)

			Convey("The jobs endpoint returns 400 invalid json message ", func() {

				importAPI.PUT("/jobs/{id}/files", jobID).WithBytes([]byte(invalidSyntaxJSON)).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})
}

func TestPutAddFilesToJob_JobIDDoesNotExists(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("A get request for a job that does not exist returns 404 not found", t, func() {

		importAPI.PUT("/jobs/{id}/files", "f708ca2232-641c12-40e62-b0dcsd-68db0aa0e007ed").WithBytes([]byte(validJSON)).
			Expect().Status(http.StatusNotFound)
	})
}
