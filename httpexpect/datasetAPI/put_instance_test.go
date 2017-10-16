package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPutInstance_UpdatesInstance(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated)

		Convey("Create an Instance", func() {

			response := datasetAPI.POST("/instances").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateInstanceJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()

			instanceID := response.Value("id").String().Raw()

			Convey("Update an instance properties", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPUTUpdateInstanceJSON)).
					Expect().Status(http.StatusOK).JSON().Object()

			})
		})
	})
}
