package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"

)

func TestSuccessfullyGetAListOfAllCodesWithinCodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a code list and codes exists", t, func() {
		Convey("When you request a list of all codes", func() {
			Convey("Then the list of codes within a code list should appear", func() {

				response := codeListAPI.GET("/code-lists/{id}/editions/one-off/codes", firstCodeListID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Element(0).Object().Value("id").Equal(firstCode)
				response.Value("items").Array().Element(0).Object().Value("label").Equal(firstCodeListFirstLabel)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/editions/")

				response.Value("items").Array().Element(1).Object().Value("id").Equal(secondCode)
				response.Value("items").Array().Element(1).Object().Value("label").Equal(firstCodeListSecondLabel)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/editions/")

				response.Value("items").Array().Element(2).Object().Value("id").Equal(thirdCode)
				response.Value("items").Array().Element(2).Object().Value("label").Equal(firstCodeListThirdLabel)
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/editions/")

			})
		})
	})

}

func TestFailureToGetAListOfAllCodesWithinCodeList(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given a code list exists", t, func() {
		Convey("When you pass a code list that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}/codes", invalidCodeListID).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}
