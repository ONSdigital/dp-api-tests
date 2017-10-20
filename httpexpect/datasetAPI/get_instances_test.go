package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAListOfInstances(t *testing.T) {

	mongo.Teardown(database, "instances", "_id", instanceID)
	mongo.Setup(database, "instances", "_id", instanceID, validPublishedInstanceData)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of instances", t, func() {
		Convey("when the user is authorised", func() {
			response := datasetAPI.GET("/instances").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
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

			mongo.TeardownMany(d)

			mongo.SetupMany(d)
			response := datasetAPI.GET("/instances").WithQuery("state", "completed").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
				response.Value("items").Array().Element(i).Object().Value("id").NotEqual("779")
				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == "799" {
					response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")
				}
			}

			mongo.TeardownMany(d)
		})

		Convey("when the user filters by multiple 'state' values", func() {
			response := datasetAPI.GET("/instances").WithQuery("state", "completed,edition-confirmed").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
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
			datasetAPI.GET("/instances").WithQuery("state", "foo").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusBadRequest)
		})
	})

	mongo.Teardown(database, "instances", "_id", instanceID)
}
