package searchAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

const resourceNotFound = "esource not found"

func TestSuccessfullyDeleteSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)

	Convey("Given an elasticsearch index exists for an instance", t, func() {
		if err := createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, dimensionKeyAggregate); err != nil {
			log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
			os.Exit(1)
		}

		Convey("When a DELETE request is made to search API with valid authentication header", func() {
			Convey("Then the index is removed and response returns status ok (200)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, internalTokenID).Expect().Status(http.StatusOK)
			})
		})
	})
}

func TestFailToDeleteSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)
	path := cfg.ElasticSearchAPIURL + "/" + instanceID + "_" + dimensionKeyAggregate

	Convey("Given an elasticsearch index does not exist for an instance", t, func() {
		statusCode, err := elasticsearch.DeleteIndex(path)
		if err != nil {
			if statusCode != http.StatusNotFound {
				log.ErrorC("failed to delete index", err, log.Data{"path": path})
				os.Exit(1)
			}
		}

		Convey("When a DELETE request is made to search API with valid authentication header", func() {
			Convey("Then the response returns status not found (404)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, internalTokenID).Expect().Status(http.StatusNotFound).Body().Contains(resourceNotFound)
			})
		})
	})

	Convey("Given an elasticsearch index exist for an instance", t, func() {
		if err := createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, dimensionKeyAggregate); err != nil {
			log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
			os.Exit(1)
		}
		Convey("When a DELETE request is made to search API without an authentication header", func() {
			Convey("Then the response returns status not found (404)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					Expect().Status(http.StatusNotFound).Body().Contains(resourceNotFound)
			})
		})

		Convey("When a DELETE request is made to search API with Invalid authentication header", func() {
			Convey("Then the response returns status not found (404)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, "grey").Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	statusCode, err := elasticsearch.DeleteIndex(path)
	if err != nil {
		if statusCode != http.StatusNotFound {
			log.ErrorC("failed to delete index", err, log.Data{"path": path})
			os.Exit(1)
		}
	}
}
