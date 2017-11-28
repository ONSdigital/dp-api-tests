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

func TestSuccessfullyGetMetadataRelevantToVersion(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published and unpublished version", t, func() {
		d, err := setupMetadataDocs(datasetID, editionID, edition, instanceID, unpublishedInstanceID)
		if err != nil {
			log.ErrorC("Failed to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated request is made to get the unpublished version", func() {
			Convey("Then the response body contains the expected metadata", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/metadata", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
				response.Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
				response.Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
				response.Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("(.+)/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("distribution").Array().Element(0).Equal("json")
				response.Value("distribution").Array().Element(1).Equal("csv")
				response.Value("distribution").Array().Element(2).Equal("xls")
				response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-2-csv$")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
				response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-2-xls$")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
				response.Value("keywords").Array().Element(0).Equal("cpi")
				response.Value("keywords").Array().Element(1).Equal("boy")
				response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
				response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
				response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
				response.Value("license").Equal("ONS license")
				response.Value("links").Object().Value("access_rights").Object().Value("href").Equal("http://ons.gov.uk/accessrights")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("2")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2/metadata$")
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

		Convey("When an unauthenticated request is made to get the metadata relevant to a published version ", func() {
			Convey("Then the response body contains the expected metadata", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", datasetID, edition).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
				response.Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
				response.Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
				response.Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("(.+)/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("distribution").Array().Element(0).Equal("json")
				response.Value("distribution").Array().Element(1).Equal("csv")
				response.Value("distribution").Array().Element(2).Equal("xls")
				response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-1-csv$")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
				response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-1-xls$")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
				response.Value("keywords").Array().Element(0).Equal("cpi")
				response.Value("keywords").Array().Element(1).Equal("boy")
				response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
				response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
				response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
				response.Value("license").Equal("ONS license")
				response.Value("links").Object().Value("access_rights").Object().Value("href").Equal("http://ons.gov.uk/accessrights")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/metadata$")
				response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
				response.Value("methodologies").Array().Element(0).Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.")
				response.Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("methodologies").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("national_statistic").Equal(true)
				response.Value("next_release").Equal("2017-10-10")
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

		if err := mongo.TeardownMany(d); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetMetadataRelevantToVersion(t *testing.T) {

	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	unpublishedDatasetID := uuid.NewV4().String()
	unpublishedEditionID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given the dataset, edition and version do not exist", t, func() {
		Convey("When an authorised request to get the metadata relevant to a version", func() {
			Convey("Then return status bad request (400) with message `Dataset not found`", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Dataset not found\n")
			})
		})
	})

	Convey("Given an unpublished dataset exist", t, func() {
		if err := mongo.Setup(database, collection, "_id", unpublishedDatasetID, validAssociatedDatasetData(unpublishedDatasetID)); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("but an edition and version do not exist", func() {
			Convey("When a request to get the metadata relevant to a version", func() {
				Convey("Then return status bad request (400) with message `Edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", unpublishedDatasetID, edition).WithHeader(internalToken, internalTokenID).
						Expect().Status(http.StatusBadRequest).Body().Contains("Edition not found\n")
				})
			})
		})

		Convey("and an unpublished edition exist", func() {
			if err := mongo.Setup(database, "editions", "_id", unpublishedEditionID, validUnpublishedEditionData(unpublishedDatasetID, unpublishedEditionID, edition)); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("but a version does not exist", func() {
				Convey("When a request to get the metadata relevant to a version", func() {
					Convey("Then return status bad request (404) with message `Version not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", unpublishedDatasetID, edition).WithHeader(internalToken, internalTokenID).
							Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")
					})
				})
			})

			if err := mongo.Teardown(database, "editions", "_id", unpublishedEditionID); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		})

		if err := mongo.Teardown(database, collection, "_id", unpublishedDatasetID); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	// Similar tests for unauthorised requests
	Convey("Given an unpublished dataset", t, func() {
		if err := mongo.Setup(database, collection, "_id", unpublishedDatasetID, validAssociatedDatasetData(unpublishedDatasetID)); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthorised request to get the metadate relevant to a version", func() {
			Convey("Then return status bad request (400) with message `Dataset not found`", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", datasetID, edition).
					Expect().Status(http.StatusBadRequest).Body().Contains("Dataset not found\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", unpublishedDatasetID); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given a published dataset", t, func() {
		if err := mongo.Setup(database, collection, "_id", datasetID, validPublishedDatasetData(datasetID)); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("and an unpublished edition", func() {
			if err := mongo.Setup(database, "editions", "_id", unpublishedEditionID, validUnpublishedEditionData(datasetID, unpublishedEditionID, edition)); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthorised request to get the metadata relevant to a version", func() {
				Convey("Then return status bad request (400) with message `Edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", datasetID, edition).
						Expect().Status(http.StatusBadRequest).Body().Contains("Edition not found\n")
				})
			})

			if err := mongo.Teardown(database, "editions", "_id", unpublishedEditionID); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		})

		Convey("and a published edition but an unpublished version", func() {
			if err := mongo.Setup(database, "editions", "_id", editionID, validPublishedEditionData(datasetID, editionID, edition)); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			if err := mongo.Setup(database, "instances", "_id", unpublishedInstanceID, validAssociatedInstanceData(datasetID, editionID, unpublishedInstanceID)); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("When an unauthorised request to get the metadata relevant to a version", func() {
				Convey("Then return status not found (404) with message `Version not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/metadata", datasetID, edition).
						Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")
				})
			})

			if err := mongo.Teardown(database, "editions", "_id", editionID); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			if err := mongo.Teardown(database, "instances", "_id", unpublishedInstanceID); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}
		})

		if err := mongo.Teardown(database, collection, "_id", datasetID); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})
}

func setupMetadataDocs(datasetID, editionID, edition, instanceID, unpublishedInstanceID string) (*mongo.ManyDocs, error) {
	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	publishedEditionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	publishedVersionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	unpublishedVersionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID),
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, publishedVersionDoc, unpublishedVersionDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.SetupMany(d); err != nil {
		return nil, err
	}

	return d, nil
}
