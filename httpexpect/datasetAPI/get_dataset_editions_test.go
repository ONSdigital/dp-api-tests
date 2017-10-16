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

func TestGetListOfDatasetEditions_ReturnsListOfDatasetEditions(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of datasets", t, func() {

		response := datasetAPI.GET("/datasets").
			Expect().Status(http.StatusOK).JSON().Object()

		editionsHref := response.Value("items").Array().Element(0).Object().Value("links").Object().Value("editions").Object().Value("href").String().Raw()

		datasetID := strings.TrimLeft(strings.TrimRight(editionsHref, "/editions"), "http://localhost:22000/datasets/")

		Convey("Get a list of dataset editions", func() {

			response := datasetAPI.GET("/datasets/{id}/editions", datasetID).Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Element(0).Object().Value("edition").Equal("2016")
			response.Value("items").Array().Element(0).Object().Value("id").Equal("a051a058-58a9-4ba4-8374-fbb7315d3b78")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("dataset").Object().Value("href").String().Contains(datasetID)

			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains(datasetID).Contains("editions")
			response.Value("items").Array().Element(0).Object().Value("links").Object().Value("versions").Object().Value("href").String().Contains(datasetID).Contains("editions").Contains("versions")
			response.Value("items").Array().Element(0).Object().Value("state").Equal("published")

		})

	})
}
