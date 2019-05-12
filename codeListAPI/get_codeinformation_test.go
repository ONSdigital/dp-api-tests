package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"

)

func TestSuccessfullyGetCodeInformationAboutACode(t *testing.T) {

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a code list and codes exists", t, func() {
		Convey("When you request a specific code information", func() {
			Convey("Then that particular code information about that code should appear", func() {

				// first code information
				firstCodeResponse := codeListAPI.GET("/code-lists/{}/editions/{}/codes/{}", firstCodeListID, firstCodeListEdition, firstCodeListFirstCodeID).
					Expect().Status(http.StatusOK).JSON().Object()

				firstCodeResponse.Value("id").Equal(firstCodeListFirstCodeID)
				firstCodeResponse.Value("label").Equal(firstCodeListFirstLabel)
				firstCodeResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID)
				firstCodeResponse.Value("links").Object().Value("self").Object().Value("href").String().
					Match("(.+)/code-lists/" + firstCodeListID + "/editions/" + firstCodeListEdition + "/codes/" + firstCodeListFirstCodeID)


				// second code information
				secondCodeResponse := codeListAPI.GET("/code-lists/{}/editions/{}/codes/{}", firstCodeListID, firstCodeListEdition, firstCodeListSecondCodeID).
					Expect().Status(http.StatusOK).JSON().Object()

				secondCodeResponse.Value("id").Equal(firstCodeListSecondCodeID)
				secondCodeResponse.Value("label").Equal(firstCodeListSecondLabel)
				secondCodeResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID)
				secondCodeResponse.Value("links").Object().Value("self").Object().Value("href").String().
					Match("(.+)/code-lists/" + firstCodeListID + "/editions/" + firstCodeListEdition + "/codes/" + firstCodeListSecondCodeID)


				// third code information
				thirdCodeResponse := codeListAPI.GET("/code-lists/{}/editions/{}/codes/{}", firstCodeListID, firstCodeListEdition, firstCodeListThirdCodeID).
					Expect().Status(http.StatusOK).JSON().Object()

				thirdCodeResponse.Value("id").Equal(firstCodeListThirdCodeID)
				thirdCodeResponse.Value("label").Equal(firstCodeListThirdLabel)
				thirdCodeResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID)
				thirdCodeResponse.Value("links").Object().Value("self").Object().Value("href").String().
					Match("(.+)/code-lists/" + firstCodeListID + "/editions/" + firstCodeListEdition + "/codes/" + firstCodeListThirdCodeID)

			})
		})
	})

}

func TestFailureToGetInDepthInformationAboutACode(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given a code list and codes exists", t, func() {
		Convey("When you pass a code list that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}/editions/{}/codes/{code_id}", invalidCodeListID, firstCodeListEdition, firstCode).
					Expect().Status(http.StatusNotFound)
			})
		})

		Convey("When you pass a code that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}/editions/{}/codes/{code_id}", firstCodeListID, firstCodeListEdition, invalidCode).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}
