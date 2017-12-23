package codeListAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetACodeList(t *testing.T) {

	secondCodeList := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      secondCodeListID,
		Update:     validSecondCodeListData,
	}

	if err := mongo.Teardown(secondCodeList); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup(secondCodeList); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a code list exists", t, func() {
		Convey("When you request a specific code list with id", func() {
			Convey("Then the code list data should appear", func() {

				response := codeListAPI.GET("/code-lists/{id}", secondCodeListID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("id").Equal(secondCodeListID)
				response.Value("name").Equal("Second Code List")
				response.Value("links").Object().Value("self").Object().Value("id").Equal(secondCodeListID)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + secondCodeListID + "$")

				response.Value("links").Object().Value("codes").Object().Value("id").Equal("code")
				response.Value("links").Object().Value("codes").Object().Value("href").String().Match("(.+)/code-lists/" + secondCodeListID + "/codes$")

			})
		})
	})

	if err := mongo.Teardown(secondCodeList); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToGetACodeList(t *testing.T) {

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given a code list exists", t, func() {
		Convey("When you pass a code list that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}", invalidCodeListID).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}
