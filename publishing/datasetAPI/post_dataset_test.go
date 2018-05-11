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
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyPostDataset(t *testing.T) {
	datasetID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
	}

	Convey("Given a dataset with the an id of ["+datasetID+"] does not exist", t, func() {

		Convey("When an authorised POST request is made to create a dataset resource", func() {
			Convey("Then return a status ok and the expected response body", func() {
				response := datasetAPI.POST("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTCreateDatasetJSON)).
					Expect().Status(http.StatusCreated).JSON().Object()

				response.Value("id").Equal(datasetID)
				response.Value("next").Object().Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
				response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
				response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
				response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
				response.Value("next").Object().Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
				response.Value("next").Object().Value("keywords").Array().Element(0).Equal("cpi")
				response.Value("next").Object().Value("id").Equal(datasetID)
				response.Value("next").Object().Value("license").Equal("ONS license")
				response.Value("next").Object().Value("links").Object().Value("access_rights").Object().Value("href").Equal("http://ons.gov.uk/accessrights")
				response.Value("next").Object().Value("links").Object().Value("editions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions$")
				response.Value("next").Object().Value("links").Object().Value("self").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
				response.Value("next").Object().Value("methodologies").Array().Element(0).Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.")
				response.Value("next").Object().Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("next").Object().Value("methodologies").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("next").Object().Value("national_statistic").Boolean().True()
				response.Value("next").Object().Value("next_release").Equal("17 October 2017")
				response.Value("next").Object().Value("publications").Array().Element(0).Object().Value("description").Equal("Price indices, percentage changes and weights for the different measures of consumer price inflation.")
				response.Value("next").Object().Value("publications").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
				response.Value("next").Object().Value("publications").Array().Element(0).Object().Value("title").Equal("UK consumer price inflation: August 2017")
				response.Value("next").Object().Value("publisher").Object().Value("name").Equal("Automation Tester")
				response.Value("next").Object().Value("publisher").Object().Value("type").Equal("publisher")
				response.Value("next").Object().Value("publisher").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
				response.Value("next").Object().Value("qmi").Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall")
				response.Value("next").Object().Value("qmi").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
				response.Value("next").Object().Value("qmi").Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
				response.Value("next").Object().Value("related_datasets").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices")
				response.Value("next").Object().Value("related_datasets").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation time series dataset")
				response.Value("next").Object().Value("release_frequency").Equal("Monthly")
				response.Value("next").Object().Value("state").Equal("created")
				response.Value("next").Object().Value("theme").Equal("Goods and services")
				response.Value("next").Object().Value("title").Equal("CPI")
				response.Value("next").Object().Value("unit_of_measure").Equal("Pounds Sterling")
				response.Value("next").Object().Value("uri").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation")

				if err := mongo.Teardown(dataset); err != nil {
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}
			})
		})
	})
}

func TestFailureToPostDataset(t *testing.T) {

	datasetID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given the dataset does not already exist", t, func() {
		Convey("When an authorised POST request is made to create dataset resource with an invalid body", func() {

			Convey("Then return a status of bad request with a message `Failed to parse json body`", func() {

				datasetAPI.POST("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, florenceToken).WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).Body().Contains("Failed to parse json body")
			})
		})

		Convey("When an unauthorised POST request is made to create a dataset resource with an invalid authentication header", func() {
			Convey("Then return a status unauthorized (401)", func() {

				datasetAPI.POST("/datasets/{id}", datasetID).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					WithBytes([]byte(validPOSTCreateDatasetJSON)).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When no authentication header is provided in POST request to create a dataset resource", func() {
			Convey("Then return a status not found (404) with a message `requested resource not found`", func() {

				datasetAPI.POST("/datasets/{id}", datasetID).
					WithBytes([]byte(validPOSTCreateDatasetJSON)).
					Expect().Status(http.StatusNotFound).Body().Contains("requested resource not found")
			})
		})
	})

	Convey("Given a dataset does exist", t, func() {
		publishedDataset := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: collection,
			Key:        "_id",
			Value:      datasetID,
			Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
		}

		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised POST request to create the same dataset resource is made", func() {
			Convey("Then return a status of forbidden with a message `forbidden - dataset already exists`", func() {

				datasetAPI.POST("/datasets/{id}", datasetID).
					WithBytes([]byte(validPOSTCreateDatasetJSON)).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusForbidden).Body().Contains("forbidden - dataset already exists")
			})
		})

		if err := mongo.Teardown(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	})
}
