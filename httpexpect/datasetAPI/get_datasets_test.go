package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetListOfDatasets_ReturnsListOfDatasets(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of datasets", t, func() {

		response := datasetAPI.GET("/datasets").
			Expect().Status(http.StatusOK).JSON().Object()

		response.Value("items").Array().Element(0).Object().Value("contact").Object().Value("email").Equal("jsinclair@test.co.uk")
		response.Value("items").Array().Element(0).Object().Value("contact").Object().Value("name").Equal("john sinclair")
		response.Value("items").Array().Element(0).Object().Value("contact").Object().Value("telephone").Equal("01633 123456")

		response.Value("items").Array().Element(0).Object().Value("collection_id").Equal("95c4669b-3ae9-4ba7-b690-87e890a1c543")
		response.Value("items").Array().Element(0).Object().Value("description").Equal("census covers the ethnicity of people living in the uk")

		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("editions").Object().Value("href").String().Contains("editions").Contains("95c4669b-3ae9-4ba7-b690-87e890a1c67c")

		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("latest_version").Object().Value("id").Equal("63294ed7-dccf-4f30-ad57-62365f038fb7")

		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("latest_version").Object().Value("href").String().Contains("95c4669b-3ae9-4ba7-b690-87e890a1c67c").Contains("editions").Contains("versions")

		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains("95c4669b-3ae9-4ba7-b690-87e890a1c67c")

		response.Value("items").Array().Element(0).Object().Value("next_release").Equal("2017-08-23")
		response.Value("items").Array().Element(0).Object().Value("periodicity").Equal("yearly")

		response.Value("items").Array().Element(0).Object().Value("publisher").Object().Value("name").Equal("The office of national statistics")
		response.Value("items").Array().Element(0).Object().Value("publisher").Object().Value("type").Equal("goverment department")
		response.Value("items").Array().Element(0).Object().Value("publisher").Object().Value("href").Equal("https://www.ons.gov.uk/")

		response.Value("items").Array().Element(0).Object().Value("state").Equal("published")
		response.Value("items").Array().Element(0).Object().Value("theme").Equal("population")
		response.Value("items").Array().Element(0).Object().Value("title").Equal("CPI")
	})
}
