package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestGetVersions_ReturnsListOfVersions(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	Convey("Given a dataset edition has a published and unpublished version", t, func() {

		docs, err := setUpDatasetEditionVersions(ids.DatasetPublished, ids.EditionPublished, edition, ids.InstancePublished, ids.InstanceAssociated, ids.UniqueTimestamp)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

		Convey("When user is authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, edition).
				WithHeader(florenceTokenName, florenceToken).Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then response contains a list of all versions of the dataset edition", func() {
				response.Value("items").Array().Length().Equal(2)
				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {

					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstancePublished {
						// check the published test version document has the expected returned fields and values
						response.Value("items").Array().Element(i).Object().Value("id").Equal(ids.InstancePublished)
						checkVersionResponse(ids.DatasetPublished, ids.EditionPublished, ids.InstancePublished, edition, response.Value("items").Array().Element(i).Object())
						checkNeitherPublicOrPrivateLinksExistInResponse(response.Value("items").Array().Element(i).Object().Value("downloads").Object())
					}

					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstanceAssociated {
						response.Value("items").Array().Element(i).Object().Value("id").Equal(ids.InstanceAssociated)
						response.Value("items").Array().Element(i).Object().Value("state").Equal("associated")
						checkNeitherPublicOrPrivateLinksExistInResponse(response.Value("items").Array().Element(i).Object().Value("downloads").Object())
					}
				}
			})
		})

		Convey("When the caller of the request is the download service", func() {
			headers := make(map[string]string)
			headers[downloadServiceTokenName] = downloadServiceToken
			headers[downloadServiceAuthTokenName] = downloadServiceAuthToken
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, edition).
				WithHeaders(headers).Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then response contains a list of all versions of the dataset edition with there respective public and private download links", func() {
				response.Value("items").Array().Length().Equal(2)
				for i := 0; i < len(response.Value("items").Array().Iter()); i++ {

					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstancePublished {
						// check the published test version document has the expected returned fields and values
						response.Value("items").Array().Element(i).Object().Value("id").Equal(ids.InstancePublished)
						checkVersionResponse(ids.DatasetPublished, ids.EditionPublished, ids.InstancePublished, edition, response.Value("items").Array().Element(i).Object())
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csv").Object().Value("private").String().Match("private/myfile.csv")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csv").Object().Value("public").String().Match("public/myfile.csv")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csvw").Object().Value("private").String().Match("private/myfile.csv-metadata.json")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csvw").Object().Value("public").String().Match("public/myfile.csv-metadata.json")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("xls").Object().Value("private").String().Match("private/myfile.xls")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("xls").Object().Value("public").String().Match("public/myfile.xls")
					}

					if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == ids.InstanceAssociated {
						response.Value("items").Array().Element(i).Object().Value("id").Equal(ids.InstanceAssociated)
						response.Value("items").Array().Element(i).Object().Value("state").Equal("associated")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csv").Object().Value("private").String().Match("private/myfile.csv")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csv").Object().Value("public").String().Match("public/myfile.csv")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csvw").Object().Value("private").String().Match("private/myfile.csv-metadata.json")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("csvw").Object().Value("public").String().Match("public/myfile.csv-metadata.json")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("xls").Object().Value("private").String().Match("private/myfile.xls")
						response.Value("items").Array().Element(i).Object().Value("downloads").Object().Value("xls").Object().Value("public").String().Match("public/myfile.xls")
					}
				}
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestGetVersions_Failed(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	publishedEdition := "2017"
	unpublishedEdition := "2018"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      ids.DatasetPublished,
		Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionPublished,
		Update:     ValidPublishedEditionData(ids.DatasetPublished, ids.EditionPublished, publishedEdition),
	}

	publishedInstance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      ids.InstancePublished,
		Update:     validPublishedInstanceData(ids.DatasetPublished, publishedEdition, ids.InstancePublished, ids.UniqueTimestamp),
	}

	unpublishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionUnpublished,
		Update:     validUnpublishedEditionData(ids.DatasetPublished, ids.EditionUnpublished, unpublishedEdition),
	}

	Convey("Given the dataset and subsequently the edition does not exist", t, func() {
		Convey("When an authenticated request is made to get a list of versions of the dataset edition", func() {
			Convey("Then return status not found (404) with message `dataset not found`", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, publishedEdition).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("dataset not found")
			})
		})
	})

	Convey("Given the dataset exists", t, func() {
		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("but the edition does not", func() {
			Convey("When an authenticated request is made to get a list of versions of the dataset edition", func() {
				Convey("Then return status not found (404) with message `edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, publishedEdition).WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).Body().Contains("edition not found")
				})
			})
		})

		Convey("and the edition does exist but there are no versions", func() {
			if err := mongo.Setup(publishedEditionDoc); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an authenticated request is made to get a list of versions of the dataset edition", func() {
				Convey("Then return status not found (404) with message `version not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, publishedEdition).WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).Body().Contains("version not found")
				})
			})
		})

		if err := mongo.Teardown(publishedDataset, publishedEditionDoc); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	// Make sure an unauthorised user cannot find the dataset
	Convey("Given an unpublished dataset exists", t, func() {
		unpublishedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      ids.DatasetAssociated,
			Update:     validAssociatedDatasetData(ids.DatasetAssociated),
		}

		if err := mongo.Setup(unpublishedDataset); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
			Convey("Then return status unauthorised (401)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, publishedEdition).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(unpublishedDataset); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a published dataset exists", t, func() {
		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("but only an unpublished edition exists", func() {
			if err := mongo.Setup(unpublishedEditionDoc); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
				Convey("Then return status unauthorized (401)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, unpublishedEdition).
						Expect().Status(http.StatusUnauthorized)
				})
			})

			if err := mongo.Teardown(unpublishedEditionDoc); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})

		Convey("and a published edition exists", func() {
			if err := mongo.Setup(publishedEditionDoc); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("but only unpublished versions exist for the dataset edition", func() {

				unpublishedInstance := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: "instances",
					Key:        "_id",
					Value:      ids.InstanceAssociated,
					Update:     validAssociatedInstanceData(ids.DatasetPublished, publishedEdition, ids.InstanceAssociated, ids.UniqueTimestamp),
				}

				if err := mongo.Setup(unpublishedInstance); err != nil {
					log.ErrorC("Unable to setup test data", err, nil)
					os.Exit(1)
				}

				Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
					Convey("Then return status unauthorized (401)", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, publishedEdition).
							Expect().Status(http.StatusUnauthorized)
					})
				})

				if err := mongo.Teardown(unpublishedInstance); err != nil {
					log.ErrorC("Unable to remove test data from mongo db", err, nil)
					os.Exit(1)
				}
			})

			Convey("and a published version exists", func() {
				if err := mongo.Setup(publishedInstance); err != nil {
					log.ErrorC("Unable to setup test data", err, nil)
					os.Exit(1)
				}

				Convey("When an unauthenticated request is made to get a list of versions of the dataset edition", func() {
					Convey("Then return status unauthorized (401)", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions", ids.DatasetPublished, publishedEdition).WithHeader(florenceTokenName, unauthorisedAuthToken).
							Expect().Status(http.StatusUnauthorized)
					})
				})

				if err := mongo.Teardown(publishedInstance, publishedEditionDoc, publishedDataset); err != nil {
					log.ErrorC("Unable to remove test data from mongo db", err, nil)
					os.Exit(1)
				}
			})
		})
	})
}

func checkVersionResponse(datasetID, editionID, instanceID, edition string, response *httpexpect.Object) {
	response.Value("id").Equal(instanceID)
	response.Value("alerts").Array().Element(0).Object().Value("date").String().Equal("2017-12-10")
	response.Value("alerts").Array().Element(0).Object().Value("description").String().Equal("A correction to an observation for males of age 25, previously 11 now changed to 12")
	response.Value("alerts").Array().Element(0).Object().Value("type").String().Equal("Correction")
	response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("An aggregate of the data")
	response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD$")
	response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("508064B3-A808-449B-9041-EA3A2F72CFAD")
	response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("aggregate")
	response.Value("downloads").Object().Value("csv").Object().Value("href").String().Match("/aws/census-2017-1-csv$")
	response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10")
	response.Value("downloads").Object().Value("csvw").Object().Value("href").String().Match("/aws/census-2017-1-csv-metadata.json$")
	response.Value("downloads").Object().Value("csvw").Object().Value("size").Equal("10")
	response.Value("downloads").Object().Value("xls").Object().Value("href").String().Match("/aws/census-2017-1-xls$")
	response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24")
	response.Value("edition").Equal("2017")
	response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
	response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
	response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
	response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
	response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions$")
	response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
	response.Value("links").Object().Value("edition").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
	response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
	response.Value("release_date").Equal("2017-12-12") // TODO Should be isodate
	response.Value("state").Equal("published")
	response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
	response.Value("version").Equal(1)
}

func checkNeitherPublicOrPrivateLinksExistInResponse(response *httpexpect.Object) {
	response.Value("csv").Object().NotContainsKey("public")
	response.Value("csv").Object().NotContainsKey("private")
	response.Value("csvw").Object().NotContainsKey("public")
	response.Value("csvw").Object().NotContainsKey("private")
	response.Value("xls").Object().NotContainsKey("public")
	response.Value("xls").Object().NotContainsKey("private")
}

func setUpDatasetEditionVersions(datasetID, editionID, edition, instanceID, unpublishedInstanceID string, uniqueTimestamp bson.MongoTimestamp) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	unpublishedInstanceDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID, uniqueTimestamp),
	}

	docs = append(docs, datasetDoc, editionDoc, instanceDoc, unpublishedInstanceDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
