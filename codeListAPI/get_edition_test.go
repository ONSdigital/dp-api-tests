package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulRetrievalOfEdition(t *testing.T) {

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := ds.SetupCodelistEditions(); err != nil {
		t.Error(err)
		t.Fail()
	}

	codelistAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given an edition of a codelist eqists", t, func() {
		Convey("When a request is made for the edition", func() {
			Convey("Then the codelist is returned to the user", func() {
				response := codelistAPI.GET("/code-lists/ABCDEF/editions/2018").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("id").Equal("ABCDEF")
				response.Value("edition").Equal("2018")
				response.Value("label").Equal("Tottenham")
				response.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ABCDEF/editions/2018")
				response.Value("links").Object().Value("self").Object().Value("id").Equal("2018")
				response.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ABCDEF/editions")
				response.Value("links").Object().Value("codes").Object().Value("href").Equal("/code-lists/ABCDEF/editions/2018/codes")

			})
		})
	})

	Convey("Given a user wants to know the edition of a code list that doesn't exist", t, func() {
		Convey("When a request is made for the edition", func() {
			Convey("Then the a not found status is returned to the user", func() {

				codelistAPI.GET("/code-lists/BCDSOE/editions/2018").Expect().Status(http.StatusNotFound)

				codelistAPI.GET("/code-lists/ABCDEF/editions/2019").Expect().Status(http.StatusNotFound)

			})
		})
	})

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

}
