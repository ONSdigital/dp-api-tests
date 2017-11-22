package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetADataset(t *testing.T) {

	setupDataset(datasetID, validPublishedDatasetData)
	defer removeDataset(datasetID)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a dataset", t, func() {
		Convey("When the user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}", datasetID).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(datasetID)
			checkDatasetDoc(response.Value("current").Object())

			response.Value("next").NotNull()
			response.Value("next").Object().Value("state").Equal("created")
		})

		Convey("When the user is unauthenticated", func() {

			response := datasetAPI.GET("/datasets/{id}", datasetID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal(datasetID)
			checkDatasetDoc(response)
		})
	})
}

func TestFailureToGetADataset(t *testing.T) {

	removeExistingDataset(datasetID)
	defer removeDataset(datasetID)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Fail to get a dataset document", t, func() {
		Convey("and return status not found", func() {
			Convey("When the dataset document does not exist", func() {
				datasetAPI.GET("/datasets/{id}", datasetID).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound)
			})
			Convey("When the user is not authenticated and the dataset document is not published", func() {
				mongo.Setup(database, collection, "_id", "133", validUnpublishedDatasetData)
				datasetAPI.GET("/datasets/{id}", datasetID).
					Expect().Status(http.StatusNotFound)
			})
		})
	})
}

func checkDatasetDoc(response *httpexpect.Object) {
	response.Value("access_right").Equal("http://ons.gov.uk/accessrights")
	response.Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
	response.Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
	response.Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
	response.Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
	response.Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
	response.Value("keywords").Array().Element(0).Equal("cpi")
	response.Value("license").Equal("ONS license")
	response.Value("links").Object().Value("editions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions$")
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
	response.Value("methodologies").Array().Element(0).Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.")
	response.Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
	response.Value("methodologies").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
	response.Value("national_statistic").Boolean().True()
	response.Value("next_release").Equal("2017-10-10")
	response.Value("publications").Array().Element(0).Object().Value("description").Equal("Price indices, percentage changes and weights for the different measures of consumer price inflation.")
	response.Value("publications").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
	response.Value("publications").Array().Element(0).Object().Value("title").Equal("UK consumer price inflation: August 2017")
	response.Value("publisher").Object().Value("name").Equal("Automation Tester")
	response.Value("publisher").Object().Value("type").Equal("publisher")
	response.Value("publisher").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017")
	response.Value("qmi").Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall")
	response.Value("qmi").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
	response.Value("qmi").Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
	response.Value("related_datasets").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices")
	response.Value("related_datasets").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation time series dataset")
	response.Value("release_frequency").Equal("Monthly")
	response.Value("state").Equal("published")
	response.Value("theme").Equal("Goods and services")
	response.Value("title").Equal("CPI")
	response.Value("uri").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation")
}
