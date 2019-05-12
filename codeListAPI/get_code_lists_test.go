package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"

)

func TestSuccessfullyGetASetOfCodeLists(t *testing.T) {

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a set of code list exists", t, func() {
		Convey("When you request all code lists", func() {
			Convey("Then set of code lists containing information about dimensions should appear", func() {

				response := codeListAPI.GET("/code-lists").
					Expect().Status(http.StatusOK).JSON().Object()

				//checking array length is alwaysgreather than 3
				response.Value("items").Array().Length().Equal(2)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("id").Equal(firstCodeListID)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("editions").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/editions$")

				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("id").Equal(secondCodeListID)

				// This functionality is not implemented yet.
				//response.Value("number_of_results").Equal(6)
			})
		})
	})
	
}

// TODO Need to write failure tests
