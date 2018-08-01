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

const (
	unauthorizedReq = "unauthorized request"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/searchAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyDeleteSearchIndex(t *testing.T) {
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)

	Convey("Given an elasticsearch index exists for an instance", t, func() {
		if err := createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, dimensionKeyAggregate); err != nil {
			log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
			t.FailNow()
		}

		Convey("When a DELETE request is made to search API with valid authentication header", func() {
			Convey("Then the index is removed and response returns status ok (200)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, serviceToken).Expect().Status(http.StatusOK)
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
				t.FailNow()
			}
		}

		Convey("When a DELETE request is made to search API with valid authentication header", func() {
			Convey("Then the response returns status not found (404)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, serviceToken).Expect().Status(http.StatusNotFound).Body().Contains("search index not found")
			})
		})
	})

	Convey("Given an elasticsearch index exist for an instance", t, func() {
		if err := createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, dimensionKeyAggregate); err != nil {
			log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
			t.FailNow()
		}
		Convey("When a DELETE request is made to search API without an authentication header", func() {
			Convey("Then the response returns status unauthorized (401)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a DELETE request is made to search API with Invalid authentication header", func() {
			Convey("Then the response returns status unauthorized (401)", func() {

				searchAPI.DELETE("/search/instances/{instanceID}/dimensions/{dimension}", instanceID, dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, "grey").
					Expect().
					Status(http.StatusUnauthorized)
			})
		})
	})

	statusCode, err := elasticsearch.DeleteIndex(path)
	if err != nil {
		if statusCode != http.StatusNotFound {
			log.ErrorC("failed to delete index", err, log.Data{"path": path})
			t.FailNow()
		}
	}
}
