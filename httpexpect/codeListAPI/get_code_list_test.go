package codeListAPI

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get information about a code list
// 200 - Json object containing information about the code list
func TestGetACodeList_ReturnsSingleCodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, config.CodeListAPIURL())

	Convey("Given a set of code lists", t, func() {

		response := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()

		cpiCodeID := response.Value("items").Array().Element(0).Object().Value("id").String().Raw()
		timeID := response.Value("items").Array().Element(1).Object().Value("id").String().Raw()
		geographyID := response.Value("items").Array().Element(2).Object().Value("id").String().Raw()

		Convey("Then get the information about CPI code list", func() {
			cpiCodesResponse := codeListAPI.GET("/code-lists/{id}", cpiCodeID).Expect().Status(http.StatusOK).JSON().Object()

			cpiCodesResponse.Value("id").Equal(cpiCodeID)
			cpiCodesResponse.Value("name").Equal("CPI Codes")
			cpiCodesResponse.Value("links").Object().Value("self").Object().Value("id").Equal(cpiCodeID)

			cpiCodesResponse.Value("links").Object().Value("self").Object().Value("href").String().Contains(cpiCodeID)
			cpiCodesResponse.Value("links").Object().Value("codes").Object().Value("id").Equal("code")
			cpiCodesResponse.Value("links").Object().Value("codes").Object().Value("href").String().Contains(cpiCodeID).Contains("codes")

		})

		Convey("And get the information about time code list", func() {
			timeListResponse := codeListAPI.GET("/code-lists/{id}", timeID).Expect().Status(http.StatusOK).JSON().Object()

			timeListResponse.Value("id").Equal(timeID)
			timeListResponse.Value("name").Equal("time")
			timeListResponse.Value("links").Object().Value("self").Object().Value("id").Equal(timeID)

			timeListResponse.Value("links").Object().Value("self").Object().Value("href").String().Contains(timeID)
			timeListResponse.Value("links").Object().Value("codes").Object().Value("id").Equal("code")
			timeListResponse.Value("links").Object().Value("codes").Object().Value("href").String().Contains(timeID).Contains("codes")

		})

		Convey("And get the information about geography code list", func() {
			geographyListResponse := codeListAPI.GET("/code-lists/{id}", geographyID).Expect().Status(http.StatusOK).JSON().Object()

			geographyListResponse.Value("id").Equal(geographyID)
			geographyListResponse.Value("name").Equal("geography")
			geographyListResponse.Value("links").Object().Value("self").Object().Value("id").Equal(geographyID)

			geographyListResponse.Value("links").Object().Value("self").Object().Value("href").String().Contains(geographyID)
			geographyListResponse.Value("links").Object().Value("codes").Object().Value("id").Equal("code")
			geographyListResponse.Value("links").Object().Value("codes").Object().Value("href").String().Contains(geographyID).Contains("codes")

		})

	})
}

// Bug Raised
// 404 - Code list not found
func TestGetACodeList_InvalidCodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, config.CodeListAPIURL())

	Convey("Given a set of code lists", t, func() {

		codeListID := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object().Value("items").Array().Element(0).Object().Value("id").String().Raw()

		Convey("With an invalid code list throws 404 error", func() {
			invalidCodeListID := strings.Replace(codeListID, "-", "", 9)

			codeListAPI.GET("/code-lists/{id}", invalidCodeListID).Expect().Status(http.StatusNotFound)

		})

	})
}
