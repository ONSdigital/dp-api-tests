package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetInstance_ReturnsSingleInstance(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create an Instance", t, func() {

		response := datasetAPI.POST("/instances").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateInstanceJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()

		instanceID := response.Value("id").String().Raw()

		Convey("A get request for an instance with instance id", func() {

			response := datasetAPI.GET("/instances/{id}", instanceID).Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(instanceID)
			response.Value("edition").Equal("2017")

			response.Value("headers").Array().Element(0).Equal("time")
			response.Value("headers").Array().Element(1).Equal("geography")

			response.Value("total_inserted_observations").Equal(1000)

			response.Value("links").Object().Value("job").Object().Value("id").Equal("042e216a-7822-4fa0-a3d6-e3f5248ffc35")
			response.Value("links").Object().Value("job").Object().Value("href").String().Match("(.+)/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35$")

			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/34B13D18-B4D8-4227-9820-492B2971E221$")

			response.Value("state").Equal("completed")
			response.Value("total_observations").Equal(1000)
			response.Value("last_updated").NotNull()

		})

	})
}
