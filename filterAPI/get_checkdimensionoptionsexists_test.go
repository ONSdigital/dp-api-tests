package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Check if a option exists within a dimension
// 204 - Option exists within the dimension
func TestGetCheckDimensionOptionsExists_OptionsExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Check If age dimension options exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).Expect().Status(http.StatusNoContent)
		})

		Convey("Check If sex dimension options exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options/male", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options/female", filterJobID).Expect().Status(http.StatusNoContent)

		})

		Convey("Check If goods and services dimension options exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options/Education", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options/health", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options/communication", filterJobID).Expect().Status(http.StatusNoContent)

		})

		Convey("Check If time dimension options exists for the filter job", func() {

			filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/March 1997", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/April 1997", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/June 1997", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/September 1997", filterJobID).Expect().Status(http.StatusNoContent)
			filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options/December 1997", filterJobID).Expect().Status(http.StatusNoContent)

		})

	})
}

// 400 - Filter job or dimension name not found
func TestGetCheckDimensionOptionsExists_FilterJobAndDimensionDoesNotExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)
	Convey("A GET request to check if an option exists for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options/27", invalidFilterJobID).
			Expect().Status(http.StatusBadRequest).Body().Contains("filter or dimension not found")
	})

	Convey("A GET request to check if an option exists for a dimension that does not exist returns 404 not found", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions/ages/options/27", filterJobID).
			Expect().Status(http.StatusBadRequest).Body().Contains("filter or dimension not found")
	})
}

// 404 - Dimension option was not found
func TestGetCheckDimensionOptionsExists_DimensionOptionDoesNotExists(t *testing.T) {

	setupDatastores()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	Convey("A get request for a filter job with a dimension option that did not exists", t, func() {

		filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options/unknown", filterJobID).
			Expect().Status(http.StatusNotFound)
	})
}
