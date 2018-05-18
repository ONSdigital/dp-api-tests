package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetHealthcheck(t *testing.T) {

	datasetAPIClient := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given the DatasetAPI is running", t, func() {
		Convey("When you ask the healthcheck endpoint", func() {
			Convey("Then you should see a status of OK", func() {

				response := datasetAPIClient.GET("/healthcheck", nil).Expect().Status(http.StatusOK).JSON().Object()
				response.Value("status").Equal("OK")
			})
		})
	})
}
