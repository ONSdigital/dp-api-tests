package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulRetrievalOfEditions(t *testing.T) {

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := ds.SetupCodelistEditions(); err != nil {
		t.Error(err)
		t.Fail()
	}

	codelistAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given two editions of the same codelist exists", t, func() {
		Convey("When a request is made for the editions", func() {
			Convey("Then the codelists are returned to the user", func() {
				response := codelistAPI.GET("/code-lists/ABCDEF/editions").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("number_of_results").Equal(2)

				item1 := response.Value("items").Array().Element(0).Object()
				item1.Value("id").Equal("ABCDEF")
				item1.Value("edition").Equal("2018")
				item1.Value("label").Equal("Tottenham")
				item1.Value("release_date").Equal("01 Jan 2018")
				item1.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ABCDEF/editions/2018")
				item1.Value("links").Object().Value("self").Object().Value("id").Equal("2018")
				item1.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ABCDEF/editions")
				item1.Value("links").Object().Value("codes").Object().Value("href").Equal("/code-lists/ABCDEF/editions/2018/codes")

				item2 := response.Value("items").Array().Element(1).Object()
				item2.Value("id").Equal("ABCDEF")
				item2.Value("edition").Equal("2017")
				item2.Value("label").Equal("Tottenham")
				item2.Value("release_date").Equal("01 Jan 2017")
				item2.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ABCDEF/editions/2017")
				item2.Value("links").Object().Value("self").Object().Value("id").Equal("2017")
				item2.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ABCDEF/editions")
				item2.Value("links").Object().Value("codes").Object().Value("href").Equal("/code-lists/ABCDEF/editions/2017/codes")

			})
		})
	})

	Convey("Given a user wants to know the editions of a code list that doesn't exist", t, func() {
		Convey("When a request is made for the editions", func() {
			Convey("Then a not found status is returned to the user", func() {

				codelistAPI.GET("/code-lists/BCDSOE/editions").Expect().Status(http.StatusNotFound)

			})
		})
	})

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

}
