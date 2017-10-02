package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Remove a dimension and any options set within the dimension
// The dimension was removed
// 200 - The dimension was removed
func TestDeleteRemoveDimension_RemovesDimension(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()

		Convey("Remove a dimension and any options set within the dimension", func() {

			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/Goods and services", filterJobID).Expect().Status(http.StatusOK)

		})

	})
}

// 400 - Filter job was not found
func TestDeleteRemoveDimension_FilterJobIDDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)
	Convey("A get request for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age", invalidFilterJobID).
			Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter job not found")
	})
}

// 403 - Forbidden, the filter job has been locked as it has been submitted to be processed
func TestDeleteRemoveDimension_UpdatingASubmittedFilterJobThrowsForbiddenError(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter with submitted state", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterSubmittedJobJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Updating a submitted job thorws forbidden error", func() {

			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age", filterJobID).Expect().Status(http.StatusForbidden)

		})

	})
}

// 404 - Dimension name was not found
func TestDeleteRemoveDimension_DimensionDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Deleting a dimension that does not exists throws 404 error", func() {

			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/agewqeqw", filterJobID).Expect().Status(http.StatusNotFound)

		})

	})
}
