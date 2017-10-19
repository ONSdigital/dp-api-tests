package datasetAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAListOfDatasets(t *testing.T) {

	mongo.Teardown(database, "instances", "_id", instanceID)
	mongo.Setup(database, "instances", "_id", instanceID, validPublishedInstanceData)

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of instances", t, func() {
		Convey("when the user is authorised", func() {
			response := datasetAPI.GET("/instances").WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Element(0).Object().Value("id").NotNull()
		})
	})

	Convey("Fail to get a list of instances", t, func() {
		Convey("When the user is unauthorised", func() {
			datasetAPI.GET("/instances").
				Expect().Status(http.StatusUnauthorized)
		})
	})

	mongo.Teardown(database, "instances", "_id", instanceID)
}
