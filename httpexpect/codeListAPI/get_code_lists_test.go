package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get a set of codes lists containing information about dimensions which are used for all datasets at the ONS
// 200 - A Json message containing a set of code lists
func TestGetCodeLists_ReturnsSetOfCodeLists(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Get a set of code lists", t, func() {

		response := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()

		// CPI Codes

		cpiCodeID := response.Value("items").Array().Element(0).Object().Value("id").String().Raw()

		response.Value("items").Array().Element(0).Object().Value("id").NotNull()
		response.Value("items").Array().Element(0).Object().Value("name").Equal("CPI Codes")
		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("id").Equal(cpiCodeID)

		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains(cpiCodeID)
		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("codes").Object().Value("id").Equal("code")
		response.Value("items").Array().Element(0).Object().Value("links").Object().Value("codes").Object().Value("href").String().Contains(cpiCodeID).Contains("codes")

		// Time

		timeID := response.Value("items").Array().Element(1).Object().Value("id").String().Raw()

		response.Value("items").Array().Element(1).Object().Value("id").NotNull()
		response.Value("items").Array().Element(1).Object().Value("name").Equal("time")
		response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("id").Equal(timeID)

		response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains(timeID)
		response.Value("items").Array().Element(1).Object().Value("links").Object().Value("codes").Object().Value("id").Equal("code")
		response.Value("items").Array().Element(1).Object().Value("links").Object().Value("codes").Object().Value("href").String().Contains(timeID).Contains("codes")

		// Geography

		geographyID := response.Value("items").Array().Element(2).Object().Value("id").String().Raw()

		response.Value("items").Array().Element(2).Object().Value("id").NotNull()
		response.Value("items").Array().Element(2).Object().Value("name").Equal("geography")
		response.Value("items").Array().Element(2).Object().Value("links").Object().Value("self").Object().Value("id").Equal(geographyID)

		response.Value("items").Array().Element(2).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains(geographyID)
		response.Value("items").Array().Element(2).Object().Value("links").Object().Value("codes").Object().Value("id").Equal("code")
		response.Value("items").Array().Element(2).Object().Value("links").Object().Value("codes").Object().Value("href").String().Contains(geographyID).Contains("codes")

	})
}

// 400 - Missing parameters within request
// Need to write test for the above response
