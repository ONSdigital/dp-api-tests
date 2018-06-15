package datasetAPI

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/gedge/mgo"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyPostInstanceEvent(t *testing.T) {

	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance has been created by an import job", t, func() {

		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      instanceID,
			Update:     validSubmittedInstanceData(datasetID, edition, instanceID),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When a POST request is made to create an event on an instance resource", func() {
			Convey("Then the instance resource is updated witht he event and response returns status ok (200)", func() {

				b, err := createValidPOSTEventJSON(time.Now().UTC())
				if err != nil {
					log.ErrorC("Unable to create event test data", err, nil)
					os.Exit(1)
				}
				datasetAPI.POST("/instances/{instance_id}/events", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes(b).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", instanceID)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": instanceID})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, instanceID)

				expectedEvent := mongo.InstanceEvent{Message: "unable to add observation to neo4j", MessageOffset: "5", Type: "error"}
				checkInstanceEvent(instance.Events, expectedEvent)
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToPostInstanceEvent(t *testing.T) {
	datasetID := uuid.NewV4().String()
	edition := "2017"

	instances := make(map[string]string)
	instances[submitted] = uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
	b, err := createValidPOSTEventJSON(time.Now().UTC())
	if err != nil {
		log.ErrorC("Unable to create event test data", err, nil)
		os.Exit(1)
	}

	Convey("Given an instance does not exist", t, func() {
		Convey("When an authorised POST request to create an event against an instance resource", func() {
			Convey("Then the response return a status not found (404) with message `Instance not found`", func() {

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes(b).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Instance not found")

			})
		})
	})

	Convey("Given a created instance exists", t, func() {
		docs, err := setupInstances(datasetID, edition, instances)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey(`When an unauthorised POST request to create an event against an
			instance resource with an invalid authentication header`, func() {

			Convey("Then fail to add event to instance resource and return a status unauthorized (401)", func() {

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithBytes(b).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in POST request to create
			an event against an instance resource`, func() {

			Convey("Then fail to add event to instance resource and return a status unauthorized (401)", func() {

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithBytes(b).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised POST request to create an event against an
				instance resource but the json is invalid`, func() {

			Convey(`Then the response return a status bad request (400)
							with message 'Failed to parse json body'`, func() {

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("Failed to parse json body")

			})
		})

		Convey(`When an authorised POST request to create an event against an
				instance resource but the json is missing 'event.messsage'`, func() {

			Convey(`Then the response return a status bad request (400)
							with message 'missing properties'`, func() {

				newBytes, err := createInvalidPOSTEventJSONWithoutMessage(time.Now().UTC())
				if err != nil {
					log.ErrorC("Unable to create event test data", err, nil)
					os.Exit(1)
				}

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes(newBytes).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing properties")

			})
		})

		Convey(`When an authorised POST request to create an event against an
				instance resource but the json is missing 'event.messsage_offset'`, func() {

			Convey(`Then the response return a status bad request (400)
							with message 'missing properties'`, func() {

				newBytes, err := createInvalidPOSTEventJSONWithoutMessageOffset(time.Now().UTC())
				if err != nil {
					log.ErrorC("Unable to create event test data", err, nil)
					os.Exit(1)
				}

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes(newBytes).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing properties")

			})
		})

		Convey(`When an authorised POST request to create an event against an
				instance resource but the json is missing 'event.time'`, func() {

			Convey(`Then the response return a status bad request (400)
							with message 'missing properties'`, func() {

				newBytes, err := createInvalidPOSTEventJSONWithoutTime()
				if err != nil {
					log.ErrorC("Unable to create event test data", err, nil)
					os.Exit(1)
				}

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes(newBytes).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing properties")

			})
		})

		Convey(`When an authorised POST request to create an event against an
				instance resource but the json is missing 'event.type'`, func() {

			Convey(`Then the response return a status bad request (400)
							with message 'missing properties'`, func() {

				newBytes, err := createInvalidPOSTEventJSONWithoutType(time.Now().UTC())
				if err != nil {
					log.ErrorC("Unable to create event test data", err, nil)
					os.Exit(1)
				}

				datasetAPI.POST("/instances/{instance_id}/events", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes(newBytes).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing properties")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func checkInstanceEvent(events *[]mongo.InstanceEvent, expectedEvent mongo.InstanceEvent) {
	So(events, ShouldResemble, &[]mongo.InstanceEvent{expectedEvent})

	return
}
