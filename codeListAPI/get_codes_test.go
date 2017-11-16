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

func TestSuccessfullyGetAListOfAllCodesWithinCodeList(t *testing.T) {

	var docs []mongo.Doc

	firstCodeListDoc := mongo.Doc{
		Database:   "codelists",
		Collection: "codelists",
		Key:        "_id",
		Value:      firstCodeListID,
		Update:     validFirstCodeListData,
	}

	firstCodeListCodesDoc := mongo.Doc{
		Database:   "codelists",
		Collection: "codes",
		Key:        "_id",
		Value:      firstCodeListFirstCodeID,
		Update:     validFirstCodeListFirstCodeData,
	}

	secondCodeListCodesDoc := mongo.Doc{
		Database:   "codelists",
		Collection: "codes",
		Key:        "_id",
		Value:      firstCodeListSecondCodeID,
		Update:     validFirstCodeListSecondCodeData,
	}

	thirdCodeListCodesDoc := mongo.Doc{
		Database:   "codelists",
		Collection: "codes",
		Key:        "_id",
		Value:      firstCodeListThirdCodeID,
		Update:     validFirstCodeListThirdCodeData,
	}
	docs = append(docs, firstCodeListDoc, firstCodeListCodesDoc, secondCodeListCodesDoc, thirdCodeListCodesDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.SetupMany(d); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a code list and codes exists", t, func() {
		Convey("When you request a list of all codes", func() {
			Convey("Then the list of codes within a code list should appear", func() {

				response := codeListAPI.GET("/code-lists/{id}/codes", firstCodeListID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Element(0).Object().Value("id").Equal(firstCode)
				response.Value("items").Array().Element(0).Object().Value("label").Equal("First Code List label one")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().ValueEqual("id", firstCodeListID)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")

				response.Value("items").Array().Element(1).Object().Value("id").Equal(secondCode)
				response.Value("items").Array().Element(1).Object().Value("label").Equal("First Code List label two")
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().ValueEqual("id", firstCodeListID)
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")
				response.Value("items").Array().Element(1).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")

				response.Value("items").Array().Element(2).Object().Value("id").Equal(thirdCode)
				response.Value("items").Array().Element(2).Object().Value("label").Equal("First Code List label three")
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("code_list").Object().ValueEqual("id", firstCodeListID)
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")
				response.Value("items").Array().Element(2).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")

			})
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

// This test fails.
// Bug Raised-When you given an invalid codelist id, its returning 200 instead of 400.
// 404 - Code list not found
func TestFailureToGetAListOfAllCodesWithinCodeList(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a code list exists", t, func() {
		Convey("When you pass a code list that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}/codes", invalidCodeListID).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}
