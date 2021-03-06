package searchAPI

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

// this is to cover "Resource" and "resource" not found
const (
	resourceNotFound = "esource not found"
	unauthorizedReq  = "unauthorized request"
)

func TestPublishingEndpointsAreHiddenForWeb(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)
	dimension := dimensionKeyAggregate

	Convey("Given an elasticsearch index exists for an instance", t, func() {
		if err := createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, dimensionKeyAggregate); err != nil {
			log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
			os.Exit(1)
		}

		Convey("When a DELETE request is made to search API with valid authentication header", func() {
			Convey("Then response returns status not found (404)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(florenceTokenName, florenceToken).Expect().Status(http.StatusNotFound)
			})
		})

		deleteIndex(t, instanceID, dimension, true)
	})

	Convey("Given an elasticsearch index does not exists for an instance dimension", t, func() {
		Convey("When a PUT request is made to search API with valid authentication header", func() {
			Convey("Then a message is sent to kafka to create an index and the response returns status ok (200)", func() {

				searchAPI.PUT("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(florenceTokenName, florenceToken).Expect().Status(http.StatusNotFound)
			})
		})
	})

	// Even though kafka message is created, the index will never be created as the
	// instance id will not match any instance that exists against an environment/local
	// database, but just in the unlikely scenario this does create an index, delete it
	deleteIndex(t, instanceID, dimension, false)
}

func deleteIndex(t *testing.T, ID, dimension string, expectedIndex bool) {
	// Due to deletion may be attempted before the kafka message is acted on,
	// delay the request to delete index
	time.Sleep(2 * time.Second)

	path := cfg.ElasticSearchAPIURL + "/" + ID + "_" + dimension
	statusCode, err := elasticsearch.DeleteIndex(path)
	if err != nil {
		if statusCode != http.StatusNotFound {
			log.ErrorC("failed to delete index", err, log.Data{"path": path})
			t.FailNow()
		}
	} else {
		if !expectedIndex {
			err = errors.New("unexpected index")
			log.ErrorC("deleted an index that should not have existed", err, log.Data{"path": path})
			t.FailNow()
		}
	}
}
