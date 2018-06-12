package searchAPI

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/gedge/mgo"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/publishing/datasetAPI"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
)

const (
	timeout               = 5 * time.Second
	retryPause            = 750 * time.Millisecond
	dimensionKeyAggregate = "aggregate"
)

func TestSuccessfullyGetDimensionViaSearch(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()

	edition := "2017"

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     datasetAPI.ValidPublishedEditionData(datasetID, editionID, edition),
	}

	versionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	if err := mongo.Setup(datasetDoc, editionDoc, versionDoc); err != nil {
		log.ErrorC("was unable to run test", err, nil)
		os.Exit(1)
	}

	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)

	if err := createSearchIndex(cfg.ElasticSearchAPIURL, instanceID, dimensionKeyAggregate); err != nil {
		log.ErrorC("Unable to setup elasticsearch index with test data", err, nil)
		os.Exit(1)
	}

	Convey("Given an existing version for an edition of a dataset is published", t, func() {
		Convey("When a GET request is made with a query term matching a dimension code", func() {
			Convey("Then the response returns a json document containing a list of results with a status ok (200)", func() {

				exitSearchCompleteLoop := make(chan bool)
				go func() {
					time.Sleep(timeout)
					close(exitSearchCompleteLoop)
				}()

				foundCount := false
			searchCompleteLoop:
				for {
					select {
					case <-exitSearchCompleteLoop:
						break searchCompleteLoop
					default:
						response := searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
							WithQuery("q", "cpih1dim1S10201").
							WithHeader(common.AuthHeaderKey, serviceToken).
							Expect().Status(http.StatusOK).
							JSON().Object()

						if count, ok := response.Value("count").Raw().(float64); ok {
							if count != 0 {
								foundCount = true
								break searchCompleteLoop
							}
						}
						log.DebugC("searchCompleteLoop", "got empty search results", log.Data{"resp": response.Raw()})
						time.Sleep(retryPause) // Don't want to batter the api
					}
				}

				if !foundCount && false {
					err := errors.New("timed out")
					log.ErrorC("Timed out - failed to get list of search results", err, log.Data{"timeout": timeout})
					os.Exit(1)
				}

				response := searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithQuery("q", "cpih1dim1S10201").
					WithHeader(common.AuthHeaderKey, serviceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("count").Equal(1)
				response.Value("items").Array().Length().Equal(1)
				response.Value("items").Array().Element(0).Object().Value("code").Equal("cpih1dim1S10201")
				response.Value("items").Array().Element(0).Object().Value("dimension_option_url").Equal("http://localhost:22600/hierarchies/bb109aff-8b36-4ada-8279-70304927b2bb/aggregate/cpih1dim1S10201")
				response.Value("items").Array().Element(0).Object().Value("has_data").Equal(false)
				response.Value("items").Array().Element(0).Object().Value("label").Equal("01.2.1 Coffee, tea and cocoa")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().NotContainsKey("label")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("code").Array().Length().Equal(1)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("code").Array().Element(0).Object().Value("start").Equal(1)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("code").Array().Element(0).Object().Value("end").Equal(15)
				response.Value("items").Array().Element(0).Object().Value("number_of_children").Equal(0)
				response.Value("limit").Equal(20)
				response.Value("offset").Equal(0)
			})
		})

		Convey("When a GET request is made with a query term matching a dimension label", func() {
			Convey("Then the response returns a json document cotaining a list of results with a status ok (200)", func() {

				response := searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithQuery("q", "Overall Index").
					WithHeader(common.AuthHeaderKey, serviceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("count").Equal(1)
				response.Value("items").Array().Length().Equal(1)
				response.Value("items").Array().Element(0).Object().Value("code").Equal("cpih1dim1A0")
				response.Value("items").Array().Element(0).Object().Value("dimension_option_url").Equal("http://localhost:22400/code-list/cpih1dim1aggid/code/cpih1dim1A0")
				response.Value("items").Array().Element(0).Object().Value("has_data").Equal(false)
				response.Value("items").Array().Element(0).Object().Value("label").Equal("Overall Index")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().NotContainsKey("code")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Length().Equal(2)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("start").Equal(1)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("end").Equal(7)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("start").Equal(9)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("end").Equal(13)
				response.Value("items").Array().Element(0).Object().Value("number_of_children").Equal(12)
				response.Value("limit").Equal(20)
				response.Value("offset").Equal(0)
			})
		})

		Convey("When a GET request is made with a query term matching multiple dimensions by label", func() {
			Convey("Then the response returns a json document cotaining a list of results with a status ok (200)", func() {
				response := searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithQuery("q", "Furniture and furnishings").
					WithHeader(common.AuthHeaderKey, serviceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("count").Equal(7)
				response.Value("items").Array().Length().Equal(7)
				response.Value("items").Array().Element(0).Object().Value("code").Equal("cpih1dim1S50101")
				response.Value("items").Array().Element(0).Object().Value("dimension_option_url").Equal("http://localhost:22600/hierarchies/bb109aff-8b36-4ada-8279-70304927b2bb/aggregate/cpih1dim1S50101")
				response.Value("items").Array().Element(0).Object().Value("has_data").Equal(false)
				response.Value("items").Array().Element(0).Object().Value("label").Equal("05.1.1 Furniture and furnishings")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().NotContainsKey("code")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Length().Equal(3)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("start").Equal(8)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("end").Equal(16)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("start").Equal(18)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("end").Equal(20)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(2).Object().Value("start").Equal(22)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(2).Object().Value("end").Equal(32)
				response.Value("items").Array().Element(0).Object().Value("number_of_children").Equal(0)
				response.Value("items").Array().Element(1).Object().Value("code").Equal("cpih1dim1T50000")
				response.Value("items").Array().Element(1).Object().Value("dimension_option_url").Equal("http://localhost:22600/hierarchies/bb109aff-8b36-4ada-8279-70304927b2bb/aggregate/cpih1dim1T50000")
				response.Value("items").Array().Element(1).Object().Value("has_data").Equal(false)
				response.Value("items").Array().Element(1).Object().Value("label").Equal("05 Furniture, household equipment and maintenance")
				response.Value("items").Array().Element(1).Object().Value("matches").Object().NotContainsKey("code")
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Length().Equal(2)
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("start").Equal(4)
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("end").Equal(12)
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("start").Equal(35)
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("end").Equal(37)
				response.Value("items").Array().Element(1).Object().Value("number_of_children").Equal(6)
				response.Value("limit").Equal(20)
				response.Value("offset").Equal(0)
			})
		})

		Convey("When a GET request is made with a query term matching multiple dimensions by label but has an offset of 2 and a limit of 2", func() {
			Convey("Then the response returns a json document cotaining a list of results with a status ok (200)", func() {
				response := searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithQuery("q", "Furniture and furnishings").
					WithQuery("offset", 2).
					WithQuery("limit", 2).
					WithHeader(common.AuthHeaderKey, serviceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("count").Equal(2)
				response.Value("items").Array().Length().Equal(2)
				response.Value("items").Array().Element(0).Object().Value("code").Equal("cpih1dim1S10201")
				response.Value("items").Array().Element(0).Object().Value("dimension_option_url").Equal("http://localhost:22600/hierarchies/bb109aff-8b36-4ada-8279-70304927b2bb/aggregate/cpih1dim1S10201")
				response.Value("items").Array().Element(0).Object().Value("has_data").Equal(false)
				response.Value("items").Array().Element(0).Object().Value("label").Equal("01.2.1 Coffee, tea and cocoa")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().NotContainsKey("code")
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Length().Equal(1)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("start").Equal(20)
				response.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("end").Equal(22)
				response.Value("items").Array().Element(0).Object().Value("number_of_children").Equal(0)
				response.Value("items").Array().Element(1).Object().Value("code").Equal("cpih1dim1T10000")
				response.Value("items").Array().Element(1).Object().Value("dimension_option_url").Equal("http://localhost:22600/hierarchies/bb109aff-8b36-4ada-8279-70304927b2bb/aggregate/cpih1dim1T10000")
				response.Value("items").Array().Element(1).Object().Value("has_data").Equal(false)
				response.Value("items").Array().Element(1).Object().Value("label").Equal("01 Food and non-alcoholic beverages")
				response.Value("items").Array().Element(1).Object().Value("matches").Object().NotContainsKey("code")
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Length().Equal(1)
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("start").Equal(9)
				response.Value("items").Array().Element(1).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("end").Equal(11)
				response.Value("items").Array().Element(1).Object().Value("number_of_children").Equal(2)
				response.Value("limit").Equal(2)
				response.Value("offset").Equal(2)
			})
		})
	})

	if skipTeardown {
		return
	}

	// delete mongo test data
	if err := mongo.Teardown(datasetDoc, editionDoc, versionDoc); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("was unable to remove test data", err, nil)
			os.Exit(1)
		}
	}

	path := cfg.ElasticSearchAPIURL + "/" + instanceID + "_" + dimensionKeyAggregate
	// delete search index
	log.Debug("deleteIndex", log.Data{"path": path})
	status, err := elasticsearch.DeleteIndex(path)
	if err != nil {
		log.ErrorC("failed to remove elastic search index", err, log.Data{"status_code": status})
		os.Exit(1)
	}
}

func TestFailureToGetDimensionViaSearch(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()

	edition := "2017"

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     datasetAPI.ValidPublishedEditionData(datasetID, editionID, edition),
	}

	versionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)

	Convey("Given a version for an edition of a dataset is not published", t, func() {
		Convey("When a GET request is made to search API", func() {
			Convey("Then the response returns unauthorized (401)", func() {

				searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithQuery("q", "Overall Index").
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	Convey("Given a version for an edition of a dataset is published", t, func() {
		if err := mongo.Setup(datasetDoc, editionDoc, versionDoc); err != nil {
			log.ErrorC("was unable to run test", err, nil)
			os.Exit(1)
		}
		Convey("When a GET request is made to search API without the query parameter 'q'", func() {
			Convey("Then the response returns Bad request (400)", func() {

				searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithHeader(common.AuthHeaderKey, serviceToken).
					WithQuery("limit", 10).
					Expect().Status(http.StatusBadRequest).Body().Contains("search term empty\n")
			})
		})

		Convey("When a GET request is made to search API with query parameter 'q' and an offset of 1000", func() {
			Convey("Then the response returns Bad request (400)", func() {

				searchAPI.GET("/search/datasets/{datasetID}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetID, edition, "1", dimensionKeyAggregate).
					WithQuery("q", "Overall Index").
					WithQuery("offset", 1000).
					WithHeader(common.AuthHeaderKey, serviceToken).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("the maximum offset has been reached, the offset cannot be more than 1000\n")
			})
		})

		if skipTeardown {
			return
		}

		if err := mongo.Teardown(datasetDoc, editionDoc, versionDoc); err != nil {
			log.ErrorC("was unable to remove test data", err, nil)
			os.Exit(1)
		}
	})
}
