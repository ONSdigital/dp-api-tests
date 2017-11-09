package datasetAPI

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

func TestGetAListOfInstances(t *testing.T) {

	if err := mongo.Teardown(database, "instances", "_id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup(database, "instances", "_id", instanceID, validPublishedInstanceData); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of instances", t, func() {
		Convey("when the user is authorised", func() {
			response := datasetAPI.GET("/instances").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Element(0).Object().Value("id").NotNull()
		})

		Convey("when the user filters by a 'state' value", func() {
			var docs []mongo.Doc

			completedDoc := mongo.Doc{
				Database:   database,
				Collection: "instances",
				Key:        "_id",
				Value:      "799",
				Update:     validCompletedInstanceData,
			}
			editionConfirmedDoc := mongo.Doc{
				Database:   database,
				Collection: "instances",
				Key:        "_id",
				Value:      "779",
				Update:     validEditionConfirmedInstanceData,
			}

			docs = append(docs, completedDoc, editionConfirmedDoc)

			d := &mongo.ManyDocs{
				Docs: docs,
			}

			if err := mongo.TeardownMany(d); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}
			}

			if err := mongo.SetupMany(d); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			response := datasetAPI.GET("/instances").WithQuery("state", "completed").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
				response.Value("items").Array().Element(i).Object().Value("id").NotEqual("779")
				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == "799" {
					response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")
				}
			}

			if err := mongo.TeardownMany(d); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}
			}
		})

		Convey("when the user filters by multiple 'state' values", func() {
			response := datasetAPI.GET("/instances").WithQuery("state", "completed,edition-confirmed").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == "799" {
					response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")
				}
				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == "779" {
					response.Value("items").Array().Element(i).Object().Value("state").Equal("edition-confirmed")
				}
			}
		})
	})

	Convey("Fail to get a list of instances", t, func() {
		Convey("When the user is unauthorised", func() {
			datasetAPI.GET("/instances").
				Expect().Status(http.StatusUnauthorized)
		})

		Convey("When the user filters by the wrong 'state' value", func() {
			datasetAPI.GET("/instances").WithQuery("state", "foo").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusBadRequest)
		})
	})

	if err := mongo.Teardown(database, "instances", "_id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}
