package codeListAPI

import (
	"os"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ONSdigital/dp-api-tests/config"
)

func Test_GetEditionsSuccess(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.bolt.Close()

	if err := db.setUp(allTestData()...); err != nil {
		log.ErrorC("test setup failure", err, nil)
		os.Exit(1)
	}

	Convey("given 2 editions of a code list exist", t, func() {

		Convey("when a valid getEditions request is sent", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)

			Convey("then the expected response is received", func() {
				response := codeListAPI.GET("/code-lists/gibson-guitars/editions").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Length().Equal(2)
				response.Value("items").Array().Element(0).Object().ValueEqual("edition", gibsonGuitars2017.edition)
				response.Value("items").Array().Element(0).Object().ValueEqual("label", gibsonGuitars2017.label)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("id").Equal(gibsonGuitars2017.id)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(gibsonGuitars2017.editionLink)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("editions").Object().Value("href").String().Match(gibsonGuitars2017.editionsLink)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("codes").Object().Value("href").String().Match(gibsonGuitars2017.codesLink)

				response.Value("items").Array().Element(1).Object().ValueEqual("edition", gibsonGuitars2018.edition)
				response.Value("items").Array().Element(1).Object().ValueEqual("label", gibsonGuitars2018.label)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("id").Equal(gibsonGuitars2018.id)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(gibsonGuitars2018.editionLink)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("editions").Object().Value("href").String().Match(gibsonGuitars2018.editionsLink)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("codes").Object().Value("href").String().Match(gibsonGuitars2018.codesLink)
			})
		})

		Reset(func() {
			db.tearDown()
		})
	})
}

func Test_GetEditionsNotFound(t *testing.T) {
	Convey("given no editions of a code list exist", t, func() {
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

func Test_GetEditionSuccess(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.bolt.Close()

	if err := db.setUp(allTestData()...); err != nil {
		log.ErrorC("test setup failure", err, nil)
		os.Exit(1)
	}

	Convey("given the requested edition of a code list exists", t, func() {

		Convey("when a valid getEdition request is sent", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)

			Convey("then the expected response is received", func() {
				response := codeListAPI.GET("/code-lists/fender-guitars/editions/2017").Expect().Status(http.StatusOK).JSON().Object()
				response.Value("edition").Equal(fenderGuitars2017.edition)
				response.Value("label").Equal(fenderGuitars2017.label)
				response.Value("links").Object().Value("self").Object().Value("id").Equal(fenderGuitars2017.edition)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(fenderGuitars2017.editionLink)
				response.Value("links").Object().Value("editions").Object().Value("href").String().Match(fenderGuitars2017.editionsLink)
				response.Value("links").Object().Value("codes").Object().Value("href").String().Match(fenderGuitars2017.codesLink)
			})
		})

		Reset(func() {
			db.tearDown()
		})
	})
}

func Test_GetEditionNotFound(t *testing.T) {
	Convey("given the requested edition does not exist", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		Convey("when a getEdition request is sent", func() {
			codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

			Convey("then a 404 response is returned", func() {
				codeListAPI.GET("/code-lists/fender-guitars/editions/2020").Expect().Status(http.StatusNotFound).Text().Equal("resource not found\n")
			})
		})
	})
}

func Test_GetEditionInternalServerError(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.bolt.Close()

	// Add the same data twice - to trigger a non unique result which should result in an internal server error.
	testData := append(fender2017, fender2017...)

	if err := db.setUp(testData...); err != nil {
		log.ErrorC("test setup failure", err, nil)
		os.Exit(1)
	}

	Convey("given the requested edition is not unique", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		Convey("when a getEdition request is sent", func() {
			codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

			Convey("then a 500 response is returned", func() {
				codeListAPI.GET("/code-lists/fender-guitars/editions/2017").Expect().Status(http.StatusInternalServerError).Text().Equal("internal server error\n")
			})
		})

		Reset(func() {
			db.tearDown()
		})
	})
}
