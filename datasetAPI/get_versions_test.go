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

func TestGetVersions_ReturnsListOfVersions(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	Convey("Given a dataset edition has a published and unpublished version", t, func() {

		d, err := setUpDatasetEditionVersions(datasetID, editionID, edition, instanceID, unpublishedInstanceID)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

		Convey("When user is authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
				WithHeader("Internal-Token", internalTokenID).Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then response contains a list of all versions of the dataset edition", func() {
				response.Value("items").Array().Length().Equal(2)
				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {

					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == instanceID {
						// check the published test version document has the expected returned fields and values
						response.Value("items").Array().Element(i).Object().Value("id").Equal(instanceID)
						checkVersionResponse(datasetID, editionID, instanceID, edition, response.Value("items").Array().Element(i).Object())
					}

					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == unpublishedInstanceID {
						response.Value("items").Array().Element(i).Object().Value("id").Equal(unpublishedInstanceID)
						response.Value("items").Array().Element(i).Object().Value("state").Equal("associated")
					}
				}
			})
		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
				Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then response contains a list of only published versions of the dataset edition", func() {
				response.Value("items").Array().Length().Equal(1)
				checkVersionResponse(datasetID, editionID, instanceID, edition, response.Value("items").Array().Element(0).Object())
			})
		})

		if err := mongo.TeardownMany(d); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestGetVersions_Failed(t *testing.T) {

	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2018"

	unpublishedDatasetID := uuid.NewV4().String()
	unpublishedEditionID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given the dataset and subsequently the edition does not exist", t, func() {
		Convey("When an authenticated request is made to get a list of versions of the dataset edition", func() {
			Convey("Then return status bad request (400)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Dataset not found\n")
			})
		})
	})

	Convey("Given the dataset exists", t, func() {
		update := validPublishedDatasetData(datasetID)
		if err := mongo.Setup(database, collection, "_id", datasetID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("but the edition does not", func() {
			Convey("When an authenticated request is made to get a list of versions of the dataset edition", func() {
				Convey("Then return status bad request (400)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader(internalToken, internalTokenID).
						Expect().Status(http.StatusBadRequest).Body().Contains("Edition not found\n")
				})
			})
		})

		Convey("and the edition does exist but there are no versions", func() {
			update := validPublishedEditionData(datasetID, editionID, edition)
			if err := mongo.Setup(database, "editions", "_id", editionID, update); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an authenticated request is made to get a list of versions of the dataset edition", func() {
				Convey("Then return status not found (404)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).WithHeader(internalToken, internalTokenID).
						Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")
				})
			})
			if err := mongo.Teardown(database, "editions", "_id", editionID); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})

		if err := mongo.Teardown(database, collection, "_id", datasetID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	// Make sure an unauthorised user cannot find the dataset
	Convey("Given an unpublished dataset exists", t, func() {
		update := validAssociatedDatasetData(unpublishedDatasetID)
		if err := mongo.Setup(database, collection, "_id", unpublishedDatasetID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
			Convey("Then return status bad request (400)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
					Expect().Status(http.StatusBadRequest).Body().Contains("Dataset not found\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", unpublishedDatasetID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a published dataset exists", t, func() {
		update := validPublishedDatasetData(datasetID)
		if err := mongo.Setup(database, collection, "_id", datasetID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("but only an unpublished edition exists", func() {
			update := validUnpublishedEditionData(datasetID, unpublishedEditionID, edition)
			if err := mongo.Setup(database, "editions", "_id", unpublishedEditionID, update); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
				Convey("Then return status bad request (400)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
						Expect().Status(http.StatusBadRequest).Body().Contains("Edition not found\n")
				})
			})

			if err := mongo.Teardown(database, "editions", "_id", unpublishedEditionID); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})

		Convey("and a published edition exists", func() {
			update := validPublishedEditionData(datasetID, editionID, edition)
			if err := mongo.Setup(database, "editions", "_id", editionID, update); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("but only unpublished versions exist for the dataset edition", func() {
				update := validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID)
				if err := mongo.Setup(database, "instances", "_id", unpublishedInstanceID, update); err != nil {
					log.ErrorC("Unable to setup test data", err, nil)
					os.Exit(1)
				}

				Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
					Convey("Then return status not found (404)", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", datasetID, edition).
							Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")
					})
				})

				if err := mongo.Teardown(database, "instances", "_id", unpublishedInstanceID); err != nil {
					log.ErrorC("Unable to remove test data from mongo db", err, nil)
					os.Exit(1)
				}
			})

			if err := mongo.Teardown(database, "editions", "_id", editionID); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})

		if err := mongo.Teardown(database, collection, "_id", datasetID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func checkVersionResponse(datasetID, editionID, instanceID, edition string, response *httpexpect.Object) {
	response.Value("id").Equal(instanceID)
	response.Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
	response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
	response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("(.+)/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
	response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
	response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
	response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-1-csv$")
	response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
	response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-1-xls$")
	response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
	response.Value("edition").Equal("2017")
	response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
	response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions$")
	response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
	response.Value("links").Object().Value("edition").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
	response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
	response.Value("release_date").Equal("2017-12-12") // TODO Should be isodate
	response.Value("state").Equal("published")
	response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
	response.Value("version").Equal(1)
}

func setUpDatasetEditionVersions(datasetID, editionID, edition, instanceID, unpublishedInstanceID string) (*mongo.ManyDocs, error) {
	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	unpublishedInstanceDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, instanceDoc, unpublishedInstanceDoc) //, instanceDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.SetupMany(d); err != nil {
		return nil, err
	}

	return d, nil
}