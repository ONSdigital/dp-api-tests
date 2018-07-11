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
func TestSuccessfullyPutInsertedObservations(t *testing.T) {

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

		Convey(`When a PUT request to add the number of inserted
			observations against an instance resource`, func() {

			Convey("Then the instance resource is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/inserted_observations/255", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", instanceID)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": instanceID})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, instanceID)
				So(instance.ImportTasks.ImportObservations.InsertedObservations, ShouldEqual, 755)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToPutInsertedObservations(t *testing.T) {
	datasetID := uuid.NewV4().String()
	edition := "2017"

	instances := make(map[string]string)
	instances[submitted] = uuid.NewV4().String()
	instances[invalid] = uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {
		Convey(`When an authorised PUT request to add the number of inserted
			observations against an instance resource`, func() {

			Convey("Then the response return a status not found (404) with message `instance not found`", func() {

				datasetAPI.PUT("/instances/{instance_id}/inserted_observations/255", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
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

		Convey(`When an unauthorised PUT request to add the number of inserted
				observations against an instance resource with an invalid authentication header`, func() {

			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/inserted_observations/255", instances[submitted]).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in PUT request to add
					the number of inserted observations against an instance resource`, func() {

			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/inserted_observations/255", instances[submitted]).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised PUT request to add the number of inserted
					observations against an instance resource but number observations in path is not a number`, func() {

			Convey(`Then the response return a status internal server error (500)
						with message 'internal error'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/inserted_observations/twohundredandfiftyfive", instances[invalid]).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("strconv.ParseInt: parsing \"twohundredandfiftyfive\": invalid syntax")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}
