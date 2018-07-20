package codeListAPI

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/gavv/httpexpect"
	"net/http"
	"github.com/ONSdigital/dp-bolt/bolt"
	"github.com/ONSdigital/dp-api-tests/config"
)

func TestGetCodesSuccess(t *testing.T) {
	db := NewDB(t)
	defer db.bolt.Close()

	Convey("given codes exist for the requested codeList edition", t, func() {
		db.Setup(AllTestData()...)

		Convey("when a valid getCode request is sent", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)
			response := codeListAPI.GET("/code-lists/fender-guitars/editions/2018/codes").Expect().Status(http.StatusOK).JSON().Object()

			Convey("then the expected response is returned", func() {
				response.Value("number_of_results").Equal(3)
				response.Value("items").Array().Length().Equal(3)

				response.Value("items").Array().Element(0).Object().Value("id").Equal("Jazzmaster")
				response.Value("items").Array().Element(0).Object().Value("label").Equal("Jazzmaster")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match(fenderGuitars2018.codeListLink)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(fenderGuitars2018.CodeLink("Jazzmaster"))

				response.Value("items").Array().Element(1).Object().Value("id").Equal("Tele")
				response.Value("items").Array().Element(1).Object().Value("label").Equal("Telecaster")
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match(fenderGuitars2018.codeListLink)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(fenderGuitars2018.CodeLink("Tele"))

				response.Value("items").Array().Element(2).Object().Value("id").Equal("Strat")
				response.Value("items").Array().Element(2).Object().Value("label").Equal("Stratocaster")
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match(fenderGuitars2018.codeListLink)
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("self").Object().Value("href").String().Match(fenderGuitars2018.CodeLink("Strat"))
			})
		})

		Reset(func() {
			db.TearDown()
		})
	})
}

func TestGetCodesNotFound(t *testing.T) {
	db := NewDB(t)
	defer db.bolt.Close()
	db.Setup(bolt.Stmt{Query: "CREATE (node:`_code_list`:`_code_list_gibson-guitars`:`_api_test` { label:'gibson', edition: {ed}})", Params: map[string]interface{}{"ed": edition2018}})

	Convey("given no codes exist for the requested codeList edition", t, func() {
		// create just the codelist with no codes.

		Convey("when a getCode request is sent", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)
			response := codeListAPI.GET("/code-lists/gibson-guitars/editions/2018/codes").Expect()

			Convey("the a 404 status is returned", func() {
				response.Status(http.StatusNotFound).Text().Equal("codes not found\n")
			})
		})

		Reset(func() {
			db.TearDown()
		})
	})
}

func TestGetCodesEditionNotFound(t *testing.T) {
	Convey("given the requested edition does not exist", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		// create just the codelist with no codes.

		Convey("when a getCode request is sent", func() {
			codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)
			response := codeListAPI.GET("/code-lists/gibson-guitars/editions/2018/codes").Expect()

			Convey("the a 404 status is returned", func() {
				response.Status(http.StatusNotFound).Text().Equal("edition not found\n")
			})
		})
	})
}

func TestGetCodeSuccess(t *testing.T) {
	db := NewDB(t)
	defer db.bolt.Close()

	Convey("given a code exists for the requested codeList edition", t, func() {
		db.Setup(AllTestData()...)

		Convey("when a valid getCode request is sent", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)
			response := codeListAPI.GET("/code-lists/fender-guitars/editions/2018/codes/Tele").Expect().Status(http.StatusOK).JSON().Object()

			Convey("then the expected response is returned", func() {
				response.Value("id").Equal("Tele")
				response.Value("label").Equal("Telecaster")
				response.Value("links").Object().Value("code_list").Object().Value("href").String().Match(fenderGuitars2018.codeListLink)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(fenderGuitars2018.CodeLink("Tele"))
			})
		})

		Reset(func() {
			db.TearDown()
		})
	})
}

func TestGetCodeNotFound(t *testing.T) {
	db := NewDB(t)
	defer db.bolt.Close()

	Convey("given a code exists for the requested codeList edition", t, func() {
		db.Setup(AllTestData()...)

		Convey("when a getCode request is sent", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)
			response := codeListAPI.GET("/code-lists/fender-guitars/editions/2018/codes/BATMAN").Expect()

			Convey("then the expected response is returned", func() {
				response.Status(http.StatusNotFound).Text().Equal("code not found\n")
			})
		})

		Reset(func() {
			db.TearDown()
		})
	})
}
