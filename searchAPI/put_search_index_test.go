package searchAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyCreateSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)
	dimension := "aggregate"

	Convey("Given an elasticsearch index does not exists for an instance dimension", t, func() {
		Convey("When a PUT request is made to search API with valid authentication header", func() {
			Convey("Then a message is sent to kafka to create an index and the response returns status ok (200)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, "aggregate").
					WithHeader(internalTokenHeader, internalTokenID).Expect().Status(http.StatusOK)
			})
		})
	})

	// Eventhough kafka message is created, the index will never be created as the
	// instance id will not match any instance that exists against an environment/local
	// database, but just in the unlikely scenario this does create an index, delete it
	deleteIndex(instanceID, dimension)
}

func TestFailToCreateSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)
	dimension := "aggregate"

	Convey("Given an elasticsearch index does not exists for an instance dimension", t, func() {
		deleteIndex(instanceID, dimension)

		Convey("When a PUT request is made to search API without an authentication header", func() {
			Convey("Then the response returns status not found (404)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, "aggregate").
					Expect().Status(http.StatusNotFound).Body().Contains("Resource not found\n")
			})
		})

		Convey("When a PUT request is made to search API with Invalid authentication header", func() {
			Convey("Then the response returns status not found (404)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, "aggregate").
					WithHeader(internalTokenHeader, "grey").Expect().Status(http.StatusNotFound).Body().Contains("Resource not found\n")
			})
		})
	})

	deleteIndex(instanceID, dimension)
}

func deleteIndex(ID, dimension string) {
	path := cfg.ElasticSearchAPIURL + "/" + ID + "_" + dimension
	statusCode, err := elasticsearch.DeleteIndex(path)
	if err != nil {
		if statusCode != http.StatusNotFound {
			log.ErrorC("failed to delete index", err, log.Data{"path": path})
			os.Exit(1)
		}
	}
}
