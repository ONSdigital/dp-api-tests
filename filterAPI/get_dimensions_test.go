package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get all dimensions used by this filter
// Return a list of all dimensions which are going to be used to filter on
// 200 - A list of dimension URLs
func TestGetDimensions_AllDimensionUrlsReturns(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()

		Convey("Get the list of all dimensions in a filter job", func() {

			actual := filterAPI.GET("/filters/{filter_job_id}/dimensions", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			actual.Element(0).Object().Value("dimension_url").NotNull()
			actual.Element(0).Object().Value("name").Equal("age")
			actual.Element(1).Object().Value("dimension_url").NotNull()
			actual.Element(1).Object().Value("name").Equal("sex")
			actual.Element(2).Object().Value("dimension_url").NotNull()
			actual.Element(2).Object().Value("name").Equal("Goods and services")
			actual.Element(3).Object().Value("dimension_url").NotNull()
			actual.Element(3).Object().Value("name").Equal("time")
		})

	})
}

// 404- Filter job not found
func TestGetDimensionsForFilterJob_FilterJobIDDoesNotExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)
	Convey("A get request for dimensons for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions", invalidFilterJobID).
			Expect().Status(http.StatusNotFound)
	})
}
