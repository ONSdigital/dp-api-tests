package codeListAPI

import (
	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/gavv/httpexpect"
	"net/http"
)

func Test_GetCodeListsSuccess(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.bolt.Close()

	if err := db.setUp(allTestData()...); err != nil {
		log.ErrorC("test setup failure", err, nil)
		os.Exit(1)
	}

	Convey("given code lists exist", t, func() {
		Convey("when a valid request is sent to get CodeLists", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)

			response := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()

			Convey("then a 200 response status is returned", func() {
				response.Value("number_of_results").Equal(2)
				response.Value("items").Array().Length().Equal(2)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("id").Equal(gibsonGuitars2017.codeList)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(gibsonGuitars2017.codeListLink)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("editions").Object().Value("href").String().Match(gibsonGuitars2017.editionsLink)

				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("id").Equal(fenderGuitars2018.codeList)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(fenderGuitars2018.codeListLink)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("editions").Object().Value("href").String().Match(fenderGuitars2018.editionsLink)
			})
		})

		Reset(func() {
			db.tearDown()
		})
	})
}

func Test_GetCodeListsNotFound(t *testing.T) {
	cfg, _ := config.Get()
	Convey("given non code lists exist", t, func() {
		Convey("when a valid request is sent to get CodeLists", func() {
			codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

			Convey("then a 404 response is returned", func() {
				codeListAPI.GET("/code-lists").Expect().Status(http.StatusNotFound).Text().Equal("resource not found\n")
			})
		})
	})
}

func Test_GetCodeListSuccess(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.bolt.Close()

	if err := db.setUp(allTestData()...); err != nil {
		log.ErrorC("test setup failure", err, nil)
		os.Exit(1)
	}

	Convey("given a code lists exist", t, func() {
		Convey("when a valid request is sent to get CodeList", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)

			Convey("then the expected response is received", func() {
				response := codeListAPI.GET("/code-lists/" + gibsonGuitars2017.codeList).Expect().Status(http.StatusOK).JSON().Object()
				response.Value("links").Object().Value("self").Object().Value("id").Equal(gibsonGuitars2017.codeList)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(gibsonGuitars2017.codeListLink)
				response.Value("links").Object().Value("editions").Object().Value("href").String().Match(gibsonGuitars2017.editionsLink)
			})
		})

		Reset(func() {
			db.tearDown()
		})
	})
}

func Test_GetCodeListNotFound(t *testing.T) {
	Convey("given the request code lists does not exist", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		Convey("when a request is sent to get CodeList", func() {
			codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

			Convey("then a 404 response is returned", func() {
				codeListAPI.GET("/code-lists/fender-guitars").Expect().Status(http.StatusNotFound).Text().Equal("resource not found\n")
			})
		})
	})
}
