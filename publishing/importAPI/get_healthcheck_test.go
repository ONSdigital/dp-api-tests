package importAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHealthcheck(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request with no authentication headers", t, func() {
		Convey("When get healthcheck is called", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.
					GET("/healthcheck").
					Expect().Status(http.StatusOK)
			})
		})
	})
}
