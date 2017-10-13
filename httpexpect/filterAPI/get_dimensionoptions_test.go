package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get all options from a dimension which have been set
// Get a list of all options which will be used to filter the dimension
// 200 - A list of all options for a dimension was returned
func TestGetDimensionOptions_ListOfOptionsReturns(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()

		Convey("Check If age dimension options exists for the filter job", func() {

			response := filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			response.Element(0).Object().Value("option").Equal("27")
			response.Element(0).Object().Value("dimension_option_url").NotNull()
		})

		Convey("Check If sex dimension options exists for the filter job", func() {

			response := filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			response.Element(0).Object().Value("option").Equal("female")
			response.Element(0).Object().Value("dimension_option_url").NotNull()

			response.Element(1).Object().Value("option").Equal("male")
			response.Element(1).Object().Value("dimension_option_url").NotNull()
		})

		Convey("Check If goods and services dimension options exists for the filter job", func() {

			response := filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			response.Element(0).Object().Value("option").Equal("Education")
			response.Element(0).Object().Value("dimension_option_url").NotNull()

			response.Element(1).Object().Value("option").Equal("health")
			response.Element(1).Object().Value("dimension_option_url").NotNull()

			response.Element(2).Object().Value("option").Equal("communication")
			response.Element(2).Object().Value("dimension_option_url").NotNull()
		})

		Convey("Check If time dimension options exists for the filter job", func() {

			response := filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			//response.Value("options").Elements("March 1997", "April 1997")
			response.Element(0).Object().Value("option").Equal("March 1997")
			// response.Element(0).Object().Value("dimension_option_url").NotNull()

			// response.Element(1).Object().Value("option").Equal("April 1997")
			// response.Element(1).Object().Value("dimension_option_url").NotNull()

			// response.Element(2).Object().Value("option").Equal("June 1997")
			// response.Element(2).Object().Value("dimension_option_url").NotNull()

			// response.Element(3).Object().Value("option").Equal("September 1997")
			// response.Element(3).Object().Value("dimension_option_url").NotNull()

			// response.Element(4).Object().Value("option").Equal("December 1997")
			// response.Element(4).Object().Value("dimension_option_url").NotNull()
		})

	})
}

// 400 - Filter job was not found
func TestGetDimensionOptions_FilterJobIDDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	Convey("A post request to add a dimension for a filter job that does not exist returns 404 not found", t, func() {

		invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)

		filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", invalidFilterJobID).
			Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
	})
}

// 404 - Dimension name was not found
func TestGetDimensionOptions_DimensionNameDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	Convey("A get request for a filter job with a dimension name and options that does not exists", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions/agefdsfdsf/options", filterJobID).
			Expect().Status(http.StatusNotFound)
	})
}
