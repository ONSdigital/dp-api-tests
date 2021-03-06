package searchAPI

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/searchAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyCreateSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)
	dimension := dimensionKeyAggregate

	Convey("Given an elasticsearch index does not exists for an instance dimension", t, func() {
		Convey("When a PUT request is made to search API with valid authentication header", func() {
			Convey("Then a message is sent to kafka to create an index and the response returns status ok (200)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, serviceToken).Expect().Status(http.StatusOK)
			})
		})
	})

	// Eventhough a kafka message is created, the index will never be created as the
	// instance id will not match any instance that exists against an environment/local
	// database, but just in the unlikely scenario this does create an index, delete it
	deleteIndex(t, instanceID, dimension)
}

func TestFailToCreateSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)
	dimension := dimensionKeyAggregate

	Convey("Given an elasticsearch index does not exists for an instance dimension", t, func() {
		deleteIndex(t, instanceID, dimension)

		Convey("When a PUT request is made to search API without an authentication header", func() {
			Convey("Then the response returns status unauthorized (401)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					Expect().
					Status(http.StatusUnauthorized)
			})
		})

		Convey("When a PUT request is made to search API with Invalid authentication header", func() {
			Convey("Then the response returns status unauthorized (401)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, "grey").
					Expect().
					Status(http.StatusUnauthorized)
			})
		})

		deleteIndex(t, instanceID, dimension)
	})
}

func deleteIndex(t *testing.T, ID, dimension string) {
	path := cfg.ElasticSearchAPIURL + "/" + ID + "_" + dimension
	statusCode, err := elasticsearch.DeleteIndex(path)
	if err != nil {
		if statusCode != http.StatusNotFound {
			log.ErrorC("failed to delete index", err, log.Data{"path": path})
			t.FailNow()
		}
	}
}
