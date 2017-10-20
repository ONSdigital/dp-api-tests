package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get a description of a filter job
// Get document describing the filter job
// 200 - The filter job was found and document is returned
func TestGetFilterJob_JobFoundAndDocumentReturned(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		expectedFilterJobID := expected.Value("filter_job_id").String().Raw()
		expectedDatasetFilterID := expected.Value("instance_id").String().Raw()
		expectedDimensionlistURL := expected.Value("dimension_list_url").String().Raw()
		expectedState := expected.Value("state").String().Raw()

		Convey("Get the description of a filter job", func() {

			actual := filterAPI.GET("/filters/{filter_job_id}", filterJobID).Expect().Status(http.StatusOK).JSON().Object()

			actualFilterJobID := actual.Value("filter_job_id").String().Raw()
			actualDatasetFilterID := actual.Value("instance_id").String().Raw()
			actualState := actual.Value("state").String().Raw()
			actualDimensionlistURL := actual.Value("dimension_list_url").String().Raw()

			// comparing the values given in POST with GET requests
			So(actualFilterJobID, ShouldEqual, expectedFilterJobID)
			So(actualDatasetFilterID, ShouldEqual, expectedDatasetFilterID)
			So(actualState, ShouldEqual, expectedState)
			So(actualDimensionlistURL, ShouldEqual, expectedDimensionlistURL)

		})

	})
}

// 404 - Filter job not found
func TestGetFilterJob_FilterJobIDDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("A get request for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.GET("/filters/{filter_job_id}", "c387798b1rf-0cb623e-43ddf4-bc5df45-78d2284c45fgh").
			Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
	})

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)

	Convey("A get request to a filter job that does not exists returns 404 not found", t, func() {

		filterAPI.GET("/filters/{filter_job_id}", invalidFilterJobID).
			Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
	})
}
