package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gedge/mgo"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

// This updates the instance resource with a dimension object to
// it's list of dimensions in dimension array
func TestSuccessfullyPutInstanceDimension(t *testing.T) {

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
			Update:     validSubmittedInstanceData(datasetID, edition, instanceID),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		docs = append(docs, instance)

		Convey("When a PUT request is made to update dimension on an instance resource", func() {
			Convey("Then the dimension is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/geography", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTGeographyDimensionJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", instanceID)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": instanceID})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, instanceID)
				checkInstanceDimensions(&instance)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToPutInstanceDimension(t *testing.T) {
	datasetID := uuid.NewV4().String()
	edition := "2017"

	instances := make(map[string]string)
	instances[submitted] = uuid.NewV4().String()
	instances[invalid] = uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {
		Convey("When an authorised PUT request is made to update dimension on an instance resource", func() {
			Convey("Then the response return a status not found (404) with message `instance not found`", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/geography", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTGeographyDimensionJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("instance not found")

			})
		})
	})

	Convey("Given a created instance exists", t, func() {
		docs, err := setupInstances(datasetID, edition, instances)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey(`When an unauthorised PUT request is made to update dimension on an
			instance resource with an invalid authentication header`, func() {

			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/geography", instances[submitted]).
					WithBytes([]byte(validPUTGeographyDimensionJSON)).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in PUT request to update
					dimension on an instance resource`, func() {
			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/geography", instances[submitted]).
					WithBytes([]byte(validPUTGeographyDimensionJSON)).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised PUT request is made to update dimension on an
			instance resource but the json is invalid`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'failed to parse json body'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/geography", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("unexpected end of JSON input")

			})
		})

		Convey(`When an authorised PUT request is made to update dimension on an
			instance resource but instance has an invalid state`, func() {

			Convey(`Then the response return a status internal server error (500)
				with message 'internal error'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/geography", instances[invalid]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTGeographyDimensionJSON)).
					Expect().Status(http.StatusInternalServerError).
					Body().Contains("internal error")

			})
		})

		Convey(`When an authorised PUT request is made to update dimension on an
			instance resource but the 'dimension' does not exist`, func() {

			Convey(`Then the response return a status not found (404)
						with message 'dimension not found'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/dimensions/age", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTAgeDimensionJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("dimension not found")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func checkInstanceDimensions(instanceDoc *mongo.Instance) {

	So(len(instanceDoc.Dimensions), ShouldEqual, 1)
	So(instanceDoc.Dimensions[0].Description, ShouldEqual, "The sites in which this dataset spans")
	So(instanceDoc.Dimensions[0].HRef, ShouldContainSubstring, "/codelists/708064B3-A808-449B-9041-EA3A2F72CFAF")
	So(instanceDoc.Dimensions[0].ID, ShouldEqual, "708064B3-A808-449B-9041-EA3A2F72CFAF")
	So(instanceDoc.Dimensions[0].Label, ShouldEqual, "geo-sites")
	So(instanceDoc.Dimensions[0].Name, ShouldEqual, "geography")

	return
}
