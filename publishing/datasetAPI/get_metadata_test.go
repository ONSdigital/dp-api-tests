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

func TestSuccessfullyGetMetadataRelevantToVersion(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published and unpublished version", t, func() {
		docs, err := setupMetadataDocs(ids.DatasetPublished, ids.EditionPublished, edition, ids.InstancePublished, ids.InstanceAssociated, ids.UniqueTimestamp)
		if err != nil {
			log.ErrorC("Failed to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated request is made to get the unpublished version", func() {
			Convey("Then the response body contains the expected metadata", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/metadata", ids.DatasetPublished, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
				response.Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
				response.Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
				response.Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("An aggregate of the data")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("508064B3-A808-449B-9041-EA3A2F72CFAD")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("aggregate")
				response.Value("distribution").Array().Element(1).Equal("csv")
				response.Value("distribution").Array().Element(2).Equal("csvw")
				response.Value("distribution").Array().Element(3).Equal("xls")
				response.Value("downloads").Object().Value("csv").Object().Value("href").String().Match("/aws/census-2017-2-csv$")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10")
				response.Value("downloads").Object().Value("csvw").Object().Value("href").String().Match("/aws/census-2017-2-csv-metadata.json$")
				response.Value("downloads").Object().Value("csvw").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("csvw").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("csvw").Object().Value("size").Equal("10")
				response.Value("downloads").Object().Value("xls").Object().Value("href").String().Match("/aws/census-2017-2-xls$")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24")
				response.Value("keywords").Array().Element(0).Equal("cpi")
				response.Value("keywords").Array().Element(1).Equal("boy")
				response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
				response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
				response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
				response.Value("license").Equal("ONS license")
				response.Value("links").Object().Value("access_rights").Object().Value("href").Equal("http://ons.gov.uk/accessrights")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("2")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + edition + "/versions/2/metadata$")
				response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
				response.Value("methodologies").Array().Element(0).Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.")
				response.Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("methodologies").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("national_statistic").Equal(true)
				response.Value("next_release").Equal("2018-10-10")
				response.Value("publications").Array().Element(0).Object().Value("description").Equal("Price indices, percentage changes and weights for the different measures of consumer price inflation.")
				response.Value("publications").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
				response.Value("publications").Array().Element(0).Object().Value("title").Equal("UK consumer price inflation: August 2017")
				response.Value("publisher").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
				response.Value("publisher").Object().Value("name").Equal("Automation Tester")
				response.Value("publisher").Object().Value("type").Equal("publisher")
				response.Value("qmi").Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall")
				response.Value("qmi").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("qmi").Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("related_datasets").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices")
				response.Value("related_datasets").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation time series dataset")
				response.Value("release_date").Equal("2017-12-12")
				response.Value("release_frequency").Equal("Monthly")
				response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
				response.Value("theme").Equal("Goods and services")
				response.Value("title").Equal("CPI")
				response.Value("unit_of_measure").Equal("Pounds Sterling")
				response.Value("uri").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation")
			})
		})

		Convey("When an authenticated request is made to get the metadata relevant to a published version ", func() {
			Convey("Then the response body contains the expected metadata", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetPublished, edition).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
				response.Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
				response.Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
				response.Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("An aggregate of the data")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("508064B3-A808-449B-9041-EA3A2F72CFAD")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("aggregate")
				response.Value("distribution").Array().Element(1).Equal("csv")
				response.Value("distribution").Array().Element(2).Equal("csvw")
				response.Value("distribution").Array().Element(3).Equal("xls")
				response.Value("downloads").Object().Value("csv").Object().Value("href").String().Match("/aws/census-2017-1-csv$")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10")
				response.Value("downloads").Object().Value("csvw").Object().Value("href").String().Match("/aws/census-2017-1-csv-metadata.json$")
				response.Value("downloads").Object().Value("csvw").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("csvw").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("csvw").Object().Value("size").Equal("10")
				response.Value("downloads").Object().Value("xls").Object().Value("href").String().Match("/aws/census-2017-1-xls$")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24")
				response.Value("keywords").Array().Element(0).Equal("cpi")
				response.Value("keywords").Array().Element(1).Equal("boy")
				response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
				response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
				response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
				response.Value("license").Equal("ONS license")
				response.Value("links").Object().Value("access_rights").Object().Value("href").Equal("http://ons.gov.uk/accessrights")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + edition + "/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + ids.DatasetPublished + "/editions/" + edition + "/versions/1/metadata$")
				response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
				response.Value("methodologies").Array().Element(0).Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.")
				response.Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("methodologies").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("national_statistic").Equal(true)
				response.Value("next_release").Equal("2018-10-10")
				response.Value("publications").Array().Element(0).Object().Value("description").Equal("Price indices, percentage changes and weights for the different measures of consumer price inflation.")
				response.Value("publications").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
				response.Value("publications").Array().Element(0).Object().Value("title").Equal("UK consumer price inflation: August 2017")
				response.Value("publisher").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
				response.Value("publisher").Object().Value("name").Equal("Automation Tester")
				response.Value("publisher").Object().Value("type").Equal("publisher")
				response.Value("qmi").Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall")
				response.Value("qmi").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("qmi").Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("related_datasets").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices")
				response.Value("related_datasets").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation time series dataset")
				response.Value("release_date").Equal("2017-12-12")
				response.Value("release_frequency").Equal("Monthly")
				response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
				response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
				response.Value("theme").Equal("Goods and services")
				response.Value("title").Equal("CPI")
				response.Value("unit_of_measure").Equal("Pounds Sterling")
				response.Value("uri").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetMetadataRelevantToVersion(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	unpublishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      ids.DatasetAssociated,
		Update:     validAssociatedDatasetData(ids.DatasetAssociated),
	}

	unpublishedEdition := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionUnpublished,
		Update:     ValidUnpublishedEditionData(ids.DatasetAssociated, ids.EditionUnpublished, edition),
	}

	Convey("Given the dataset, edition and version do not exist", t, func() {
		Convey("When an authorised request to get the metadata relevant to a version", func() {
			Convey("Then return status not found (404) with message `dataset not found`", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetAssociated, edition).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("dataset not found")
			})
		})
	})

	Convey("Given an unpublished dataset exist", t, func() {
		if err := mongo.Setup(unpublishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("but an edition and version do not exist", func() {
			Convey("When a request to get the metadata relevant to a version", func() {
				Convey("Then return status not found (404) with message `edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetAssociated, edition).WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).Body().Contains("edition not found")
				})
			})
		})

		Convey("and an unpublished edition exist", func() {
			if err := mongo.Setup(unpublishedEdition); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("but a version does not exist", func() {
				Convey("When a request to get the metadata relevant to a version", func() {
					Convey("Then return status bad request (404) with message `version not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetAssociated, edition).WithHeader(florenceTokenName, florenceToken).
							Expect().Status(http.StatusNotFound).Body().Contains("version not found")
					})
				})
			})
		})

		if err := mongo.Teardown(unpublishedDataset, unpublishedEdition); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	// Similar tests for unauthorised requests
	Convey("Given an unpublished dataset", t, func() {
		if err := mongo.Setup(unpublishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthorised request to get the metadate relevant to a version", func() {
			Convey("Then return status unauthorized (401)", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetAssociated, edition).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(unpublishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a published dataset", t, func() {
		publishedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      ids.DatasetPublished,
			Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
		}

		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("and an unpublished edition", func() {
			unpublishedEdition.Update = ValidUnpublishedEditionData(ids.DatasetPublished, ids.EditionUnpublished, edition)
			if err := mongo.Setup(unpublishedEdition); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthorised request to get the metadata relevant to a version", func() {
				Convey("Then return status unauthorized (401)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetPublished, edition).
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
				Value:      ids.EditionPublished,
				Update:     ValidPublishedEditionData(ids.DatasetPublished, ids.EditionPublished, edition),
			}

			associatedInstance := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      ids.InstanceAssociated,
				Update:     validAssociatedInstanceData(ids.DatasetPublished, ids.EditionPublished, ids.InstanceAssociated, ids.UniqueTimestamp),
			}

			if err := mongo.Setup(publishedEdition, associatedInstance); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthorised request to get the metadata relevant to a version", func() {
				Convey("Then return status unauthorized (401)", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", ids.DatasetPublished, edition).
						Expect().Status(http.StatusUnauthorized)
				})
			})

			if err := mongo.Teardown(publishedEdition, associatedInstance); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		})

		if err := mongo.Teardown(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})
}

func setupMetadataDocs(datasetID, editionID, edition, instanceID, unpublishedInstanceID string, uniqueTimestamp bson.MongoTimestamp) ([]*mongo.Doc, error) {
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
		Update:     validPublishedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	unpublishedVersionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID, uniqueTimestamp),
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, publishedVersionDoc, unpublishedVersionDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
