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

	if err := db.setUp(); err != nil {
		log.ErrorC("test setup failure", err, nil)
		os.Exit(1)
	}

	Convey("given code lists exist", t, func() {
		Convey("when a valid request is sent to get CodeLists", func() {
			codeListAPI := httpexpect.New(t, db.cfg.CodeListAPIURL)

			Convey("then the expected response is received", func() {
				response := codeListAPI.GET("/code-lists").Expect().Status(http.StatusOK).JSON().Object()
				response.Value("items").Array().Length().Equal(1)
				response.Value("items").Array().Element(0).Object().ValueEqual("name", codeListName)

				links := response.Value("items").Array().Element(0).Object().Value("links")

				links.Object().Value("self").Object().Value("id").Equal(codeListID)
				links.Object().Value("self").Object().Value("href").String().Match("(.*)\\/code-lists\\/" + codeListID + "$")

				links.Object().Value("editions").Object().Value("href").String().Match("(.*)\\/code-lists\\/" + codeListID + "\\/editions" + "$")

				links.Object().Value("latest").Object().Value("id").Equal(edition)
				links.Object().Value("latest").Object().Value("href").String().Match("(.*)\\/code-lists\\/" + codeListID + "\\/editions\\/" + edition + "$")
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
