package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Add a dimension to filter on with a list of options
// The dimension can only be added into the job if the state is still set to created
// otherwise 403 status code is returned

// 201 - The dimension was created
func TestPostAddDimension_CreatesDimension(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Add a dimension to the filter job", func() {

			filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).WithBytes([]byte(validPOSTAddDimensionToFilterJobJSON)).Expect().Status(http.StatusCreated)

		})

	})
}

// 400 - Invalid request body
func TestPostAddDimension_InvalidInput(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Add a dimension to the filter job with invalid JSON", func() {

			filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).WithBytes([]byte(invalidPOSTAddDimensionToFilterJobJSON)).Expect().Status(http.StatusBadRequest)

		})

	})
}

// 403 - Forbidden, the filter job has been locked as it has been submitted to be processed
func TestPostAddDimension_SubmittedJobForbiddenError(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter with submitted state", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterSubmittedJobJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Adding a dimension to a submitted job should throw forbidden error", func() {

			filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", filterJobID).WithBytes([]byte(validPOSTAddDimensionToFilterJobJSON)).Expect().Status(http.StatusForbidden)

		})

	})
}

// 404 - Filter job was not found
func TestPostAddDimension_FilterJobIDDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)
	Convey("A post request to add a dimension for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type", invalidFilterJobID).WithBytes([]byte(validPOSTAddDimensionToFilterJobJSON)).
			Expect().Status(http.StatusNotFound).Body().Contains("Filter job not found")
	})
}
