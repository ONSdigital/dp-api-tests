package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPUTVersion_UpdatesVersion(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated)

		Convey("Get a list of all editions of a dataset", func() {

			response := datasetAPI.GET("/datasets/{id}/editions", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			edition := response.Value("items").Array().Element(0).Object().Value("edition").String().Raw()

			Convey("Get an edition of a dataset", func() {

				datasetAPI.GET("/datasets/{id}/editions/{edition}", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
					Expect().Status(http.StatusOK)

				Convey("Get a list of all versions from an edition of a dataset", func() {

					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
						Expect().Status(http.StatusOK)

					Convey("Get a specific version and edition of a dataset", func() {

						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateInstanceJSON)).
							Expect().Status(http.StatusOK)

						Convey("Update a version for an edition of a dataset", func() {

							datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPUTUpdateVersionJSON)).
								Expect().Status(http.StatusOK)
						})
					})
				})
			})
		})
	})

}
