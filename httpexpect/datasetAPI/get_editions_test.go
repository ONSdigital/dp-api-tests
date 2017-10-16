package datasetAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get a list of editions from a type of dataset
// 200 - A json list containing all editions for a dataset

func TestGetEditions_ReturnsListOfEditions(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated)

		Convey("Get a list of all editions of a dataset", func() {

			response := datasetAPI.GET("/datasets/{id}/editions", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateInstanceJSON)).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Element(0).Object().Value("edition").Equal("2017")
			response.Value("items").Array().Element(0).Object().Value("id").Equal("bed1f712-aadf-433a-8284-f7992e01ffc3")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/34B13D18-B4D8-4227-9820-492B2971E221$")

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/34B13D18-B4D8-4227-9820-492B2971E221/editions/2017$")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("versions").Object().Value("href").String().Match("(.+)/datasets/34B13D18-B4D8-4227-9820-492B2971E221/editions/2017/versions$")

			response.Value("items").Array().Element(0).Object().Value("state").Equal("created")

		})
	})
}

// Get a list of editions from a type of dataset
// 400 - No editions were found for the dataset id provided
func TestGetEditions_InvalidDatasetID(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated)

		Convey("Getting the list of editions with invalid dataset id throws 404 error", func() {

			invalidDatasetID := strings.Replace(datasetID, "-", "", 9)

			datasetAPI.GET("/datasets/{id}/editions", invalidDatasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusNotFound)

		})
	})
}
