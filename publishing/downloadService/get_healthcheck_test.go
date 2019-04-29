package downloadService

import (
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"fmt"
)

func TestHealthcheck(t *testing.T) {

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given a request with no authentication headers", t, func() {
		Convey("When get healthcheck is called", func() {
			Convey("Then the response returns a status of either 200 or 429", func() {

				response := downloadService.
					GET("/healthcheck").
					WithHeader(authHeader, serviceToken).
					Expect().Raw()

				isValidResponse := false

				if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusTooManyRequests {
					isValidResponse = true
				}

				fmt.Println(response.StatusCode)

				So(isValidResponse, ShouldBeTrue)
			})
		})
	})
}
