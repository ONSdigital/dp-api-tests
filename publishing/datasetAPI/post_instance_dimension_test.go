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

func TestSuccessfullyPostInstanceDimension(t *testing.T) {

	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"

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

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		docs = append(docs, instance)

		Convey("When a POST request is made to add dimension for instance", func() {
			Convey("Then the dimension option is created and response returns status ok (200)", func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					Expect().Status(http.StatusOK)

				dimensionOption, err := mongo.GetDimensionOption(cfg.MongoDB, "dimension.options", "instance_id", instanceID)
				if err != nil {
					log.ErrorC("Was unable to retrieve dimension option test data", err, log.Data{"instance_id": instanceID})
					os.Exit(1)
				}

				So(dimensionOption.InstanceID, ShouldEqual, instanceID)
				checkDimensionOptionDoc(instanceID, &dimensionOption)

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

func TestFailureToPostDimension(t *testing.T) {
	datasetID := uuid.NewV4().String()
	edition := "2017"

	instances := make(map[string]string)
	instances[created] = uuid.NewV4().String()
	instances[invalid] = uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {
		Convey("When an authorised POST request is made to add dimension option for an instance", func() {
			Convey("Then the response return a status not found (404) with message `Instance not found`", func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["created"]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
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

		Convey(`When an unauthorised POST request is made to add dimension option
			 for an instance with an invalid authentication header`, func() {
			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["created"]).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in POST request to add
			dimension option for an instance`, func() {
			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["created"]).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised POST request is made to create dimension option
			for an instance but the json is invalid`, func() {

			Convey(`Then the response return a status not found (400)
				with message 'failed to parse json body'`, func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["created"]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("failed to parse json body")

			})
		})

		Convey(`When an authorised POST request is made to create dimension option
			for an instance but instance has an invalid state`, func() {

			Convey(`Then the response return a status internal server error (500)
				with message 'internal error'`, func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["invalid"]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTAgeDimensionJSON)).
					Expect().Status(http.StatusInternalServerError).
					Body().Contains("internal error")

			})
		})

		Convey(`When an authorised POST request is made to create dimension option
			for an instance but the json is missing mandatory field 'dimension'`, func() {

			Convey(`Then the response return a status not found (400)
				with message 'missing properties in JSON'`, func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["created"]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(invalidPOSTDimensionJSONMissingDimension)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing properties in JSON")

			})
		})

		Convey(`When an authorised POST request is made to create dimension option
			but the json is missing mandatory both 'option' and 'codelist'`, func() {

			Convey(`Then the response return a status not found (400)
				with message 'missing properties in JSON'`, func() {

				datasetAPI.POST("/instances/{instance_id}/dimensions", instances["created"]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(invalidPOSTDimensionJSONMissingOptionAndCodelist)).
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

func checkDimensionOptionDoc(instanceID string, dimensionOptionDoc *datasetAPIModel.DimensionOption) {

	So(dimensionOptionDoc.Name, ShouldEqual, "age")
	So(dimensionOptionDoc.Label, ShouldEqual, "25")
	So(dimensionOptionDoc.LastUpdated, ShouldNotBeEmpty)
	So(dimensionOptionDoc.Links.Code.ID, ShouldEqual, "ABC123DEF456")
	So(dimensionOptionDoc.Links.Code.HRef, ShouldEqual, cfg.CodeListAPIURL+"/code-lists/age-list/codes/ABC123DEF456")
	So(dimensionOptionDoc.Links.CodeList.ID, ShouldEqual, "age-list")
	So(dimensionOptionDoc.Links.CodeList.HRef, ShouldEqual, cfg.CodeListAPIURL+"/code-lists/age-list")
	So(dimensionOptionDoc.NodeID, ShouldBeEmpty)
	So(dimensionOptionDoc.Option, ShouldEqual, "25")

	return
}
