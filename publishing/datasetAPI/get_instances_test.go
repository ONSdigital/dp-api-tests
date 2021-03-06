package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestSuccessfullyGetAListOfInstances(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published instance exists", t, func() {
		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstancePublished,
			Update:     validPublishedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised request to get a list of instances is received", func() {
			Convey("Then a list of instances are returned", func() {
				response := datasetAPI.GET("/instances").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Element(0).Object().Value("last_updated").NotNull()
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a `completed` and `edition-confirmed` instances exist", t, func() {
		var docs []*mongo.Doc

		completedDoc := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstanceCompleted,
			Update:     validCompletedInstanceData(ids.DatasetPublished, "2018", ids.InstanceCompleted, ids.UniqueTimestamp),
		}

		editionConfirmedDoc := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstanceEditionConfirmed,
			Update:     validEditionConfirmedInstanceData(ids.DatasetPublished, "2017", ids.InstanceEditionConfirmed, ids.UniqueTimestamp),
		}

		docs = append(docs, completedDoc, editionConfirmedDoc)

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When the authorised request contains a query parameter 'state' of value completed", func() {
			Convey("Then return only instances that contain a 'state' of value completed", func() {

				response := datasetAPI.GET("/instances").WithQuery("state", "completed").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				var foundInstance bool

				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
					response.Value("items").Array().Element(i).Object().Value("id").NotEqual(ids.InstanceEditionConfirmed)
					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstanceCompleted {
						foundInstance = true
						response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")
					}
				}

				So(foundInstance, ShouldEqual, true)
			})
		})

		Convey("When the authorised request contains a query parameter 'state' of value `completed` and `edition-confirmed`", func() {
			Convey("Then return all instances that contain a 'state' of value `completed` or `edition-confirmed`", func() {

				response := datasetAPI.GET("/instances").WithQuery("state", "completed,edition-confirmed").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				count := 0

				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstanceCompleted {
						response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")

						count++
					}
					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstanceEditionConfirmed {
						response.Value("items").Array().Element(i).Object().Value("state").Equal("edition-confirmed")

						count++
					}
				}

				// Check both resources were found in response
				So(count, ShouldEqual, 2)
			})
		})

		Convey("When the authorised request contains a query parameter 'dataset'", func() {
			Convey("Then return all instances have a 'dataset' ID of the given value", func() {

				response := datasetAPI.GET("/instances").WithQuery("dataset", ids.DatasetPublished).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				count := 0

				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
					if response.Value("items").Array().Element(i).Object().Value("links").Object().Value("dataset").Object().Value("id").String().Raw() == ids.DatasetPublished {
						count++
					}
				}

				// Check both resources were found in response
				So(count, ShouldEqual, 2)
			})
		})

		Convey("When the authorised request contains a query parameter 'dataset' and a 'state' of value `completed`", func() {
			Convey("Then return all instances have a 'dataset' ID of the given value and 'state' of value `completed`", func() {

				response := datasetAPI.GET("/instances").WithQuery("dataset", ids.DatasetPublished).WithQuery("state", "completed").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				count := 0

				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
					if response.Value("items").Array().Element(i).Object().Value("links").Object().Value("dataset").Object().Value("id").String().Raw() == ids.DatasetPublished {
						response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")

						count++
					}
				}

				// Check both resources were found in response
				So(count, ShouldEqual, 1)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetAListOfInstances(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      ids.InstancePublished,
		Update:     validPublishedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
	}

	Convey("Given an instance with state `published` exists", t, func() {
		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When no authentication header is provided in request to get list of resources", func() {
			Convey("Then return a status unauthorized (401)", func() {
				datasetAPI.GET("/instances").
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an unauthorised request is made to get resource", func() {
			Convey("Then return a status of unauthorized (401)", func() {
				datasetAPI.GET("/instances").WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an authorised request to get a list of resources is made with an invalid filter value for 'state'", func() {
			Convey("Then return a status of bad request (400) with message `Bad request - invalid filter state values`", func() {
				datasetAPI.GET("/instances").WithQuery("state", "foo").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusBadRequest).Body().Contains("bad request - invalid filter state values: [foo]")
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}
