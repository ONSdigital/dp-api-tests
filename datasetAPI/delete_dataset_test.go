package datasetAPI

import (
	"testing"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
)

func TestDeleteDataset(t *testing.T) {

	datasetID := uuid.NewV4().String()
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a dataset with the an id of ["+datasetID+"] exists", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).
			WithHeader(internalToken, internalTokenID).
			WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID).
				WithHeader(internalToken, internalTokenID)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusNoContent)
			})
		})
	})
}

func TestDeleteDataset_Idempotent(t *testing.T) {

	datasetID := uuid.NewV4().String()
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a dataset with the an id of ["+datasetID+"] does not already exist", t, func() {

		Convey("When an authorised DELETE request is made to delete a dataset resource", func() {

			request := datasetAPI.DELETE("/datasets/{id}", datasetID).
				WithHeader(internalToken, internalTokenID)

			Convey("Then the expected response is returned", func() {
				request.Expect().Status(http.StatusNoContent)
			})
		})
	})
}
