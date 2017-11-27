package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyPutInstance(t *testing.T) {

	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance has been created by an import job", t, func() {
		d, err := setupInstance(datasetID, edition, instanceID)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When a PUT request is made to update instance meta data", func() {
			Convey("Then the instance is updated and return a status ok (200)", func() {
				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTFullInstanceJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(database, "instances", "_id", instanceID)
				if err != nil {
					if err != mgo.ErrNotFound {
						log.ErrorC("Was unable to remove test data", err, nil)
						os.Exit(1)
					}
				}

				log.Debug("next instance", log.Data{"instance": instance})

				So(instance.InstanceID, ShouldEqual, instanceID)
				checkInstanceDoc(datasetID, instanceID, "completed", instance)
			})
		})

		Convey("and is updated to a state of `completed`", func() {

			datasetAPI.PUT("/instances/{instance_id}", instanceID).
				WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTCompletedInstanceJSON)).
				Expect().Status(http.StatusOK)

			instance, err := mongo.GetInstance(database, "instances", "_id", instanceID)
			if err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("Was unable to remove test data", err, nil)
					os.Exit(1)
				}
			}

			So(instance.InstanceID, ShouldEqual, instanceID)
			So(instance.State, ShouldEqual, "completed")

			Convey("When a PUT request is made to update instance meta data and set state to `edition-confirmed`", func() {
				Convey("Then the instance is updated and return a status ok (200)", func() {
					datasetAPI.PUT("/instances/{instance_id}", instanceID).
						WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTEditionConfirmedInstanceJSON)).
						Expect().Status(http.StatusOK)

					instance, err := mongo.GetInstance(database, "instances", "_id", instanceID)
					if err != nil {
						if err != mgo.ErrNotFound {
							log.ErrorC("Was unable to remove test data", err, nil)
							os.Exit(1)
						}
					}

					So(instance.InstanceID, ShouldEqual, instanceID)
					checkInstanceDoc(datasetID, instanceID, "edition-confirmed", instance)
					So(instance.Version, ShouldEqual, 1)

					// Check edition document has been created
					edition, err := mongo.GetEdition(database, "editions", "links.self.href", instance.Links.Edition.HRef)
					if err != nil {
						if err != mgo.ErrNotFound {
							log.ErrorC("Was unable to remove test data", err, nil)
							os.Exit(1)
						}
					}

					checkEditionDoc(datasetID, instanceID, edition)

					if err := mongo.Teardown(database, "editions", "links.self.href", instance.Links.Edition.HRef); err != nil {
						if err != mgo.ErrNotFound {
							log.ErrorC("Was unable to remove test data", err, nil)
							os.Exit(1)
						}
					}
				})
			})
		})

		if err := mongo.TeardownMany(d); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToPutInstance(t *testing.T) {

	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance does not exist", t, func() {

		Convey("When an authorised PUT request is made to update instance with meta data", func() {
			Convey("Then the response return a status not found (404) with message `Instance not found`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTFullInstanceJSON)).
					Expect().Status(http.StatusNotFound).Body().Contains("Instance not found\n")
			})
		})
	})

	Convey("Given a created instance exists", t, func() {
		d, err := setupInstance(datasetID, edition, instanceID)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised PUT request is made to update instance with invalid json", func() {
			Convey("Then the response return a status not found (400) with message `Failed to parse json body: unexpected end of JSON input`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).Body().Contains("Failed to parse json body: unexpected end of JSON input\n")
			})
		})

		Convey("When an unauthorised PUT request is made to update an instance resource with an invalid authentication header", func() {
			Convey("Then fail to update resource and return a status unauthorized (401) with a message `Unauthorised access to API`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).WithBytes([]byte(validPUTFullInstanceJSON)).
					WithHeader(internalToken, invalidInternalTokenID).Expect().Status(http.StatusUnauthorized).
					Body().Contains("Unauthorised access to API\n")
			})
		})

		Convey("When no authentication header is provided in PUT request to update an instance resource", func() {
			Convey("Then fail to update resource and return a status of unauthorized (401) with a message `No authentication header provided`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).WithBytes([]byte(validPUTFullInstanceJSON)).
					Expect().Status(http.StatusUnauthorized).
					Body().Contains("No authentication header provided\n")
			})
		})

		Convey("When a PUT request is made to update instance state to `edition-confirmed`", func() {
			Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update resource, expected resource to have a state of completed`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte(`{"state": "edition-confirmed"}`)).
					Expect().Status(http.StatusForbidden).Body().Contains("Unable to update resource, expected resource to have a state of completed\n")
			})
		})

		Convey("When a PUT request is made to update instance state to `associated`", func() {
			Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update resource, expected resource to have a state of edition-confirmed`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte(`{"state": "associated"}`)).
					Expect().Status(http.StatusForbidden).Body().Contains("Unable to update resource, expected resource to have a state of edition-confirmed\n")
			})
		})

		Convey("When a PUT request is made to update instance state to `published`", func() {
			Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update resource, expected resource to have a state of edition-confirmed`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte(`{"state": "published"}`)).
					Expect().Status(http.StatusForbidden).Body().Contains("Unable to update resource, expected resource to have a state of associated\n")
			})
		})

		Convey("When a PUT request is made to update instance state to `fake-state`", func() {
			Convey("Then fail to update resource and return a status of bad request (400) with a message `Bad request - invalid filter state values: [fake-state]`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(internalToken, internalTokenID).WithBytes([]byte(`{"state": "fake-state"}`)).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - invalid filter state values: [fake-state]\n")
			})
		})

		if err := mongo.TeardownMany(d); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func setupInstance(datasetID, edition, instanceID string) (*mongo.ManyDocs, error) {
	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validCreatedInstanceData(datasetID, edition, instanceID),
	}

	docs = append(docs, datasetDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.SetupMany(d); err != nil {
		return nil, err
	}

	return d, nil
}

func checkInstanceDoc(datasetID, instanceID, state string, instance mongo.Instance) {
	dimension := mongo.CodeList{
		Description: "The age ranging from 16 to 75+",
		HRef:        "http://localhost:22400//code-lists/43513D18-B4D8-4227-9820-492B2971E7T5",
		ID:          "43513D18-B4D8-4227-9820-492B2971E7T5",
		Name:        "age",
	}

	links := mongo.InstanceLinks{
		Job: &mongo.IDLink{
			ID:   "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			HRef: "http://localhost:22000/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		},
		Dataset: &mongo.IDLink{
			ID:   datasetID,
			HRef: "http://localhost:22000/datasets/" + datasetID,
		},
		Self: &mongo.IDLink{
			HRef: cfg.DatasetAPIURL + "/instances/" + instanceID,
		},
		Spatial: &mongo.IDLink{
			HRef: "http://ons.gov.uk/geography-list",
		},
	}

	if state == "edition-confirmed" {
		links.Dimensions = &mongo.IDLink{
			HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/1/dimensions",
		}
		links.Edition = &mongo.IDLink{
			ID:   "2017",
			HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017",
		}
		links.Version = &mongo.IDLink{
			ID:   "1",
			HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/1",
		}
	}

	observations := 1000

	temporal := mongo.TemporalFrequency{
		EndDate:   "2016-10-10",
		Frequency: "monthly",
		StartDate: "2014-10-10",
	}

	So(instance.Dimensions, ShouldResemble, []mongo.CodeList{dimension})
	So(instance.Edition, ShouldEqual, "2017")
	So(instance.Headers, ShouldResemble, &[]string{"time", "geography"})
	So(instance.LastUpdated, ShouldNotBeNil)
	So(instance.Links, ShouldResemble, links)
	So(instance.ReleaseDate, ShouldEqual, "2017-11-11")
	So(instance.State, ShouldEqual, state)
	So(instance.Temporal, ShouldResemble, &[]mongo.TemporalFrequency{temporal})
	So(instance.TotalObservations, ShouldResemble, &observations)
	So(instance.InsertedObservations, ShouldResemble, &observations)

	return
}

func checkEditionDoc(datasetID, instanceID string, editionDoc mongo.Edition) {
	So(editionDoc.Edition, ShouldEqual, "2017")
	So(editionDoc.Links.Dataset.ID, ShouldEqual, datasetID)
	So(editionDoc.Links.Dataset.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID)
	So(editionDoc.Links.LatestVersion.ID, ShouldEqual, "1")
	So(editionDoc.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID+"/editions/2017/versions/1")
	So(editionDoc.Links.Self.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID+"/editions/2017")
	So(editionDoc.Links.Versions.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID+"/editions/2017/versions")
	So(editionDoc.State, ShouldEqual, "created")

	return
}