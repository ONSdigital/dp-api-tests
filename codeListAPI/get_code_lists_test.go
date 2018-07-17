package codeListAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestSuccessfullyGetASetOfCodeLists(t *testing.T) {

	var docs []*mongo.Doc

	firstCodeListDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "codelists",
		Key:        "_id",
		Value:      firstCodeListID,
		Update:     validFirstCodeListData,
	}

	secondCodeListDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "codelists",
		Key:        "_id",
		Value:      secondCodeListID,
		Update:     validSecondCodeListData,
	}

	thirdCodeListDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "codelists",
		Key:        "_id",
		Value:      thirdCodelistID,
		Update:     validThirdCodeListData,
	}

	docs = append(docs, firstCodeListDoc, secondCodeListDoc, thirdCodeListDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	codeListAPI := httpexpect.New(t, cfg.CodeListAPIURL)

	Convey("Given a set of code list exists", t, func() {
		Convey("When you request all code lists", func() {
			Convey("Then set of code lists containing information about dimensions should appear", func() {

				response := codeListAPI.GET("/code-lists").
					Expect().Status(http.StatusOK).JSON().Object()

				//checking array length is alwaysgreather than 3
				response.Value("items").Array().Length().Equal(3)
				response.Value("items").Array().Element(0).Object().ValueEqual("id", firstCodeListID)
				response.Value("items").Array().Element(0).Object().ValueEqual("name", "First Code List")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("id").Equal(firstCodeListID)
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "$")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("codes").Object().Value("id").Equal("code")
				response.Value("items").Array().Element(0).Object().Value("links").Object().Value("codes").Object().Value("href").String().Match("(.+)/code-lists/" + firstCodeListID + "/codes$")

				response.Value("items").Array().Element(1).Object().ValueEqual("id", secondCodeListID)
				response.Value("items").Array().Element(2).Object().ValueEqual("id", thirdCodelistID)

				// This functionality is not implemented yet.
				//response.Value("number_of_results").Equal(6)
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

// TODO Need to write failure tests
