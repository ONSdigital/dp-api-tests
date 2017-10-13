package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Update filter job
// Update the filter job by providing new properties
// 200 - The filter job has been updated
func TestPUTUpdateFilterJob_FilterJobUpdates(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		expectedFilterJobID := expected.Value("filter_job_id").String().Raw()
		expectedDatasetFilterID := expected.Value("dataset_filter_id").String().Raw()
		expectedDimensionlistURL := expected.Value("dimension_list_url").String().Raw()
		expectedState := expected.Value("state").String().Raw()

		Convey("Update filter job with new properties", func() {

			filterAPI.PUT("/filters/{filter_job_id}", filterJobID).WithBytes([]byte(validPUTUpdateFilterJobJSON)).Expect().Status(http.StatusOK)

			Convey("Verify filter job state is updated", func() {

				updated := filterAPI.GET("/filters/{filter_job_id}", filterJobID).Expect().Status(http.StatusOK).JSON().Object()

				updatedFilterJobID := updated.Value("filter_job_id").String().Raw()
				updatedDatasetFilterID := updated.Value("dataset_filter_id").String().Raw()
				updatedState := updated.Value("state").String().Raw()
				updatedDimensionListURL := updated.Value("dimension_list_url").String().Raw()

				So(updatedState, ShouldNotEqual, expectedState)
				So(updatedFilterJobID, ShouldEqual, expectedFilterJobID)
				So(updatedDatasetFilterID, ShouldEqual, expectedDatasetFilterID)
				So(updatedDimensionListURL, ShouldEqual, expectedDimensionlistURL)

			})
			Convey("Verify filter job dimension is updated", func() {

				dimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

				dimResponse.Element(0).Object().Value("name").Equal("age")
				dimResponse.Element(0).Object().Value("dimension_url").NotNull()
			})
			Convey("Verify filter job dimension options are updated", func() {

				dimOptionsResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

				dimOptionsResponse.Element(0).Object().Value("option").Equal("27")
				dimOptionsResponse.Element(0).Object().Value("dimension_option_url").NotNull()

				dimOptionsResponse.Element(1).Object().Value("option").Equal("28")
				dimOptionsResponse.Element(1).Object().Value("dimension_option_url").NotNull()
			})
		})

	})
}

// 400 -Invalid request body
func TestPUTUpdateFilterJob_InvalidInput(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given invalid json input to update filter job", t, func() {

		Convey("Given an existing filter job", func() {

			expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()
			filterJobID := expected.Value("filter_job_id").String().Raw()

			Convey("The filter job endpoint returns 400 invalid json message ", func() {

				filterAPI.PUT("/filters/{filter_job_id}", filterJobID).WithBytes([]byte(invalidSyntaxJSON)).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})
}

// 403 - Forbidden, the job has been locked as it has been submitted to be processed
func TestPUTUpdateSubmittedFilterJob_ForbiddenError(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter with submitted state", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterSubmittedJobJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()

		Convey("Updating a submitted job throws forbidden error", func() {

			filterAPI.PUT("/filters/{filter_job_id}", filterJobID).WithBytes([]byte(validPUTUpdateFilterJobJSON)).Expect().Status(http.StatusForbidden)

		})

	})
}

// 404 - Filter job not found
func TestPUTUpdateFilterJob_FilterJobIDDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)
	Convey("A post request to add a dimension for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.PUT("/filters/{filter_job_id}", invalidFilterJobID).WithBytes([]byte(validPUTUpdateFilterJobJSON)).
			Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
	})
}

// 1 more tests need to write for 401 response
