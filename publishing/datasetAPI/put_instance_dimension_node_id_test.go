package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/globalsign/mgo"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	datasetAPIModel "github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyPutInstanceDimensionOptionNodeID(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance has been created by an import job", t, func() {
		var docs []*mongo.Doc

		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstanceCreated,
			Update:     validCreatedInstanceData(ids.DatasetPublished, edition, ids.InstanceCreated, created, ids.UniqueTimestamp),
		}

		dimensionOption := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "dimension.options",
			Key:        "_id",
			Value:      ids.Dimension,
			Update:     validTimeDimensionsData(ids.Dimension, ids.InstanceCreated),
		}

		if err := mongo.Setup(instance, dimensionOption); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		docs = append(docs, instance, dimensionOption)

		Convey("When a PUT request is made to add node ID to dimension option for instance", func() {
			Convey("Then the dimension option is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.45/node_id/{node_id}", ids.InstanceCreated, ids.Node).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					Expect().Status(http.StatusOK)

				dimensionOption, err := mongo.GetDimensionOption(cfg.MongoDB, "dimension.options", "_id", ids.Dimension)
				if err != nil {
					log.ErrorC("Was unable to retrieve dimension option test data", err, log.Data{"_id": ids.Dimension, "instance_id": ids.InstanceCreated, "node_id": ids.Node})
					os.Exit(1)
				}

				checkDimensionOptionDocWithNodeID(ids.InstanceCreated, ids.Node, &dimensionOption)

				dimensionOptionDoc := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: "dimension.options",
					Key:        "instance_id",
					Value:      ids.InstanceCreated,
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
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	instances := make(map[string]string)
	instances[created] = ids.InstanceCreated
	instances[invalid] = ids.InstanceInvalid

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {
		Convey(`When an authorised PUT request is made to update dimension option
			with a 'node_id' for an instance`, func() {

			Convey(`Then the response return a status not found (404)
				with message 'instance not found'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instances[created], ids.Node).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("instance not found")

			})
		})
	})

	Convey("Given a created instance exists", t, func() {
		docs, err := setupInstances(ids.DatasetPublished, edition, ids.UniqueTimestamp, instances)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey(`When an unauthorised PUT request is made to update dimension option
				 for an instance with a node id`, func() {
			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instances[created], ids.Node).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in PUT request to update
				dimension option for an instance with a node id`, func() {
			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instances[created], ids.Node).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised PUT request is made to update dimension option
					for an instance but instance has an invalid state`, func() {

			Convey(`Then the response return a status internal server error (500)
						with message 'internal error'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/time/options/202.5/node_id/{node_id}", instances[invalid], ids.Node).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusInternalServerError).
					Body().Contains("internal error")

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
