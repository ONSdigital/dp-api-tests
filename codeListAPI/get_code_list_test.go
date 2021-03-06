package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"

)

func TestSuccessfullyGetACodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a code list exists", t, func() {
		Convey("When you request a specific code list with id", func() {
			Convey("Then the code list data should appear", func() {

				response := codeListAPI.GET("/code-lists/{id}", secondCodeListID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("links").Object().Value("self").Object().Value("id").Equal(secondCodeListID)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + secondCodeListID + "$")
				response.Value("links").Object().Value("editions").Object().Value("href").String().Match("(.+)/code-lists/" + secondCodeListID + "/editions$")

			})
		})
	})

}

func TestFailureToGetACodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given a code list exists", t, func() {
		Convey("When you pass a code list that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}", invalidCodeListID).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}
