package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Add an option to a dimension to filter on
// 201 - Option was added
func TestPostAddDimensionOptions_AddsOptions(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Add an option to a dimension to filter on and Verify options are added", func() {

			filterAPI.POST("/filters/{filter_job_id}/dimensions/age/options/28", filterJobID).Expect().Status(http.StatusCreated)
			filterAPI.POST("/filters/{filter_job_id}/dimensions/sex/options/unknown", filterJobID).Expect().Status(http.StatusCreated)
			filterAPI.POST("/filters/{filter_job_id}/dimensions/Goods and services/options/welfare", filterJobID).Expect().Status(http.StatusCreated)
			filterAPI.POST("/filters/{filter_job_id}/dimensions/time/options/February 2007", filterJobID).Expect().Status(http.StatusCreated)

			ageDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			sexDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			goodsAndServicesDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			timeDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			ageDimResponse.Element(1).Object().Value("option").Equal("28")
			sexDimResponse.Element(2).Object().Value("option").Equal("unknown")
			goodsAndServicesDimResponse.Element(3).Object().Value("option").Equal("welfare")
			timeDimResponse.Element(5).Object().Value("option").Equal("February 2007")
		})

	})
}

// BUG RAISED
// 400 - Filter job was not found
func TestPostAddDimensionOptions_FilterJobIDDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)

	Convey("A post request to add an option to a dimension for a filter job that does not exist returns 404 not found", t, func() {

		filterAPI.POST("/filters/{filter_job_id}/dimensions/age/options/30", invalidFilterJobID).
			Expect().Status(http.StatusBadRequest)
	})
}

// 403 - Forbidden, the filter job has been locked as it has been submitted to be processed
func TestPostAddDimensionOptions_SubmittedJobForbiddenError(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given an existing filter with submitted state", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterSubmittedJobJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Updating a submitted job throws forbidden error", func() {

			filterAPI.POST("/filters/{filter_job_id}/dimensions/Residence Type/options/rent", filterJobID).Expect().Status(http.StatusForbidden)

		})

	})
}

//BUG RAISED
// 404 - Dimension name was not found
func TestPostAddDimensionOptions_DimensionNameDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	Convey("A post request for a filter job to add an option for a dimension that did not exists", t, func() {

		filterAPI.POST("/filters/{filter_job_id}/dimensions/agef/options/44", filterJobID).
			Expect().Status(http.StatusNotFound)
	})
}
