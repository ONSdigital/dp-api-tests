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

func TestSuccessfullyGetCodeInformationAboutACode(t *testing.T) {

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
		Convey("When you request a specific code information", func() {
			Convey("Then that particular code information about that code should appear", func() {

				// first code information
				firstCodeResponse := codeListAPI.GET("/code-lists/{id}/codes/{code_id}", firstCodeListID, firstCode).
					Expect().Status(http.StatusOK).JSON().Object()

				firstCodeResponse.Value("id").Equal(firstCode)
				firstCodeResponse.Value("label").Equal("First Code List label one")
				firstCodeResponse.Value("links").Object().Value("code_list").Object().ValueEqual("id", firstCodeListID)
				firstCodeResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")
				firstCodeResponse.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")

				// second code information
				secondCodeResponse := codeListAPI.GET("/code-lists/{id}/codes/{code_id}", firstCodeListID, secondCode).
					Expect().Status(http.StatusOK).JSON().Object()

				secondCodeResponse.Value("id").Equal(secondCode)
				secondCodeResponse.Value("label").Equal("First Code List label two")
				secondCodeResponse.Value("links").Object().Value("code_list").Object().ValueEqual("id", firstCodeListID)
				secondCodeResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")
				secondCodeResponse.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")

				// third code information
				thirdCodeResponse := codeListAPI.GET("/code-lists/{id}/codes/{code_id}", firstCodeListID, thirdCode).
					Expect().Status(http.StatusOK).JSON().Object()

				thirdCodeResponse.Value("id").Equal(thirdCode)
				thirdCodeResponse.Value("label").Equal("First Code List label three")
				thirdCodeResponse.Value("links").Object().Value("code_list").Object().ValueEqual("id", firstCodeListID)
				thirdCodeResponse.Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")
				thirdCodeResponse.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")
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

func TestFailureToGetInDepthInformationAboutACode(t *testing.T) {
	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	// TODO Dont skip test once endpoint has been refactored
	SkipConvey("Given a code list and codes exists", t, func() {
		Convey("When you pass a code list that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}/codes/{code_id}", invalidCodeListID, firstCode).
					Expect().Status(http.StatusNotFound)
			})
		})

		Convey("When you pass a code that does not exist", func() {
			Convey("Then the response should be status not found (404)", func() {
				codeListAPI.GET("/code-lists/{id}/codes/{code_id}", firstCodeListID, invalidCode).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}