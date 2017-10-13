package codeListAPI

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Get in depth information about a code within a code list
// 200 - Get in depth information about a code
func TestGetCodeInformationWithInACodeList_ReturnsListOfCodes(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a set of code lists", t, func() {

		response := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()

		cpiCodeID := response.Value("items").Array().Element(0).Object().Value("id").String().Raw()
		timeID := response.Value("items").Array().Element(1).Object().Value("id").String().Raw()

		geographyID := response.Value("items").Array().Element(2).Object().Value("id").String().Raw()

		Convey("Then get the list of all code with in a CPI code", func() {
			cpiCodesResponse := codeListAPI.GET("/code-lists/{id}/codes", cpiCodeID).Expect().Status(http.StatusOK).JSON().Object()

			cpiFourthCode := cpiCodesResponse.Value("items").Array().Element(3).Object().Value("id").String().Raw()

			Convey("Then get the list of all code with in a CPI code", func() {

				codeInfoResponse := codeListAPI.GET("/code-lists/{id}/codes/{code_id}", cpiCodeID, cpiFourthCode).Expect().Status(http.StatusOK).JSON().Object()

				codeInfoResponse.Value("id").Equal(cpiFourthCode)
				codeInfoResponse.Value("label").Equal("Education")
				codeInfoResponse.Value("links").Object().Value("code_list").Object().Value("id").Equal(cpiCodeID)
				codeInfoResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Contains(cpiCodeID)

				codeInfoResponse.Value("links").Object().Value("self").Object().Value("href").String().Contains(cpiCodeID).Contains("codes").Contains(cpiFourthCode)

			})

		})

		Convey("And get the list of all codes with in a time code", func() {
			timeListResponse := codeListAPI.GET("/code-lists/{id}/codes", timeID).Expect().Status(http.StatusOK).JSON().Object()

			timeFifthCode := timeListResponse.Value("items").Array().Element(4).Object().Value("id").String().Raw()
			Convey("Then get the list of all code with in a CPI code", func() {

				timeInfoResponse := codeListAPI.GET("/code-lists/{id}/codes/{code_id}", timeID, timeFifthCode).Expect().Status(http.StatusOK).JSON().Object()

				timeInfoResponse.Value("id").Equal(timeFifthCode)

				timeInfoResponse.Value("label").Equal(timeFifthCode)

				timeInfoResponse.Value("links").Object().Value("code_list").Object().Value("id").Equal(timeID)

				timeInfoResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Contains(timeID)

				timeInfoResponse.Value("links").Object().Value("self").Object().Value("href").String().Contains(timeID).Contains("codes").Contains(timeFifthCode)

			})
		})

		Convey("And get the list of all codes with in a geography code", func() {
			geographyListResponse := codeListAPI.GET("/code-lists/{id}/codes", geographyID).Expect().Status(http.StatusOK).JSON().Object()

			geographyListResponse.Value("items").Null()

		})

	})
}

// Bug Raised
// 400 - Code List not found/Invalid code list
func TestGetCodeInformationWithInACodeList_InvalidCodeList(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a set of code lists", t, func() {

		codeListID := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object().Value("items").Array().Element(0).Object().Value("id").String().Raw()

		Convey("Then get the list of all code with in a CPI code", func() {

			cpiCodesResponse := codeListAPI.GET("/code-lists/{id}/codes", codeListID).Expect().Status(http.StatusOK).JSON().Object()

			cpiFourthCode := cpiCodesResponse.Value("items").Array().Element(3).Object().Value("id").String().Raw()

			Convey("With an invalid code list throws 404 error", func() {
				invalidCodeListID := strings.Replace(codeListID, "-", "", 9)

				codeListAPI.GET("/code-lists/{id}/codes/{code_id}", invalidCodeListID, cpiFourthCode).Expect().Status(http.StatusNotFound)

			})
		})
	})
}

// Bug Raised
// 400 - Code not found/Invalid code
func TestGetCodeInformationWithInACodeList_InvalidCode(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a set of code lists", t, func() {

		codeListID := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object().Value("items").Array().Element(0).Object().Value("id").String().Raw()

		Convey("Then get the list of all code with in a CPI code", func() {

			cpiCodesResponse := codeListAPI.GET("/code-lists/{id}/codes", codeListID).Expect().Status(http.StatusOK).JSON().Object()

			cpiFourthCode := cpiCodesResponse.Value("items").Array().Element(3).Object().Value("id").String().Raw()

			Convey("With an invalid code throws 404 error", func() {
				invalidCodeID := strings.Replace(cpiFourthCode, "_", "", 5)

				fmt.Println(invalidCodeID)
				codeListAPI.GET("/code-lists/{id}/codes/{code_id}", codeListID, invalidCodeID).Expect().Status(http.StatusNotFound)

			})
		})
	})
}
