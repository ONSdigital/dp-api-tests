package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

// This updates the instance resource with a dimension object to
// it's list of dimensions in dimension array
func TestSuccessfullyPutImportTasks(t *testing.T) {
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
			Value:      ids.InstanceSubmitted,
			Update:     validSubmittedInstanceData(ids.DatasetPublished, edition, ids.InstanceSubmitted, submitted, ids.UniqueTimestamp),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		docs = append(docs, instance)

		Convey(`When a PUT request to update an observations import task on an instance resource`, func() {

			Convey("Then the instance resource is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", ids.InstanceSubmitted).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validObservationImportTaskJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", ids.InstanceSubmitted)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": ids.InstanceSubmitted})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, ids.InstanceSubmitted)
				So(instance.ImportTasks.ImportObservations.State, ShouldEqual, "completed")
			})
		})

		Convey(`When a PUT request to update a hierarchy import task on an instance resource`, func() {

			Convey("Then the instance resource is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", ids.InstanceSubmitted).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validHierarchyImportTaskJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", ids.InstanceSubmitted)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": ids.InstanceSubmitted})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, ids.InstanceSubmitted)
				So(instance.ImportTasks.BuildHierarchyTasks[0].DimensionName, ShouldEqual, "geography")
				So(instance.ImportTasks.BuildHierarchyTasks[0].State, ShouldEqual, "completed")
				So(instance.ImportTasks.BuildHierarchyTasks[0].CodeListID, ShouldEqual, "K02000001")
			})
		})

		Convey(`When a PUT request to update a search index import task on an instance resource`, func() {

			Convey("Then the instance resource is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", ids.InstanceSubmitted).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validSearchIndexImportTaskJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", ids.InstanceSubmitted)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": ids.InstanceSubmitted})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, ids.InstanceSubmitted)
				So(instance.ImportTasks.SearchTasks[0].DimensionName, ShouldEqual, "geography")
				So(instance.ImportTasks.SearchTasks[0].State, ShouldEqual, "completed")
			})
		})

		Convey(`When a PUT request to update multiple import tasks on an instance resource`, func() {

			Convey("Then the instance resource is updated and response returns status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", ids.InstanceSubmitted).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validMultipleImportTaskJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", ids.InstanceSubmitted)
				if err != nil {
					log.ErrorC("Was unable to retrieve instance test data", err, log.Data{"instance_id": ids.InstanceSubmitted})
					os.Exit(1)
				}

				So(instance.InstanceID, ShouldEqual, ids.InstanceSubmitted)
				So(instance.ImportTasks.ImportObservations.State, ShouldEqual, "completed")
				So(instance.ImportTasks.ImportObservations.InsertedObservations, ShouldEqual, 1000)
				So(instance.ImportTasks.BuildHierarchyTasks[0].DimensionName, ShouldEqual, "geography")
				So(instance.ImportTasks.BuildHierarchyTasks[0].State, ShouldEqual, "completed")
				So(instance.ImportTasks.BuildHierarchyTasks[0].CodeListID, ShouldEqual, "K02000001")
				So(instance.ImportTasks.SearchTasks[0].DimensionName, ShouldEqual, "geography")
				So(instance.ImportTasks.SearchTasks[0].State, ShouldEqual, "completed")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToPutImportTasks(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	instances := make(map[string]string)
	instances[submitted] = ids.InstanceSubmitted
	instances[invalid] = ids.InstanceInvalid

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {
		Convey(`When an authorised PUT request to update import task against
			an instance resource`, func() {

			Convey("Then the response return a status not found (404) with message `instance not found`", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validObservationImportTaskJSON)).
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

		Convey(`When an unauthorised PUT request to update import task against
			an instance resource with an invalid authentication header`, func() {

			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					WithBytes([]byte(validObservationImportTaskJSON)).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When no authentication header is provided in PUT request
					to update import task against an instance resource`, func() {

			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithBytes([]byte(validObservationImportTaskJSON)).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json is invalid`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'unexpected end of JSON input'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("unexpected end of JSON input")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains no import tasks`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'request body does not contain any import tasks'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte("{}")).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - request body does not contain any import tasks")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains empty observation import task`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'invalid import observation task, must include state'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"import_observations": {}}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - invalid import observation task, must include state")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains empty array of hierarchies import task`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'missing hierarchy task'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_hierarchies": []}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - missing hierarchy task")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains a hierarchy import task with
			empty 'dimension_name' and 'state'`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'missing mandatory fields: [dimension_name state]'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_hierarchies": [{}]}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - missing mandatory fields: [dimension_name state]")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains a hierarchy import task with
			a state of 'gobbledygook'`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'invalid task state value: gobbledygook'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_hierarchies": [{"state": "gobbledygook", "dimension_name": "geography"}]}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - invalid task state value: gobbledygook")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains a hierarchy import task with
			a dimension name that does not exist`, func() {

			Convey(`Then the response return a status not found (404)
						with message 'age hierarchy import task does not exist'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_hierarchies": [{"state": "completed", "dimension_name": "age"}]}`)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("age hierarchy import task does not exist")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains empty array of search indexes import task`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'missing search index task'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_search_indexes": []}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - missing search index task")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains a search index import task with
			empty 'dimension_name' and 'state'`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'missing mandatory fields: [dimension_name state]'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_search_indexes": [{}]}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - missing mandatory fields: [dimension_name state]")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains a search index import task with
			a state of 'gobbledygook'`, func() {

			Convey(`Then the response return a status bad request (400)
						with message 'invalid task state value: gobbledygook'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_search_indexes": [{"state": "gobbledygook", "dimension_name": "geography"}]}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - invalid task state value: gobbledygook")

			})
		})

		Convey(`When an authorised PUT request to update import task against an
			instance resource but the json contains a search index import task with
			a dimension name that does not exist`, func() {

			Convey(`Then the response return a status not found (404)
						with message 'age hierarchy import task does not exist'`, func() {

				datasetAPI.PUT("/instances/{instance_id}/import_tasks", instances[submitted]).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"build_search_indexes": [{"state": "completed", "dimension_name": "age"}]}`)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("age search index import task does not exist")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}
