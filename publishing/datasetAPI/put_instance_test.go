package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gedge/mgo"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyPutInstance(t *testing.T) {

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	publishedInstanceID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	neo4JStore, err := neo4j.NewDatastore(cfg.Neo4jAddr, "", "")
	if err != nil {
		log.ErrorC("unable to connect to neo4j", err, nil)
		t.FailNow()
	}

	Convey("Given an instance has been created by an import job", t, func() {
		publishedInstance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      publishedInstanceID,
			Update:     validPublishedInstanceData(datasetID, edition, publishedInstanceID),
		}

		completedInstance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      instanceID,
			Update:     validCompletedInstanceData(datasetID, edition, instanceID),
		}

		editionDoc := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "editions",
			Key:        "_id",
			Value:      editionID,
			Update:     ValidPublishedEditionData(datasetID, editionID, edition),
		}

		if err := mongo.Setup(publishedInstance, completedInstance, editionDoc); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When a PUT request is made to update instance meta data", func() {
			Convey("Then the instance is updated and return a status ok (200)", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTFullInstanceJSON)).
					Expect().Status(http.StatusOK)

				instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", instanceID)
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
				WithHeader(florenceTokenName, florenceToken).
				WithBytes([]byte(validPUTCompletedInstanceJSON)).
				Expect().Status(http.StatusOK)

			instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", instanceID)
			if err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("Was unable to remove test data", err, nil)
					os.Exit(1)
				}
			}

			So(instance.InstanceID, ShouldEqual, instanceID)
			So(instance.State, ShouldEqual, "completed")

			Convey("When a PUT request is made to update instance meta data and set state to `edition-confirmed`", func() {

				count, err := neo4JStore.CreateInstanceNode(instanceID)
				if err != nil {
					t.Errorf("failed to create neo4j instance node: [%v]\n error: [%v]\n", instanceID, err)
					t.FailNow()
				}
				So(count, ShouldEqual, 1)

				Convey("Then the instance is updated and return a status ok (200)", func() {

					datasetAPI.PUT("/instances/{instance_id}", instanceID).
						WithHeader(florenceTokenName, florenceToken).
						WithBytes([]byte(validPUTEditionConfirmedInstanceJSON)).
						Expect().Status(http.StatusOK)

					instance, err := mongo.GetInstance(cfg.MongoDB, "instances", "_id", instanceID)
					if err != nil {
						if err != mgo.ErrNotFound {
							log.ErrorC("Was unable to remove test data", err, nil)
							os.Exit(1)
						}
					}

					So(instance.InstanceID, ShouldEqual, instanceID)
					checkInstanceDoc(datasetID, instanceID, "edition-confirmed", instance)
					So(instance.Version, ShouldEqual, 2)

					// Check edition document has been created
					edition, err := mongo.GetEdition(cfg.MongoDB, "editions", "next.links.self.href", instance.Links.Edition.HRef)
					if err != nil {
						if err != mgo.ErrNotFound {
							log.ErrorC("Was unable to remove test data", err, nil)
							os.Exit(1)
						}
					}

					checkEditionDoc(datasetID, instanceID, edition.Next)

					Convey("and the dataset_id, edition and version values are set a properties on the neo4j instance node", func() {

						instanceProps, err := neo4JStore.GetInstanceProperties(instanceID)
						if err != nil {
							t.Errorf("failed to get properties from neo4j instance node: [%v]\n error: [%v]\n", instanceID, err)
							t.FailNow()
						}

						So(instanceProps["dataset_id"], ShouldEqual, datasetID)
						So(instanceProps["edition"], ShouldEqual, instance.Edition)
						So(instanceProps["version"], ShouldEqual, instance.Version)
					})

					if instance.Links.Edition != nil {
						e := &mongo.Doc{
							Database:   cfg.MongoDB,
							Collection: "editions",
							Key:        "links.self.href",
							Value:      instance.Links.Edition.HRef,
						}

						if err := mongo.Teardown(e); err != nil {
							if err != mgo.ErrNotFound {
								os.Exit(1)
							}
						}

						if err := neo4JStore.CleanUpInstance(instanceID); err != nil {
							t.Errorf("failed to cleanup neo4j instances: [%v]\n error: [%v]\n", instanceID, err)
							t.FailNow()
						}
					}
				})
			})
		})

		if err := mongo.Teardown(publishedInstance, completedInstance, editionDoc); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

// TODO test to be able to update version after being published with an alert?
func TestFailureToPutInstance(t *testing.T) {

	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	// Convey("Given an instance does not exist", t, func() {
	//
	// 	Convey("When an authorised PUT request is made to update instance with meta data", func() {
	// 		Convey("Then the response return a status not found (404) with message `instance not found`", func() {
	//
	// 			datasetAPI.PUT("/instances/{instance_id}", instanceID).
	// 				WithHeader(florenceTokenName, florenceToken).
	// 				WithBytes([]byte(validPUTFullInstanceJSON)).
	// 				Expect().Status(http.StatusNotFound).
	// 				Body().Contains("instance not found")
	//
	// 		})
	// 	})
	// })

	Convey("Given a created instance exists", t, func() {
		docs, err := setupInstance(datasetID, edition, instanceID)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		// Convey("When an authorised PUT request is made to update instance with invalid json", func() {
		// 	Convey("Then the response return a status not found (400) with message `failed to parse json body: unexpected end of JSON input`", func() {
		//
		// 		datasetAPI.PUT("/instances/{instance_id}", instanceID).
		// 			WithHeader(florenceTokenName, florenceToken).
		// 			WithBytes([]byte("{")).
		// 			Expect().Status(http.StatusBadRequest).
		// 			Body().Contains("failed to parse json body")
		//
		// 	})
		// })
		//
		// Convey("When an unauthorised PUT request is made to update an instance resource with an invalid authentication header", func() {
		// 	Convey("Then fail to update resource and return a status unauthorized (401)", func() {
		//
		// 		datasetAPI.PUT("/instances/{instance_id}", instanceID).
		// 			WithBytes([]byte(validPUTFullInstanceJSON)).
		// 			WithHeader(florenceTokenName, unauthorisedAuthToken).
		// 			Expect().Status(http.StatusUnauthorized)
		//
		// 	})
		// })
		//
		// Convey("When no authentication header is provided in PUT request to update an instance resource", func() {
		// 	Convey("Then fail to update resource and return a status unauthorized (401)", func() {
		//
		// 		datasetAPI.PUT("/instances/{instance_id}", instanceID).
		// 			WithBytes([]byte(validPUTFullInstanceJSON)).
		// 			Expect().Status(http.StatusUnauthorized)
		//
		// 	})
		// })
		//
		// Convey("When a PUT request is made to update instance state to `edition-confirmed`", func() {
		// 	Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update resource, expected resource to have a state of completed`", func() {
		//
		// 		datasetAPI.PUT("/instances/{instance_id}", instanceID).
		// 			WithHeader(florenceTokenName, florenceToken).
		// 			WithBytes([]byte(`{"state": "edition-confirmed"}`)).
		// 			Expect().Status(http.StatusForbidden).Body().
		// 			Contains("unable to update resource, expected resource to have a state of completed")
		//
		// 	})
		// })
		//
		// Convey("When a PUT request is made to update instance state to `associated`", func() {
		// 	Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update resource, expected resource to have a state of edition-confirmed`", func() {
		//
		// 		datasetAPI.PUT("/instances/{instance_id}", instanceID).
		// 			WithHeader(florenceTokenName, florenceToken).
		// 			WithBytes([]byte(`{"state": "associated"}`)).
		// 			Expect().Status(http.StatusForbidden).
		// 			Body().Contains("unable to update resource, expected resource to have a state of edition-confirmed")
		//
		// 	})
		// })
		//
		// Convey("When a PUT request is made to update instance state to `published`", func() {
		// 	Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update resource, expected resource to have a state of edition-confirmed`", func() {
		//
		// 		datasetAPI.PUT("/instances/{instance_id}", instanceID).
		// 			WithHeader(florenceTokenName, florenceToken).
		// 			WithBytes([]byte(`{"state": "published"}`)).
		// 			Expect().Status(http.StatusForbidden).
		// 			Body().Contains("unable to update resource, expected resource to have a state of associated")
		//
		// 	})
		// })

		Convey("When a PUT request is made to update instance state to `fake-state`", func() {
			Convey("Then fail to update resource and return a status of bad request (400) with a message `Bad request - invalid filter state values: [fake-state]`", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"state": "fake-state"}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("bad request - invalid filter state values: [fake-state]")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestUpdatingStateOnPublishedDataset(t *testing.T) {
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an dataset has been published", t, func() {
		Convey("When a valid authorised PUT request is made to update the state to `completed`", func() {
			instance := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validPublishedInstanceData(datasetID, edition, instanceID),
			}
			if err := mongo.Setup(instance); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("Then a forbidden http status is returned", func() {

				datasetAPI.PUT("/instances/{instance_id}", instanceID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"state": "completed"}`)).
					Expect().Status(http.StatusForbidden)

			})

			if err := mongo.Teardown(instance); err != nil {
				log.ErrorC("Was unable to teardown test", err, nil)
				os.Exit(1)
			}
		})
	})
}

func setupInstance(datasetID, edition, instanceID string) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validCreatedInstanceData(datasetID, edition, instanceID, "created"),
	}

	docs = append(docs, datasetDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}

func checkInstanceDoc(datasetID, instanceID, state string, instance mongo.Instance) {
	alert := mongo.Alert{
		Date:        "2017-04-05",
		Description: "All data entries (observations) for Plymouth have been updated",
		Type:        "Correction",
	}

	dimension := mongo.CodeList{
		Description: "The age ranging from 16 to 75+",
		HRef:        "http://localhost:22400//code-lists/43513D18-B4D8-4227-9820-492B2971E7T5",
		ID:          "43513D18-B4D8-4227-9820-492B2971E7T5",
		Name:        "age",
	}

	latestChange := mongo.LatestChange{
		Description: "change to the period frequency from quarterly to monthly",
		Name:        "Changes to the period frequency",
		Type:        "Summary of Changes",
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
		Dimensions: &mongo.IDLink{
			ID:   "",
			HRef: "http://localhost:22000/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
		},
		Edition: &mongo.IDLink{
			ID:   "2017",
			HRef: "http://localhost:22000/datasets/" + datasetID + "/editions/2017",
		},
	}

	if state == "edition-confirmed" {
		links.Dimensions = &mongo.IDLink{
			HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
		}
		links.Edition = &mongo.IDLink{
			ID:   "2017",
			HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017",
		}
		links.Version = &mongo.IDLink{
			ID:   "2",
			HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/2",
		}
	}

	observations := 1000

	temporal := mongo.TemporalFrequency{
		EndDate:   "2016-10-10",
		Frequency: "monthly",
		StartDate: "2014-10-10",
	}

	So(instance.Alerts, ShouldResemble, &[]mongo.Alert{alert})
	So(instance.Dimensions, ShouldResemble, []mongo.CodeList{dimension})
	So(instance.Edition, ShouldEqual, "2017")
	So(instance.Headers, ShouldResemble, &[]string{"time", "geography"})
	So(instance.ImportTasks, ShouldNotBeNil)
	So(instance.ImportTasks.ImportObservations, ShouldNotBeNil)
	So(instance.ImportTasks.ImportObservations.InsertedObservations, ShouldEqual, observations)
	So(instance.LastUpdated, ShouldNotBeNil)
	So(instance.LatestChanges, ShouldResemble, &[]mongo.LatestChange{latestChange})
	So(instance.Links, ShouldResemble, links)
	So(instance.ReleaseDate, ShouldEqual, "2017-11-11")
	So(instance.State, ShouldEqual, state)
	So(instance.Temporal, ShouldResemble, &[]mongo.TemporalFrequency{temporal})
	So(instance.TotalObservations, ShouldEqual, observations)

	return
}

func checkEditionDoc(datasetID, instanceID string, editionDoc *mongo.Edition) {
	log.Info("edition", log.Data{"edition": editionDoc})
	So(editionDoc.Edition, ShouldEqual, "2017")
	So(editionDoc.Links.Dataset.ID, ShouldEqual, datasetID)
	So(editionDoc.Links.Dataset.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID)
	So(editionDoc.Links.LatestVersion.ID, ShouldEqual, "2")
	So(editionDoc.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID+"/editions/2017/versions/2")
	So(editionDoc.Links.Self.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID+"/editions/2017")
	So(editionDoc.Links.Versions.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetID+"/editions/2017/versions")
	So(editionDoc.State, ShouldEqual, "edition-confirmed")

	return
}
