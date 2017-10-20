package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// Update the meta data for the next release of the dataset.
// The dataset contains all high level information, for additional details see editions or versions of a dataset.
// 200 - A json object for a single Dataset
func TestPutUpdateDataset_UpdatesDataset(t *testing.T) {
	mongo.Teardown(database, collection, "_id", datasetID)

	mongo.Setup(database, "datasets", "_id", datasetID, validPublishedDatasetData)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Update the dataset", t, func() {

		datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte(validPUTUpdateDatasetJSON)).
			Expect().Status(http.StatusOK)
	})

	Convey("Get the updated dataset details", t, func() {

		response := datasetAPI.GET("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
			Expect().Status(http.StatusOK).JSON().Object()

		response.Value("id").Equal(datasetID)

		response.Value("next").Object().Value("collection_id").Equal("308064B3-A808-449B-9041-EA3A2F72CFAC")

		response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("email").Equal("rpi@onstest.gov.uk")
		response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("name").Equal("Test Automation")
		response.Value("next").Object().Value("contacts").Array().Element(0).Object().Value("telephone").Equal("+44 (0)1833 456123")

		response.Value("next").Object().Value("description").Equal("Producer Price Indices (PPIs) are a series of economic indicators that measure the price movement of goods bought and sold by UK manufacturers")

		response.Value("next").Object().Value("keywords").Array().Element(0).Equal("rpi")

		response.Value("next").Object().Value("methodologies").Array().Element(0).Object().Value("description").Equal("The Producer Price Index (PPI) is a monthly survey that measures the price changes of goods bought and sold by UK manufacturers")
		response.Value("next").Object().Value("methodologies").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/producerpriceindicesqmi")
		response.Value("next").Object().Value("methodologies").Array().Element(0).Object().Value("title").Equal("Producer price indices QMI")

		// response.Value("next").Object().Value("national_statistic").Boolean().False()

		response.Value("next").Object().Value("next_release").Equal("18 September 2017")

		response.Value("next").Object().Value("publications").Array().Element(0).Object().Value("description").Equal("Changes in the prices of goods bought and sold by UK manufacturers including price indices of materials and fuels purchased (input prices) and factory gate prices (output prices)")
		response.Value("next").Object().Value("publications").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/producerpriceinflation/september2017")
		response.Value("next").Object().Value("publications").Array().Element(0).Object().Value("title").Equal("Producer price inflation, UK: September 2017")

		response.Value("next").Object().Value("publisher").Object().Value("name").Equal("Test Automation Engineer")
		response.Value("next").Object().Value("publisher").Object().Value("type").Equal("publisher")
		response.Value("next").Object().Value("publisher").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/producerpriceinflation/september2017")

		response.Value("next").Object().Value("qmi").Object().Value("description").Equal("PPI provides an important measure of inflation")
		response.Value("next").Object().Value("qmi").Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/producerpriceindicesqmi")
		response.Value("next").Object().Value("qmi").Object().Value("title").Equal("The Producer Price Index (PPI) is a monthly survey that measures the price changes")

		response.Value("next").Object().Value("related_datasets").Array().Element(0).Object().Value("href").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex")
		response.Value("next").Object().Value("related_datasets").Array().Element(0).Object().Value("title").Equal("Producer Price Index time series dataset")

		response.Value("next").Object().Value("release_frequency").Equal("Quaterly")
		// response.Value("next").Object().Value("state").Equal("created")
		response.Value("next").Object().Value("theme").Equal("Price movement of goods")
		response.Value("next").Object().Value("title").Equal("RPI")
		response.Value("next").Object().Value("uri").Equal("https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex")
		mongo.Teardown(database, "datasets", "_id", datasetID)

	})
}

// 401 - Unauthorised to create/overwrite dataset
func TestPUTUpdateDataset_InvalidToken(t *testing.T) {

	mongo.Teardown(database, collection, "_id", datasetID)

	mongo.Setup(database, "datasets", "_id", datasetID, validPublishedDatasetData)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Update a dataset with invalid token value", t, func() {

		datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F4651").WithBytes([]byte(validPUTUpdateDatasetJSON)).
			Expect().Status(http.StatusUnauthorized)

	})

	mongo.Teardown(database, "datasets", "_id", datasetID)
}

// 401 - Unauthorised to create/overwrite dataset
func TestPUTUpdateDataset_WithoutToken(t *testing.T) {

	mongo.Teardown(database, collection, "_id", datasetID)
	mongo.Setup(database, "datasets", "_id", datasetID, validPublishedDatasetData)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Update a dataset with out token", t, func() {

		datasetAPI.POST("/datasets/{id}", datasetID).WithBytes([]byte(validPUTUpdateDatasetJSON)).
			Expect().Status(http.StatusUnauthorized)

	})
	mongo.Teardown(database, "datasets", "_id", datasetID)
}

// 400 - Invalid request body
func TestPUTUpdateDataset_InvalidBody(t *testing.T) {

	mongo.Teardown(database, collection, "_id", datasetID)
	mongo.Setup(database, "datasets", "_id", datasetID, validPublishedDatasetData)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Update a dataset with an invalid json body", t, func() {

		datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").WithBytes([]byte("{")).
			Expect().Status(http.StatusBadRequest)

	})
	mongo.Teardown(database, "datasets", "_id", datasetID)
}
