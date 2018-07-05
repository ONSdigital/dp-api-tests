package codeListAPI

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulRetrievalOfCodeLists(t *testing.T) {

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := ds.SetupCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

	codelistAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given codelists are available for retrieval", t, func() {
		Convey("When a request is made to retrieve the codelists", func() {
			Convey("Then the response contains a list of codelists", func() {

				response := codelistAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("number_of_results").Equal(3)
				response.Value("count").Equal(3)
				response.Value("limit").Equal(3)
				response.Value("offset").Equal(0)

				item1 := response.Value("items").Array().Element(0).Object()
				item1.Value("name").Equal("Tottenham")
				item1.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ABCDEF")
				item1.Value("links").Object().Value("self").Object().Value("id").Equal("ABCDEF")
				item1.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ABCDEF/editions")

				item2 := response.Value("items").Array().Element(1).Object()
				item2.Value("name").Equal("Crystal Palace")
				item2.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ZYXWVU")
				item2.Value("links").Object().Value("self").Object().Value("id").Equal("ZYXWVU")
				item2.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ZYXWVU/editions")

				item3 := response.Value("items").Array().Element(2).Object()
				item3.Value("name").Equal("England")
				item3.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ENG")
				item3.Value("links").Object().Value("self").Object().Value("id").Equal("ENG")
				item3.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ENG/editions")

			})
		})
	})

	Convey("Given codelists are available for retrieval including a geography codelist", t, func() {
		Convey("When a request is made to retrieve the codelists filtered by geography", func() {
			Convey("Then the response contains only geography codelists", func() {

				response := codelistAPI.GET("/code-lists").WithQuery("type", "geography").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("number_of_results").Equal(1)
				response.Value("count").Equal(1)
				response.Value("limit").Equal(1)
				response.Value("offset").Equal(0)

				item := response.Value("items").Array().Element(0).Object()
				item.Value("name").Equal("England")
				item.Value("links").Object().Value("self").Object().Value("href").Equal("/code-lists/ENG")
				item.Value("links").Object().Value("self").Object().Value("id").Equal("ENG")
				item.Value("links").Object().Value("editions").Object().Value("href").Equal("/code-lists/ENG/editions")

			})
		})
	})

	if err := ds.DropCodelists(); err != nil {
		t.Error(err)
		t.Fail()
	}

}
