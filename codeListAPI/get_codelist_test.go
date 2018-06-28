package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulRetrievalOfCodeList(t *testing.T) {

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := ds.SetupCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

	codelistAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a codelist available for retrieval", t, func() {
		Convey("When a request is made to retrieve the codelist", func() {
			Convey("Then the response contains the expected values", func() {

				response := codelistAPI.GET("/code-lists/ENG").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("name").Equal("England")
				response.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ENG")
				response.Value("links").Object().Value("self").Object().Value("id").Equal("ENG")
				response.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ENG/editions")
				response.Value("links").Object().Value("latest").Object().Value("href").Equal("/code-lists/ENG/editions/2006")
				response.Value("links").Object().Value("latest").Object().Value("id").Equal("2006")

			})
		})
	})

	Convey("Given a codelist available for retrieval", t, func() {
		Convey("When a request is made for a codelist that does not exist", func() {
			Convey("Then the response status is 404", func() {

				codelistAPI.GET("/code-lists/WAL").Expect().Status(http.StatusNotFound)

			})
		})
	})

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

}
