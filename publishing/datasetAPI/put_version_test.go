package datasetAPI

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyUpdateVersion(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	neo4JStore, err := neo4j.NewDatastore(cfg.Neo4jAddr, "", neo4j.GenericHierarchyCPIHTestData)
	if err != nil {
		t.Errorf("unable to connect to neo4j. error: [%v]\n", err)
		log.ErrorC("unable to connect to neo4j", err, nil)
		t.FailNow()
	}

	Convey("Given an unpublished dataset, edition and version", t, func() {
		edition := "2018"
		version := "2"

		docs, err := setupResources(ids.DatasetAssociated, ids.EditionUnpublished, edition, ids.InstanceEditionConfirmed, ids.UniqueTimestamp, 1)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			t.FailNow()
		}

		count, err := neo4JStore.CreateInstanceNode(ids.InstanceEditionConfirmed)
		if err != nil {
			t.Errorf("failed to create neo4j instance node: [%v]\n error: [%v]\n", ids.InstanceEditionConfirmed, err)
			t.FailNow()
		}
		So(count, ShouldEqual, 1)

		Convey("When a PUT request to update meta data against the version resource", func() {
			Convey("Then version resource is updated and returns a status ok (200)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", ids.DatasetAssociated, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionMetaDataJSON)).
					Expect().Status(http.StatusOK)

				updatedVersion, err := mongo.GetVersion(cfg.MongoDB, "instances", "_id", ids.InstanceEditionConfirmed)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check version has been updated
				So(updatedVersion.ID, ShouldEqual, ids.InstanceEditionConfirmed)
				So(updatedVersion.ReleaseDate, ShouldEqual, "2018-11-11")
				So(len(*updatedVersion.UsageNotes), ShouldEqual, 2)

				So((*updatedVersion.UsageNotes)[0].Title, ShouldEqual, "Coefficients of variation")

				alert := mongo.Alert{
					Description: "All data entries (observations) for Plymouth have been updated",
					Date:        "2017-04-05",
					Type:        "Correction",
				}

				alertList := &[]mongo.Alert{alert}

				So(updatedVersion.Alerts, ShouldResemble, alertList)

				latestChange := mongo.LatestChange{
					Description: "change to the period frequency from quarterly to monthly",
					Name:        "Changes to the period frequency",
					Type:        "Summary of Changes",
				}

				latestChangesList := []mongo.LatestChange{latestChange}

				So(updatedVersion.LatestChanges, ShouldResemble, latestChangesList)

				So(updatedVersion.Links.Spatial.HRef, ShouldEqual, "http://ons.gov.uk/new-geography-list")

				// Check self link does not update - the only link that can be updated is `spatial`
				So(updatedVersion.Links.Self.HRef, ShouldNotEqual, "http://bogus/bad-link")

				temporal := mongo.TemporalFrequency{
					StartDate: "2014-11-11",
					EndDate:   "2017-11-11",
					Frequency: "monthly",
				}

				temporalList := []mongo.TemporalFrequency{temporal}

				So(updatedVersion.Temporal, ShouldResemble, temporalList)
			})
		})

		Convey("When a PUT request to update version resource with a collection id and state of associated", func() {
			Convey("Then the dataset and version resources are updated and returns a status ok (200)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", ids.DatasetAssociated, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToAssociatedJSON)).
					Expect().Status(http.StatusOK)

				updatedVersion, err := mongo.GetVersion(cfg.MongoDB, "instances", "_id", ids.InstanceEditionConfirmed)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check version has been updated
				So(updatedVersion.ID, ShouldEqual, ids.InstanceEditionConfirmed)
				So(updatedVersion.CollectionID, ShouldEqual, "45454545")
				So(updatedVersion.State, ShouldEqual, "associated")

				updatedDataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", ids.DatasetAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check dataset has been updated
				So(updatedDataset.ID, ShouldEqual, ids.DatasetAssociated)
				So(updatedDataset.Next.CollectionID, ShouldEqual, "45454545")
				So(updatedDataset.Next.State, ShouldEqual, "associated")
			})
		})

		Convey("When a PUT request to update version resource with a collection id and state of published", func() {
			Convey("Then the dataset, edition and version resources are updated and returns a status ok (200)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", ids.DatasetAssociated, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToPublishedWithCollectionIDJSON)).
					Expect().Status(http.StatusOK)

				updatedVersion, err := mongo.GetVersion(cfg.MongoDB, "instances", "_id", ids.InstanceEditionConfirmed)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check version has been updated, and CollectionID removed
				So(updatedVersion.ID, ShouldEqual, ids.InstanceEditionConfirmed)
				So(updatedVersion.CollectionID, ShouldBeEmpty)
				So(updatedVersion.State, ShouldEqual, "published")

				log.Debug("edition id", log.Data{"edition_id": ids.EditionUnpublished})

				updatedEdition, err := mongo.GetEdition(cfg.MongoDB, "editions", "_id", ids.EditionUnpublished)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check edition has been updated
				So(updatedEdition.ID, ShouldEqual, ids.EditionUnpublished)
				So(updatedEdition.Next.State, ShouldEqual, "published")
				So(updatedEdition.Current, ShouldNotBeNil)
				So(updatedEdition.Current.State, ShouldEqual, "published")

				updatedDataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", ids.DatasetAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check dataset has been updated, and CollectionID removed
				So(updatedDataset.ID, ShouldEqual, ids.DatasetAssociated)
				So(updatedDataset.Current.CollectionID, ShouldBeEmpty)
				So(updatedDataset.Current.State, ShouldEqual, "published")

				instanceProps, err := neo4JStore.GetInstanceProperties(ids.InstanceEditionConfirmed)
				if err != nil {
					log.ErrorC("failed to get properties from neo4j instance node", err, nil)
					t.FailNow()
				}

				So(instanceProps["is_published"], ShouldBeTrue)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}

		if err := neo4JStore.CleanUpInstance(ids.InstanceEditionConfirmed); err != nil {
			t.Errorf("failed to cleanup neo4j instances: [%v]\n error: [%v]\n", ids.InstanceEditionConfirmed, err)
			t.FailNow()
		}
	})

	Convey("Given an unpublished dataset, edition and a version that has been associated", t, func() {
		edition := "2018"
		version := "2"

		docs, err := setupResources(ids.DatasetAssociated, ids.EditionUnpublished, edition, ids.InstanceAssociated, ids.UniqueTimestamp, 2)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		count, err := neo4JStore.CreateInstanceNode(ids.InstanceAssociated)
		if err != nil {
			t.Errorf("failed to create neo4j instance node: [%v]\n error: [%v]\n", ids.InstanceAssociated, err)
			t.FailNow()
		}
		So(count, ShouldEqual, 1)

		// TODO Remove skipped tests when code has been refactored (and hence fixed)
		// 1 test skipped
		SkipConvey("When a PUT request to update version resource to remove collection id", func() {
			Convey("Then the dataset and version resources are updated accordingly and returns a status ok (200)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", ids.DatasetAssociated, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionFromAssociatedToEditionConfirmedJSON)).
					Expect().Status(http.StatusOK)

				updatedVersion, err := mongo.GetVersion(cfg.MongoDB, "instances", "_id", ids.InstanceAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check version has been updated
				So(updatedVersion.ID, ShouldEqual, ids.InstanceAssociated)
				So(updatedVersion.CollectionID, ShouldEqual, "")
				So(updatedVersion.State, ShouldEqual, "edition-confirmed")

				updatedDataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", ids.DatasetAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check dataset has been updated
				So(updatedDataset.ID, ShouldEqual, ids.DatasetAssociated)
				So(updatedDataset.Next.CollectionID, ShouldEqual, "")
				So(updatedDataset.Next.State, ShouldEqual, "edition-confirmed")
			})
		})

		Convey("When a PUT request to update version resource with a state of published", func() {
			Convey("Then the dataset, edition and version resources are updated and returns a status ok (200)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", ids.DatasetAssociated, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
					Expect().Status(http.StatusOK)

				updatedVersion, err := mongo.GetVersion(cfg.MongoDB, "instances", "_id", ids.InstanceAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check version has been updated
				So(updatedVersion.ID, ShouldEqual, ids.InstanceAssociated)
				So(updatedVersion.State, ShouldEqual, "published")

				updatedEdition, err := mongo.GetEdition(cfg.MongoDB, "editions", "_id", ids.EditionUnpublished)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check edition has been updated
				So(updatedEdition.ID, ShouldEqual, ids.EditionUnpublished)
				So(updatedEdition.Next.State, ShouldEqual, "published")
				So(updatedEdition.Current, ShouldNotBeNil)
				So(updatedEdition.Current.State, ShouldEqual, "published")
				So(updatedEdition.Current.Links.LatestVersion.ID, ShouldEqual, "2")
				So(updatedEdition.Current.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+ids.DatasetAssociated+"/editions/2018/versions/2")

				updatedDataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", ids.DatasetAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check dataset has been updated, next sub document should be copied across to current sub doc
				So(updatedDataset.ID, ShouldEqual, ids.DatasetAssociated)
				So(updatedDataset.Current.State, ShouldEqual, "published")
				So(updatedDataset.Next.State, ShouldEqual, "published") // Check next subdoc still exists
				So(updatedDataset, ShouldResemble, expectedDatasetResource(ids.DatasetAssociated, 0))
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}

		if err := neo4JStore.CleanUpInstance(ids.InstanceAssociated); err != nil {
			t.Errorf("failed to cleanup neo4j instances: [%v]\n error: [%v]\n", ids.InstanceAssociated, err)
			t.FailNow()
		}
	})

	Convey("Given a published dataset and edition, and a version that has been associated", t, func() {
		edition := "2017"
		version := "2"

		docs, err := setupResources(ids.DatasetPublished, ids.EditionPublished, edition, ids.InstanceAssociated, ids.UniqueTimestamp, 3)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		count, err := neo4JStore.CreateInstanceNode(ids.InstanceAssociated)
		if err != nil {
			t.Errorf("failed to create neo4j instance node: [%v]\n error: [%v]\n", ids.InstanceAssociated, err)
			t.Fail()
		}
		So(count, ShouldEqual, 1)

		Convey("When a PUT request to update version resource with a state of published", func() {
			Convey("Then the dataset, edition and version resources are updated and returns a status ok (200)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", ids.DatasetPublished, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
					Expect().Status(http.StatusOK)

				updatedVersion, err := mongo.GetVersion(cfg.MongoDB, "instances", "_id", ids.InstanceAssociated)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check version has been updated
				So(updatedVersion.ID, ShouldEqual, ids.InstanceAssociated)
				So(updatedVersion.State, ShouldEqual, "published")

				updatedEdition, err := mongo.GetEdition(cfg.MongoDB, "editions", "_id", ids.EditionPublished)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check edition has been updated
				So(updatedEdition.ID, ShouldEqual, ids.EditionPublished)
				So(updatedEdition.Next.State, ShouldEqual, "published")
				So(updatedEdition.Current, ShouldNotBeNil)
				So(updatedEdition.Current.State, ShouldEqual, "published")
				So(updatedEdition.Current.Links.LatestVersion.ID, ShouldEqual, "2")
				So(updatedEdition.Current.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+ids.DatasetPublished+"/editions/2017/versions/2")

				updatedDataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", ids.DatasetPublished)
				if err != nil {
					log.ErrorC("Unable to retrieve updated version document", err, nil)
					os.Exit(1)
				}

				// Check dataset has been updated, next sub document should be copied across to current sub doc
				So(updatedDataset.ID, ShouldEqual, ids.DatasetPublished)
				So(updatedDataset.Current.State, ShouldEqual, "published")
				So(updatedDataset.Next.State, ShouldEqual, "published") // Check next subdoc still exists
				So(updatedDataset, ShouldResemble, expectedDatasetResource(ids.DatasetPublished, 1))
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}

		if err := neo4JStore.CleanUpInstance(ids.InstanceAssociated); err != nil {
			t.Errorf("failed to cleanup neo4j instances: [%v]\n error: [%v]\n", ids.InstanceAssociated, err)
			t.FailNow()
		}
	})
}

func TestFailureToUpdateVersion(t *testing.T) {
	ts, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	edition := "2018"
	version := "2"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	// test for updating a version that has no dataset (bad request)
	Convey("Given an edition and a version of state associated exist for a dataset that does not exist in datastore", t, func() {

		docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 4)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised PUT request is made to update version resource", func() {
			Convey("Then fail to update resource and return a status of not found (404) with a message `dataset not found`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
					Expect().Status(http.StatusNotFound).Body().Contains("dataset not found")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}
	})

	// test for updating a version that has no edition (bad request)
	Convey("Given a dataset and a version both of state associated exist but the edition does not", t, func() {

		docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 5)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised PUT request is made to update version resource", func() {
			Convey("Then fail to update resource and return a status of not found (404) with a message `edition not found`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("edition not found")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}
	})

	// test for updating a version that does not exist (not found)
	Convey("Given a dataset and edition exist but the version for the dataset edition does not", t, func() {

		docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 6)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised PUT request is made to update version resource", func() {
			Convey("Then fail to update resource and return a status of not found (404) with a message `version not found`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
					Expect().Status(http.StatusNotFound).
					Body().Contains("version not found")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}
	})

	// test for updating a version with a state of `published` but request missing mandatory fields
	Convey("Given a published dataset and edition and an unpublished version exist", t, func() {

		Convey("with mandatory fields missing", func() {
			docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 9)
			if err != nil {
				log.ErrorC("Was unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an authorised PUT request is made to update version resource to published", func() {
				Convey("Then fail to update resource and return a status bad request (400) with the correct message", func() {

					datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
						WithHeader(florenceTokenName, florenceToken).
						WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
						Expect().Status(http.StatusBadRequest).
						Body().Contains("missing mandatory fields: [release_date Downloads.XLS.HRef Downloads.XLS.Size Downloads.CSV.HRef Downloads.CSV.Size Downloads.CSVW.HRef Downloads.CSVW.Size]")

				})
			})

			if err := mongo.Teardown(docs...); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("Was unable to remove test data", err, nil)
					os.Exit(1)
				}
			}
		})

		Convey("with invalid fields", func() {
			docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 10)
			if err != nil {
				log.ErrorC("Was unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an authorised PUT request is made to update version resource to published", func() {
				Convey("Then fail to update resource and return a status bad request (400) with the correct message", func() {

					datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
						WithHeader(florenceTokenName, florenceToken).
						WithBytes([]byte(validPUTUpdateVersionToPublishedJSON)).
						Expect().Status(http.StatusBadRequest).
						Body().Contains("invalid fields: [Downloads.XLS.Size not a number Downloads.CSV.Size not a number Downloads.CSVW.Size not a number]")

				})
			})

			if err := mongo.Teardown(docs...); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("Was unable to remove test data", err, nil)
					os.Exit(1)
				}
			}
		})
	})

	// test for bad request (invalid json)
	Convey("Given a dataset, edition and version do not exist", t, func() {
		Convey("When an authorised PUT request is made to update version resource with invalid json", func() {
			Convey("Then fail to update resource and return a status of bad request (400) with a message `failed to parse json body`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("failed to parse json body")

			})
		})
	})

	Convey("Given a dataset, edition exist and the version for the dataset edition has a state of `edition-confirmed`", t, func() {
		edition := "2018"
		version := "2"

		docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 7)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		// test for bad request when associating version (Missing mandatory fields)
		Convey("When an authorised PUT request is made to update version resource to a state of associated", func() {
			Convey("Then fail to update resource and return a status of bad request (400) with a message `Missing collection_id for association between version and a collection`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"state": "associated"}`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing collection_id for association between version and a collection")

			})
		})

		// test for unauthorised request to update version
		Convey("When an unauthorised PUT request is made to update version resource", func() {
			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					WithBytes([]byte(validPUTUpdateVersionMetaDataJSON)).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		// test for missing auth header when making a request to update version
		Convey("When a PUT request is made to update version resource without an authentication header", func() {
			Convey("Then fail to update resource and return a status unauthorized (401)", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithBytes([]byte(validPUTUpdateVersionMetaDataJSON)).
					Expect().Status(http.StatusUnauthorized)

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}
	})

	Convey("Given a dataset, edition and version for the dataset edition are published", t, func() {
		edition := "2017"
		version := "1"

		docs, err := setupResources(datasetID, editionID, edition, instanceID, ts.UniqueTimestamp, 8)
		if err != nil {
			log.ErrorC("Was unable to setup test data", err, nil)
			os.Exit(1)
		}

		// test for reverting state against a published version (forbidden)
		Convey("When an authorised PUT request is made to update version resource to a state of `edition-confirmed`", func() {
			Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update document, already published`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"state": "edition-confirmed"}`)).
					Expect().Status(http.StatusForbidden).
					Body().Contains("unable to update version as it has been published")

			})
		})

		// test for updating meta data against a published version (forbidden)
		Convey("When an authorised PUT request is made to update a published version resource with meta data", func() {
			Convey("Then fail to update resource and return a status of forbidden (403) with a message `Unable to update document, already published`", func() {

				datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetID, edition, version).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{"links":{"spatial":{"href": "http://ons.gov.uk/spatial-notes"}}}`)).
					Expect().Status(http.StatusForbidden).
					Body().Contains("unable to update version as it has been published")

			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to remove test data", err, nil)
				os.Exit(1)
			}
		}
	})
}

func setupResources(datasetID, editionID, edition, instanceID string, uniqueTimestamp bson.MongoTimestamp, setup int) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	publishedDatasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	associatedDatasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validAssociatedDatasetData(datasetID),
	}

	createdDatasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validCreatedDatasetData(datasetID),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	unpublishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidUnpublishedEditionData(datasetID, editionID, edition),
	}

	publishedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	associatedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	editionConfirmedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validEditionConfirmedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	editionConfirmedInstanceMissingMandatoryFieldsDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     editionConfirmedInstanceMissingMandatoryFields(datasetID, edition, instanceID, uniqueTimestamp),
	}

	editionConfirmedInstanceInvalidFieldsDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     editionConfirmedInstanceInvalidFields(datasetID, edition, instanceID, uniqueTimestamp),
	}

	switch setup {
	case 1:
		docs = append(docs, createdDatasetDoc, unpublishedEditionDoc, editionConfirmedInstanceDoc)
	case 2:
		docs = append(docs, associatedDatasetDoc, unpublishedEditionDoc, associatedInstanceDoc)
	case 3:
		docs = append(docs, publishedDatasetDoc, publishedEditionDoc, associatedInstanceDoc)
	case 4:
		docs = append(docs, unpublishedEditionDoc, associatedInstanceDoc)
	case 5:
		docs = append(docs, associatedDatasetDoc, associatedInstanceDoc)
	case 6:
		docs = append(docs, associatedDatasetDoc, unpublishedEditionDoc)
	case 7:
		docs = append(docs, publishedDatasetDoc, publishedEditionDoc, editionConfirmedInstanceDoc)
	case 8:
		docs = append(docs, publishedDatasetDoc, publishedEditionDoc, publishedInstanceDoc)
	case 9:
		docs = append(docs, publishedDatasetDoc, publishedEditionDoc, editionConfirmedInstanceMissingMandatoryFieldsDoc)
	case 10:
		docs = append(docs, publishedDatasetDoc, publishedEditionDoc, editionConfirmedInstanceInvalidFieldsDoc)
	default:
		errMsg := fmt.Errorf("Failed to pick a valid setup value")
		log.Error(errMsg, log.Data{"setup": setup})
		return nil, errMsg
	}

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		return nil, err
	}

	return docs, nil
}

func expectedDatasetResource(datasetID string, resource int) mongo.DatasetUpdate {

	nationalStatistic := true

	doc := mongo.Dataset{
		Contacts:    []mongo.ContactDetails{contact},
		Description: "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
		Keywords:    []string{"cpi", "boy"},
		Links: &mongo.DatasetLinks{
			AccessRights: &mongo.LinkObject{
				HRef: "http://ons.gov.uk/accessrights",
			},
			Editions: &mongo.LinkObject{
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			},
			LatestVersion: &mongo.LinkObject{
				ID:   "2",
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2018/versions/2",
			},
			Self: &mongo.LinkObject{
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID,
			},
		},
		Methodologies:     []mongo.GeneralDetails{methodology},
		NationalStatistic: &nationalStatistic,
		NextRelease:       "2018-10-10",
		Publications:      []mongo.GeneralDetails{publication},
		Publisher: &mongo.Publisher{
			Name: "Automation Tester",
			Type: "publisher",
			HRef: "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		},
		QMI: &mongo.GeneralDetails{
			Description: "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
			HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
			Title:       "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
		},
		RelatedDatasets:  []mongo.GeneralDetails{relatedDatasets},
		ReleaseFrequency: "Monthly",
		State:            "published",
		Theme:            "Goods and services",
		Title:            "CPI",
		UnitOfMeasure:    "Pounds Sterling",
		URI:              "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
	}

	if resource == 1 {
		doc.License = "ONS license"
		doc.Links.LatestVersion.HRef = cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/2"
	}

	dataset := mongo.DatasetUpdate{
		ID:      datasetID,
		Current: &doc,
		Next:    &doc,
	}

	return dataset
}
