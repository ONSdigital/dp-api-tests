package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyDeleteDataset(t *testing.T) {

	datasetID := uuid.NewV4().String()
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a dataset with the an id of ["+datasetID+"] exists", t, func() {

		associatedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      datasetID,
			Update:     validAssociatedDatasetData(datasetID),
		}

		if err := mongo.Setup(associatedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID).
				WithHeader(florenceTokenName, florenceToken)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusNoContent)
			})
		})

		if err := mongo.Teardown(associatedDataset); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})

	// Check idempotent request, if resource is already deleted it should respond with 204
	Convey("Given a dataset with the an id of ["+datasetID+"] does not already exist", t, func() {

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID).
				WithHeader(florenceTokenName, florenceToken)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusNoContent)
			})
		})
	})
}

func TestFailureToDeleteDataset(t *testing.T) {

	datasetID := uuid.NewV4().String()
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
	secondDatasetID := uuid.NewV4().String()

	Convey("Given a published dataset with the an id of ["+datasetID+"] exists", t, func() {
		publishedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      datasetID,
			Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
		}

		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID).
				WithHeader(florenceTokenName, florenceToken)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusForbidden).Body().Contains("a published dataset cannot be deleted")
			})
		})

		if err := mongo.Teardown(publishedDataset); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})

	Convey("Given an associated dataset with the an id of ["+secondDatasetID+"] exists", t, func() {

		associatedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      secondDatasetID,
			Update:     validAssociatedDatasetData(secondDatasetID),
		}

		if err := mongo.Setup(associatedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID).
				WithHeader(florenceTokenName, unauthorisedAuthToken)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a DELETE request is made to delete a dataset resource without authentication", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(associatedDataset); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}
