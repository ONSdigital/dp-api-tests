package filterAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Remove an option from a dimension
// 200 - Option was removed
func TestDeleteRemoveDimensionOptions_RemovesOption(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Remove an option to a dimension to filter on and Verify options are removed", func() {

			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/sex/options/male", filterJobID).Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/Goods and services/options/communication", filterJobID).Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/time/options/April 1997", filterJobID).Expect().Status(http.StatusOK)

			// The below step is verifying, if a single option is there for a dimension and if you delete that single option,
			// not only option but dimension also deleted.
			// As you cant have a dimension with out an option.
			filterAPI.GET("/filters/{filter_job_id}/dimensions/age/options", filterJobID).Expect().Status(http.StatusBadRequest)
			sexDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/sex/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			goodsAndServicesDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/Goods and services/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()
			timeDimResponse := filterAPI.GET("/filters/{filter_job_id}/dimensions/time/options", filterJobID).Expect().Status(http.StatusOK).JSON().Array()

			sexDimResponse.Element(0).Object().Value("option").NotEqual("male").Equal("female")

			goodsAndServicesDimResponse.Element(0).Object().Value("option").NotEqual("communication").Equal("Education")
			goodsAndServicesDimResponse.Element(1).Object().Value("option").NotEqual("communication").Equal("health")
			timeDimResponse.Element(0).Object().Value("option").NotEqual("April 1997").Equal("March 1997")
			timeDimResponse.Element(1).Object().Value("option").NotEqual("April 1997").Equal("June 1997")
			timeDimResponse.Element(2).Object().Value("option").NotEqual("April 1997").Equal("September 1997")
			timeDimResponse.Element(3).Object().Value("option").NotEqual("April 1997").Equal("December 1997")

		})

	})
}

// 400 - This error code could be one or more of:
// . Filter job was not found
// . Dimension name was not found
func TestDeleteRemoveDimensionOptions_FilterJobIDAndDimensionDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	invalidFilterJobID := strings.Replace(filterJobID, "-", "", 9)

	Convey("A delete request to remove an option of a dimension of a filter job that does not exist returns 400 error", t, func() {

		filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/27", invalidFilterJobID).
			Expect().Status(http.StatusBadRequest)
	})
	Convey("A delete request to remove an option of a dimension that does not exist returns 400 error", t, func() {

		filterAPI.DELETE("/filters/{filter_job_id}/dimensions/ages/options/27", filterJobID).
			Expect().Status(http.StatusBadRequest)
	})
}

// 403 - Forbidden, the filter job has been locked as it has been submitted to be processed
func TestDeleteRemoveDimensionOptions_SubmittedJobForbiddenError(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter with submitted state", t, func() {

		expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterSubmittedJobJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()
		filterJobID := expected.Value("filter_job_id").String().Raw()
		Convey("Deleting an option for a submitted job should throw forbidden error", func() {

			filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/27", filterJobID).Expect().Status(http.StatusForbidden)

		})

	})
}

// 404 - Dimension option was not found
func TestDeleteRemoveDimensionOptions_DimensionOptionDoesNotExists(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	expected := filterAPI.POST("/filters").WithBytes([]byte(validPOSTMultipleDimensionsCreateFilterJSON)).
		Expect().Status(http.StatusCreated).JSON().Object()
	filterJobID := expected.Value("filter_job_id").String().Raw()

	Convey("A delete request for a filter job to remove an option for a dimension that did not exists", t, func() {

		filterAPI.DELETE("/filters/{filter_job_id}/dimensions/age/options/44", filterJobID).
			Expect().Status(http.StatusNotFound)
	})
}
