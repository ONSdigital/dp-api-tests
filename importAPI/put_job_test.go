package importAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestUpdateImportJobState(t *testing.T) {

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "id",
		Value:      instanceID,
		Update:     validCreatedInstanceData,
	}

	if err := mongo.Setup(importJob, instance); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a valid authenticated request", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status OK (200)", func() {

				importAPI.PUT("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusOK)
			})
		})
	})

	if err := mongo.Teardown(importJob, instance); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestUpdateImportJobStateUnauthorised(t *testing.T) {

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "id",
		Value:      instanceID,
		Update:     validCreatedInstanceData,
	}

	if err := mongo.Setup(importJob, instance); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request with no Authorization header", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status 404 not found", func() {
				importAPI.PUT("/jobs/{id}", jobID).
					WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	Convey("Given a request with an unauthorised Authorization header", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status 401 unauthorised", func() {

				importAPI.PUT("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, unauthorisedServiceAuthToken).
					WithBytes([]byte(validPUTJobJSON)).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	if err := mongo.Teardown(importJob, instance); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func TestFailureToUpdateAnImportJob(t *testing.T) {

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: collection,
		Key:        "id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	if err := mongo.Setup(importJob); err != nil {
		log.ErrorC("Failed to set up test data", err, nil)
		os.Exit(1)
	}

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given a request for a job that does not exist", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status not found (404)", func() {
				importAPI.PUT("/jobs/{id}", invalidJobID).
					WithBytes([]byte(validPUTJobJSON)).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	Convey("Given a request with an invalid job JSON body", t, func() {
		Convey("When update job is called", func() {
			Convey("Then the response returns status bad request (400)", func() {
				importAPI.PUT("/jobs/{id}", jobID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest)
			})
		})
	})

	if err := mongo.Teardown(importJob); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}
