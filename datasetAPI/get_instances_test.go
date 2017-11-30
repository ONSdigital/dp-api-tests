package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetAListOfInstances(t *testing.T) {

	instanceID := uuid.NewV4().String()
	secondInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published instance exists", t, func() {
		if err := mongo.Setup(database, "instances", "_id", instanceID, validPublishedInstanceData(datasetID, edition, instanceID)); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised request to get a list of instances is received", func() {
			Convey("Then a list of instances are returned", func() {
				response := datasetAPI.GET("/instances").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Element(0).Object().Value("id").NotNull()
			})
		})

		if err := mongo.Teardown(database, "instances", "_id", instanceID); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a `completed` and `edition-confirmed` instances exist", t, func() {
		var docs []mongo.Doc

		completedDoc := mongo.Doc{
			Database:   database,
			Collection: "instances",
			Key:        "_id",
			Value:      instanceID,
			Update:     validCompletedInstanceData(datasetID, "2018", instanceID),
		}
		editionConfirmedDoc := mongo.Doc{
			Database:   database,
			Collection: "instances",
			Key:        "_id",
			Value:      secondInstanceID,
			Update:     validEditionConfirmedInstanceData(datasetID, "2017", secondInstanceID),
		}

		docs = append(docs, completedDoc, editionConfirmedDoc)

		d := &mongo.ManyDocs{
			Docs: docs,
		}

		if err := mongo.SetupMany(d); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When the authorised request contains a query parameter 'state' of value completed", func() {
			Convey("Then return only instances that contain a 'state' of value completed", func() {

				response := datasetAPI.GET("/instances").WithQuery("state", "completed").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusOK).JSON().Object()

				var foundInstance bool

				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
					response.Value("items").Array().Element(i).Object().Value("id").NotEqual(secondInstanceID)
					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == instanceID {
						foundInstance = true
						response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")
					}
				}

				So(foundInstance, ShouldEqual, true)
			})
		})

		Convey("When the authorised request contains a query parameter 'state' of value `completed` and `edition-confirmed`", func() {
			Convey("Then return all instances that contain a 'state' of value `completed` or `edition-confirmed`", func() {

				response := datasetAPI.GET("/instances").WithQuery("state", "completed,edition-confirmed").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusOK).JSON().Object()

				count := 0

				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == instanceID {
						response.Value("items").Array().Element(i).Object().Value("state").Equal("completed")

						count++
					}
					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == secondInstanceID {
						response.Value("items").Array().Element(i).Object().Value("state").Equal("edition-confirmed")

						count++
					}
				}

				// Check both resources were found in response
				So(count, ShouldEqual, 2)
			})
		})

		if err := mongo.TeardownMany(d); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetAListOfInstances(t *testing.T) {

	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance with state `published` exists", t, func() {
		if err := mongo.Setup(database, "instances", "_id", instanceID, validPublishedInstanceData(datasetID, edition, instanceID)); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When no authentication header is provided in request to get list of resources", func() {
			Convey("Then return a status of unauthorised (401) with message `No authentication header provided`", func() {
				datasetAPI.GET("/instances").Expect().Status(http.StatusUnauthorized).
					Body().Contains("No authentication header provided\n")
			})
		})

		Convey("When an unauthorised request is made to get resource", func() {
			Convey("Then return a status of unauthorised (401) with message `Unauthorised access to API`", func() {
				datasetAPI.GET("/instances").WithHeader(internalToken, "wrong-header").
					Expect().Status(http.StatusUnauthorized).Body().Contains("Unauthorised access to API\n")
			})
		})

		Convey("When an authorised request to get a list of resources is made with an invalid filter value for 'state'", func() {
			Convey("Then return a status of bad request (400) with message `Bad request - invalid filter state values`", func() {
				datasetAPI.GET("/instances").WithQuery("state", "foo").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - invalid filter state values: [foo]\n")
			})
		})

		if err := mongo.Teardown(database, "instances", "_id", instanceID); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}
