package importAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get a list of all jobs
// Lists can be filtered by the job state
// 200 - A list of jobs has been returned
func TestGetJobs_ReturnsListOfJobs(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given an existing job", t, func() {

		importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
			Expect().Status(http.StatusCreated)

		Convey("Get a list of all jobs", func() {

			importAPI.GET("/jobs").Expect().Status(http.StatusOK).JSON().NotNull()

		})

	})
}
