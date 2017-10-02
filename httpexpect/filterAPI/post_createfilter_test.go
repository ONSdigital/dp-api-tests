package filterAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Create a filter job for a dataset
// Create a job so that dimensions can be added to filter a dataset
// 201- Job was created
func TestPostCreateFilter_CreatesFilter(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given a valid json input to create a filter", t, func() {

		Convey("The filters endpoint returns 201 created", func() {

			response := filterAPI.POST("/filters").WithBytes([]byte(validPOSTCreateFilterJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()

			response.Value("filter_job_id").NotNull()
			response.Value("dimension_list_url").NotNull()
		})

	})

}

// 400 -Invalid request body
func TestPostCreateFilter_InvalidInput(t *testing.T) {

	filterAPI := httpexpect.New(t, config.FilterAPIURL())

	Convey("Given invalid json input to create a filter", t, func() {

		Convey("The jobs endpoint returns 400 invalid json message ", func() {

			filterAPI.POST("/filters").WithBytes([]byte(invalidJSON)).
				Expect().Status(http.StatusBadRequest)
		})
	})
}
