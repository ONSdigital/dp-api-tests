package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gedge/mgo"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	datasetAPIModel "github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyPutInstanceDimensionOptionNodeID(t *testing.T) {

	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"
	dimensionOptionID := uuid.NewV4().String()
	nodeID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance has been created by an import job", t, func() {
		var docs []*mongo.Doc

		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      instanceID,
			Update:     validCreatedInstanceData(datasetID, edition, instanceID, "created"),
		}

		dimensionOption := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "dimension.options",
			Key:        "_id",
			Value:      dimensionOptionID,
			Update:     validTimeDimensionsData(dimensionOptionID, instanceID),
		}

		if err := mongo.Setup(instance, dimensionOption); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		docs = append(docs, instance, dimensionOption)

		Convey("When a PUT request is made to add node ID to dimension option for instance", func() {
			Convey("Then the dimension option is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.45/node_id/{node_id}", instanceID, nodeID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					Expect().Status(http.StatusOK)

				dimensionOption, err := mongo.GetDimensionOption(cfg.MongoDB, "dimension.options", "_id", dimensionOptionID)
				if err != nil {
					log.ErrorC("Was unable to retrieve dimension option test data", err, log.Data{"_id": dimensionOptionID, "instance_id": instanceID, "node_id": nodeID})
					os.Exit(1)
				}

				checkDimensionOptionDocWithNodeID(instanceID, nodeID, &dimensionOption)

				dimensionOptionDoc := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: "dimension.options",
					Key:        "instance_id",
					Value:      instanceID,
				}

				docs = append(docs, dimensionOptionDoc)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToPutDimensionOptionNodeID(t *testing.T) {
	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	invalidInstanceID := uuid.NewV4().String()
	edition := "2017"
	nodeID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {
		Convey(`When an authorised PUT request is made to update dimension option
			with a 'node_id' for an instance`, func() {

			Convey(`Then the response return a status not found (404)
				with message 'Instance not found'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instanceID, nodeID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Instance not found")

			})
		})
	})

	Convey("Given a created instance exists", t, func() {
		// use setupInstance in file post_instance_dimension_test.go
		docs, err := setupInstances(datasetID, edition, instanceID, invalidInstanceID)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey(`When an unauthorised PUT request is made to update dimension option
				 for an instance with a node id`, func() {
			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instanceID, nodeID).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in PUT request to update
				dimension option for an instance with a node id`, func() {
			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instanceID, nodeID).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised PUT request is made to update dimension option
					for an instance but instance has an invalid state`, func() {

			Convey(`Then the response return a status internal server error (500)
						with message 'internal error'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", invalidInstanceID, nodeID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusInternalServerError).
					Body().Contains("internal error")

			})
		})

		Convey(`When an authorised PUT request is made to update dimension option
					for an instance but the dimension option does not exist`, func() {

			Convey(`Then the response return a status not found (404)
						with message 'failed to parse json body'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instanceID, nodeID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Dimension option not found")

			})
		})

		Convey(`When an authorised POST request is made to create dimension option
					for an instance but the json is missing mandatory field 'dimension'`, func() {

			Convey(`Then the response return a status not found (400)
						with message 'missing properties in JSON'`, func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(invalidPOSTDimensionJSONMissingDimension)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing properties in JSON")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func checkDimensionOptionDocWithNodeID(instanceID, nodeID string, dimensionOptionDoc *datasetAPIModel.DimensionOption) {

	So(dimensionOptionDoc.Name, ShouldEqual, "time")
	So(dimensionOptionDoc.Label, ShouldEqual, "")
	So(dimensionOptionDoc.LastUpdated, ShouldNotBeEmpty)
	So(dimensionOptionDoc.Links.Code.ID, ShouldEqual, "202.45")
	So(dimensionOptionDoc.Links.Code.HRef, ShouldEqual, cfg.CodeListAPIURL+"/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45")
	So(dimensionOptionDoc.Links.CodeList.ID, ShouldEqual, "64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	So(dimensionOptionDoc.Links.CodeList.HRef, ShouldEqual, cfg.CodeListAPIURL+"/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	So(dimensionOptionDoc.NodeID, ShouldEqual, nodeID)
	So(dimensionOptionDoc.Option, ShouldEqual, "202.45")

	return
}
