package importAPI

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetJobs(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given an existing job", t, func() {

		r := importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()

		expectedJobID := r.Value("id")

		fmt.Println(expectedJobID)

		Convey("Get a list of all jobs", func() {
			//json := importAPI.GET("/jobs").Expect().Status(http.StatusOK).JSON().Array().Element(1).Object()

			//json.Values().Contains(expectedJobID)

			importAPI.GET("/jobs").Expect().Status(http.StatusOK).JSON().NotNull()

		})

	})
}
