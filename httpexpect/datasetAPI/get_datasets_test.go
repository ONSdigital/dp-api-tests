package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// This test may be slow due to iterating over results in dataset
// (which could be many)
func TestSuccessfulGetAListOfDatasets(t *testing.T) {

	d, err := setupTestDataForGetAListOfDatasets()
	if err != nil {
		log.ErrorC("Failed setting up test data", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of datasets", t, func() {
		Convey("when the user is unauthorised", func() {
			response := datasetAPI.GET("/datasets").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Element(0).Object().Value("id").NotNull()

			for i := 0; i < len(response.Value("items").Array().Iter()); i++ {
				//Unauthorised user so should NOT have an unpublished dataset in response
				response.Value("items").Array().Element(i).Object().Value("id").NotEqual("133")

				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == datasetID {
					// check the published test dataset document has the expected returned fields and values
					response.Value("items").Array().Element(i).Object().Value("id").Equal(datasetID)
					checkDatasetResponse(response.Value("items").Array().Element(i).Object())
				}
			}
		})

		Convey("when the user is authorised", func() {
			response := datasetAPI.GET("/datasets").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Element(0).Object().Value("id").NotNull()

			for i := 0; i < len(response.Value("items").Array().Iter()); i++ {

				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == datasetID {
					// check the published test dataset document has the expected returned fields and values
					checkDatasetResponse(response.Value("items").Array().Element(i).Object().Value("current").Object())
					response.Value("items").Array().Element(i).Object().Value("next").Object().NotEmpty()
				}

				if response.Value("items").Array().Element(i).Object().Value("id").String().Raw() == "133" {
					// check the published test dataset document has the expected returned fields and values
					response.Value("items").Array().Element(i).Object().NotContainsKey("current")
					response.Value("items").Array().Element(i).Object().Value("next").Object().NotEmpty()
				}
			}
		})
	})

	mongo.TeardownMany(d)
}

func checkDatasetResponse(response *httpexpect.Object) {
	response.Value("collection_id").Equal("108064B3-A808-449B-9041-EA3A2F72CFAA")
	response.Value("contacts").Array().Element(0).Object().Value("email").Equal("cpi@onstest.gov.uk")
	response.Value("contacts").Array().Element(0).Object().Value("name").Equal("Automation Tester")
	response.Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1633 123456")
	response.Value("description").Equal("Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.")
	response.Value("keywords").Array().Element(0).String().Equal("cpi")
	response.Value("keywords").Array().Element(1).String().Equal("boy")
	response.Value("links").Object().Value("editions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions$")
	response.Value("links").Object().Value("latest_version").Object().Value("id").Equal("1")
	response.Value("links").Object().Value("latest_version").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/2017/versions/1$")
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
	response.Value("methodologies").Array().Element(0).Object().Value("description").Equal("Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.")
	response.Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi")
	response.Value("methodologies").Array().Element(0).Object().Value("title").Equal("Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)")
	response.Value("national_statistic").Equal(true)
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

func setupTestDataForGetAListOfDatasets() (*mongo.ManyDocs, error) {
	var docs []mongo.Doc

	publishedDatasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData,
	}

	unpublishedDatasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      "133",
		Update:     validUnpublishedDatasetData,
	}

	docs = append(docs, publishedDatasetDoc, unpublishedDatasetDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		return nil, err
	}

	if err := mongo.SetupMany(d); err != nil {
		return nil, err
	}

	return d, nil
}
