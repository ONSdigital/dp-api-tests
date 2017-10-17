package datasetAPI

import (
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulGetAListOfDatasets(t *testing.T) {

	mongo.Teardown(database, collection, "_id", datasetID)
	mongo.Setup(database, collection, "_id", datasetID, validPublishedDatasetData)

	// TODO reinstate assertions
	//datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of datasets", t, func() {
		// Convey("when the user is unauthorised", func() {
		// 	response := datasetAPI.GET("/datasets").
		// 		Expect().Status(http.StatusOK).JSON().Object()
		//
		// 	response.Value("items").Array().Element(0).Object().Value("id").NotNull()
		// })

		// Convey("when the user is authorised", func() {
		// 	response := datasetAPI.GET("/datasets").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
		// 		Expect().Status(http.StatusOK).JSON().Object()
		//
		// 	response.Value("items").Array().Element(0).Object().Value("id").NotNull()
		// 	response.Value("items").Array().Element(0).Object().Value("current").NotNull()
		// 	response.Value("items").Array().Element(0).Object().Value("next").NotNull()
		// })
	})

	mongo.Teardown(database, collection, "_id", datasetID)
}
