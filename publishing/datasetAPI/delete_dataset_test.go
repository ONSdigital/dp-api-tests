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

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyDeleteDataset(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a dataset with the an id of ["+ids.DatasetAssociated+"] exists", t, func() {

		associatedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      ids.DatasetAssociated,
			Update:     validAssociatedDatasetData(ids.DatasetAssociated),
		}

		if err := mongo.Setup(associatedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", ids.DatasetAssociated).
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
	Convey("Given a dataset with the an id of ["+ids.DatasetAssociated+"] does not already exist", t, func() {

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", ids.DatasetAssociated).
				WithHeader(florenceTokenName, florenceToken)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusNoContent)
			})
		})
	})
}

func TestFailureToDeleteDataset(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published dataset with the an id of ["+ids.DatasetPublished+"] exists", t, func() {
		publishedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      ids.DatasetPublished,
			Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
		}

		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", ids.DatasetPublished).
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

	Convey("Given an associated dataset with the an id of ["+ids.DatasetAssociated+"] exists", t, func() {

		associatedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      ids.DatasetAssociated,
			Update:     validAssociatedDatasetData(ids.DatasetAssociated),
		}

		if err := mongo.Setup(associatedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", ids.DatasetAssociated).
				WithHeader(florenceTokenName, unauthorisedAuthToken)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a DELETE request is made to delete a dataset resource without authentication", func() {

			request := datasetAPI.DELETE("/datasets/{id}", ids.DatasetAssociated)

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
