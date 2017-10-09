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
func TestGetCodesWithInACodeList_ReturnsListOfCodesInACodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, config.CodeListAPIURL())

	Convey("Given a set of code lists", t, func() {

		response := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()

		cpiCodeID := response.Value("items").Array().Element(0).Object().Value("id").String().Raw()
		timeID := response.Value("items").Array().Element(1).Object().Value("id").String().Raw()
		geographyID := response.Value("items").Array().Element(2).Object().Value("id").String().Raw()

		Convey("Then get the list of all codes with in a CPI code-list", func() {

			cpiCodesResponse := codeListAPI.GET("/code-lists/{id}/codes", cpiCodeID).Expect().Status(http.StatusOK).JSON().Object()
			cpiFirstCode := cpiCodesResponse.Value("items").Array().Element(0).Object().Value("id").String().Raw()
			cpiCodesResponse.Value("items").Array().Element(0).Object().Value("id").NotNull()
			cpiCodesResponse.Value("items").Array().Element(0).Object().Value("label").Equal("Sugar, jam, syrups, chocolate and confectionery")
			cpiCodesResponse.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal(cpiCodeID)

			cpiCodesResponse.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Contains(cpiCodeID)

			cpiCodesResponse.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains(cpiCodeID).Contains("codes").Contains(cpiFirstCode)
			cpiCodesResponse.Value("items").Array().Length().Equal(137)

		})

		Convey("And get the list of all codes with in a time code-list", func() {

			timeListResponse := codeListAPI.GET("/code-lists/{id}/codes", timeID).Expect().Status(http.StatusOK).JSON().Object()
			timeFirstCode := timeListResponse.Value("items").Array().Element(0).Object().Value("id").String().Raw()

			timeListResponse.Value("items").Array().Element(0).Object().Value("id").NotNull()
			timeListResponse.Value("items").Array().Element(0).Object().Value("label").Equal("1998.07")
			timeListResponse.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal(timeID)

			timeListResponse.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Contains(timeID)

			timeListResponse.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Contains(timeID).Contains("codes").Contains(timeFirstCode)
			timeListResponse.Value("items").Array().Length().Equal(254)

		})

		Convey("And get the list of all codes with in a geography code", func() {
			geographyListResponse := codeListAPI.GET("/code-lists/{id}/codes", geographyID).Expect().Status(http.StatusOK).JSON().Object()
			geographyListResponse.Value("items").Null()
		})

	})
}

// Bug Raised
// 404 - Code list not found
func TestGetCodesWithInACodeList_InvalidCodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, config.CodeListAPIURL())

	Convey("Given a set of code lists", t, func() {

		codeListID := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object().Value("items").Array().Element(0).Object().Value("id").String().Raw()

		Convey("With an invalid code list throws 404 error", func() {
			invalidCodeListID := strings.Replace(codeListID, "-", "", 9)

			codeListAPI.GET("/code-lists/{id}/codes", invalidCodeListID).Expect().Status(http.StatusNotFound)

		})

	})
}
