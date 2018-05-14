package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestSuccessfullyGetVersionOfADatasetEdition(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	authHeaders := make(map[string]string)
	authHeaders[downloadServiceAuthTokenName] = downloadServiceAuthToken
	authHeaders[florenceTokenName] = florenceToken

	Convey("Given a published and unpublished version for a dataset edition exists", t, func() {
		docs, err := setupPublishedAndUnpublishedVersions(datasetID, editionID, edition, instanceID, unpublishedInstanceID)
		if err != nil {
			log.ErrorC("Failed to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated request is made to get the unpublished version", func() {
			Convey("Then the response body contains the expected version", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2", datasetID, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("id").Equal(unpublishedInstanceID)
				response.Value("collection_id").Equal("208064B3-A808-449B-9041-EA3A2F72CFAB")
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("downloads").Object().Value("csv").Object().Value("href").String().Match("/aws/census-2017-2-csv$")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10")
				response.Value("downloads").Object().Value("xls").Object().Value("href").String().Match("/aws/census-2017-2-xls$")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24")
				response.Value("edition").Equal(edition)
				response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
				response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
				response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
				response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
				response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
				response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2/dimensions$")
				response.Value("links").Object().Value("edition").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
				response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
				response.Value("release_date").Equal("2017-12-12")
				response.Value("state").Equal("associated")
				response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
				response.Value("version").Equal(2)
			})
		})

		Convey("When an authenticated request including a valid download service token is made to get the published version", func() {
			Convey("Then the response body contains the expected downloads object in version", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2", datasetID, edition).
					WithHeaders(authHeaders).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("downloads").Object().Value("csv").Object().Value("public").String().Equal("https://s3-eu-west-1.amazon.com/public/myfile.csv")
				response.Value("downloads").Object().Value("csv").Object().Value("private").String().Equal("s3://private/myfile.csv")
				response.Value("downloads").Object().Value("xls").Object().Value("public").String().Equal("https://s3-eu-west-1.amazon.com/public/myfile.xls")
				response.Value("downloads").Object().Value("xls").Object().Value("private").String().Equal("s3://private/myfile.xls")
			})
		})

		Convey("When an authenticated request is made to get the published version", func() {
			Convey("Then the response body contains the expected version", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("alerts").Array().Element(0).Object().Value("date").String().Equal("2017-12-10")
				response.Value("alerts").Array().Element(0).Object().Value("description").String().Equal("A correction to an observation for males of age 25, previously 11 now changed to 12")
				response.Value("alerts").Array().Element(0).Object().Value("type").String().Equal("Correction")
				response.Value("id").Equal(instanceID)
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("downloads").Object().Value("csv").Object().Value("href").String().Match("/aws/census-2017-1-csv$")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10")
				response.Value("downloads").Object().Value("xls").Object().Value("href").String().Match("/aws/census-2017-1-xls$")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24")
				response.Value("edition").Equal(edition)
				response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
				response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
				response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
				response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
				response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
				response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions$")
				response.Value("links").Object().Value("edition").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
				response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
				response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
				response.Value("release_date").Equal("2017-12-12")
				response.Value("state").Equal("published")
				response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
				response.Value("version").Equal(1)
			})
		})

		Convey("When a request including a valid download service token is made to get the published version", func() {
			Convey("Then the response body contains the expected downloads object in version", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
					WithHeaders(authHeaders).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("downloads").Object().Value("csv").Object().Value("public").String().Equal("https://s3-eu-west-1.amazon.com/public/myfile.csv")
				response.Value("downloads").Object().Value("csv").Object().Value("private").String().Equal("s3://private/myfile.csv")
				response.Value("downloads").Object().Value("xls").Object().Value("public").String().Equal("https://s3-eu-west-1.amazon.com/public/myfile.xls")
				response.Value("downloads").Object().Value("xls").Object().Value("private").String().Equal("s3://private/myfile.xls")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetVersionOfADatasetEdition(t *testing.T) {

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	unpublishedDatasetID := uuid.NewV4().String()
	unpublishedEditionID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	unpublishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      unpublishedDatasetID,
		Update:     validAssociatedDatasetData(unpublishedDatasetID),
	}

	publishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	unpublishedEdition := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      unpublishedEditionID,
		Update:     validUnpublishedEditionData(unpublishedDatasetID, unpublishedEditionID, edition),
	}

	unpublishedInstance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, editionID, unpublishedInstanceID),
	}

	Convey("Given an unpublished dataset, edition and version exists", t, func() {
		var docs []*mongo.Doc
		docs = append(docs, unpublishedDataset, unpublishedEdition, unpublishedInstance)

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When a request to get version of the dataset edition and an invalid token is set", func() {
			Convey("Then return status unauthorized (401)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		}
	})

	Convey("Given the dataset, edition and version do not exist", t, func() {
		Convey("When an authorised request to get the version of the dataset edition", func() {
			Convey("Then return status not found (404) with message `Dataset not found`", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("Dataset not found")
			})
		})
	})

	Convey("Given an unpublished dataset exist", t, func() {
		if err := mongo.Setup(unpublishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("but an edition and version do not exist", func() {
			Convey("When a request to get the version of the dataset edition", func() {
				Convey("Then return status not found (404) with message `Edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", unpublishedDatasetID, edition).WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).Body().Contains("Edition not found")
				})
			})
		})

		Convey("and an unpublished edition exist", func() {
			if err := mongo.Setup(unpublishedEdition); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("but a version does not exist", func() {
				Convey("When a request to get the version of the dataset edition", func() {
					Convey("Then return status bad request (404) with message `Version not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", unpublishedDatasetID, edition).WithHeader(florenceTokenName, florenceToken).
							Expect().Status(http.StatusNotFound).Body().Contains("Version not found")
					})
				})
			})
		})
	})

	if err := mongo.Teardown(unpublishedDataset, unpublishedEdition); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	// Similar tests for unauthorised requests
	Convey("Given an unpublished dataset", t, func() {
		if err := mongo.Setup(unpublishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthorised request to get the version of the dataset edition", func() {
			Convey("Then return status unauthorized (401)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(unpublishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a published dataset", t, func() {
		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("and an unpublished edition", func() {
			unpublishedEdition.Update = validUnpublishedEditionData(datasetID, unpublishedEditionID, edition)
			if err := mongo.Setup(unpublishedEdition); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthorised request to get the version of the dataset edition", func() {
				Convey("Then return status unauthorized (401)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
						Expect().Status(http.StatusUnauthorized)
				})
			})

			if err := mongo.Teardown(unpublishedEdition); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		})

		Convey("and a published edition but an unpublished version", func() {
			publishedEdition := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "editions",
				Key:        "_id",
				Value:      editionID,
				Update:     ValidPublishedEditionData(datasetID, editionID, edition),
			}

			if err := mongo.Setup(publishedEdition, unpublishedInstance); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthorised request to get the version of the dataset edition", func() {
				Convey("Then return status unauthorized (401)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1", datasetID, edition).
						Expect().Status(http.StatusUnauthorized)
				})
			})

			if err := mongo.Teardown(publishedEdition, unpublishedInstance, publishedDataset); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		})
	})
}

func setupPublishedAndUnpublishedVersions(datasetID, editionID, edition, instanceID, unpublishedInstanceID string) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	publishedVersionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	unpublishedVersionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID),
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, publishedVersionDoc, unpublishedVersionDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
