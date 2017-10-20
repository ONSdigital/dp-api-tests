package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Create a dataset provided by the ONS that can be filtered using the filter API
// 200 - A json object containing a dataset which has been created
func TestPostCreateDataset_CreatesDataset(t *testing.T) {
	mongo.Teardown(database, collection, "_id", datasetID)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset", t, func() {

		response := datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()

		response.Value("id").Equal(datasetID)

		response.Value("next").Object().Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")

		response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
		response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
		response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")

		response.Value("next").Object().Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")

		response.Value("next").Object().Value("keywords").Array().Element(0).Equal("cpi")

		response.Value("next").Object().Value("id").Equal(datasetID)
		response.Value("next").Object().Value("links").Object().Value("editions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions$")

		response.Value("next").Object().Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")

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
		response.Value("next").Object().Value("uri").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation")
		mongo.Teardown(database, collection, "_id", datasetID)
	})
}

// 401 - Unauthorised to create/overwrite dataset
func TestPostCreateDataset_InvalidToken(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset with invalid token value", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F4651").WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusUnauthorized)

	})
}

// 401 - Unauthorised to create/overwrite dataset
func TestPostCreateDataset_WithoutToken(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset with out token", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusUnauthorized)

	})
}

// 400 - Invalid request body
func TestPostCreateDataset_InvalidBody(t *testing.T) {

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Create a dataset with an invalid json body", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte("{")).
			Expect().Status(http.StatusBadRequest)

	})
}

// TODO
// 403 - Forbidden to overwrite dataset, already published
