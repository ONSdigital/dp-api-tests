package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Check if a dimension exists within a filter job
// 204 - Dimension exists for filter job
func TestGetCheckDimensionExists_DimensionExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Check If age dimension exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/age", filterJobID).Expect().Status(http.StatusNoContent)

		})

		Convey("Check If sex dimension exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/sex", filterJobID).Expect().Status(http.StatusNoContent)

		})

		Convey("Check If goods and services dimension exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services", filterJobID).Expect().Status(http.StatusNoContent)

		})

		Convey("Check If time dimension exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/time", filterJobID).Expect().Status(http.StatusNoContent)

		})

	})
}

// 400 - Filter job was not found
func TestGetCheckDimensionsExists_FilterJobIDDoesNotExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)
	Convey("A get request for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions/age", invalidFilterJobID).
			Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
	})
}

// 404 - Dimension name was not found
func TestGetDimensionsExists_DimensionNameDoesNotExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	Convey("A get request for a filter job with a dimension name that does not exists", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions/agefdsfdsf", filterJobID).
			Expect().Status(http.StatusNotFound)
	})
}
