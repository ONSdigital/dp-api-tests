package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gedge/mgo"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

const (
	observationTestData          = "../../testDataSetup/neo4j/instance.cypher"
	expectedNumberOfObservations = 137
)

func TestSuccessfullyGetObservationsForVersion(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, observationTestData)
	if err != nil {
		log.ErrorC("Unable to connect to a neo4j instance", err, nil)
		os.Exit(1)
	}

	if err = publishedGraphData.Setup(); err != nil {
		log.ErrorC("Unable to setup graph data", err, nil)
		os.Exit(1)
	}

	unpublishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, unpublishedInstanceID, observationTestData)
	if err != nil {
		log.ErrorC("Unable to connect to a neo4j instance", err, nil)
		os.Exit(1)
	}

	if err := unpublishedGraphData.Setup(); err != nil {
		log.ErrorC("Unable to setup graph data", err, nil)
		os.Exit(1)
	}

	Convey("Given a published and unpublished version", t, func() {
		docs, err := setupObservationDocs(datasetID, editionID, edition, instanceID, unpublishedInstanceID)
		if err != nil {
			log.ErrorC("Failed to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated request is made to get an observation resource for a published version", func() {
			Convey("Then the response body contains the expected observation data", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/observations", datasetID, edition).
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1G50100").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD/codes/cpi1dim1G50100$")
				response.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("id").Equal("cpi1dim1G50100")
				response.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("href").String().Match("/codelists/708064B3-A808-449B-9041-EA3A2F72CFAF/codes/K02000001$")
				response.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("id").Equal("K02000001")
				response.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("href").String().Match("/codelists/608064B3-A808-449B-9041-EA3A2F72CFAE/codes/Aug-16$")
				response.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("id").Equal("Aug-16")
				response.Value("limit").Equal(10000)
				response.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1/metadata$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetID + "/editions/" + edition + "/versions/1/observations\\?aggregate=cpi1dim1G50100&geography=K02000001&time=Aug-16$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("observations").Array().Length().Equal(1)
				response.Value("observations").Array().Element(0).Object().Value("observation").Equal("117.9")
				response.Value("offset").Equal(0)
				response.Value("total_observations").Equal(1)
				response.Value("unit_of_measure").Equal("Pounds Sterling")
			})
		})

		Convey("When an authenticated request is made to get an observation resource for an unpublished version", func() {
			Convey("Then the response body contains the expected observation data", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD/codes/cpi1dim1S40403$")
				response.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("id").Equal("cpi1dim1S40403")
				response.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("href").String().Match("/codelists/708064B3-A808-449B-9041-EA3A2F72CFAF/codes/K02000001$")
				response.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("id").Equal("K02000001")
				response.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("href").String().Match("/codelists/608064B3-A808-449B-9041-EA3A2F72CFAE/codes/Aug-16$")
				response.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("id").Equal("Aug-16")
				response.Value("limit").Equal(10000)
				response.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2/metadata$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetID + "/editions/" + edition + "/versions/2/observations\\?aggregate=cpi1dim1S40403&geography=K02000001&time=Aug-16$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("2")
				response.Value("observations").Array().Length().Equal(1)
				response.Value("observations").Array().Element(0).Object().Value("observation").Equal("154.6")
				response.Value("offset").Equal(0)
				response.Value("total_observations").Equal(1)
				response.Value("unit_of_measure").Equal("Pounds Sterling")
			})
		})

		Convey("When a request is made to get an observations resource containing more than one observation for an unpublished version", func() {
			Convey("Then the response body contains the expected observations data", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=*").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("href").String().Match("/codelists/708064B3-A808-449B-9041-EA3A2F72CFAF/codes/K02000001$")
				response.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("id").Equal("K02000001")
				response.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("href").String().Match("/codelists/608064B3-A808-449B-9041-EA3A2F72CFAE/codes/Aug-16$")
				response.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("id").Equal("Aug-16")
				response.Value("limit").Equal(10000)
				response.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2/metadata$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetID + "/editions/" + edition + "/versions/2/observations\\?aggregate=\\%2A&geography=K02000001&time=Aug-16$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("2")
				response.Value("observations").Array().Length().Equal(expectedNumberOfObservations)

				// check two observations in observations array
				var firstObservation, secondObservation bool
				count := make(map[string]int)
				for _, observation := range response.Value("observations").Array().Iter() {

					if observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("id").String().Raw() == "cpi1dim1S10107" {
						observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD/codes/cpi1dim1S10107")
						observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("id").Equal("cpi1dim1S10107")
						observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("label").Equal("01.1.7 Vegetables including potatoes and tubers")
						observation.Object().Value("dimensions").Object().NotContainsKey("geography")
						observation.Object().Value("observation").Equal("136.8")
						firstObservation = true
					}

					if observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("id").String().Raw() == "cpi1dim1G100000" {
						observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD/codes/cpi1dim1G100000")
						observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("id").Equal("cpi1dim1G100000")
						observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("label").Equal("10.0 Education")
						observation.Object().Value("dimensions").Object().NotContainsKey("geography")
						observation.Object().Value("observation").Equal("244.3")
						secondObservation = true
					}

					count[observation.Object().Value("dimensions").Object().Value("Aggregate").Object().Value("id").String().Raw()] = 1
				}

				if !firstObservation || !secondObservation {
					t.Errorf("failed to find observations, \nfirst observation: [%v]\nsecond observation: [%v]\n", firstObservation, secondObservation)
					t.Fail()
				}

				if len(count) != expectedNumberOfObservations {
					t.Errorf("failed to find [%d] unique observations via unique codes, found: [%d]", expectedNumberOfObservations, len(count))
					t.Fail()
				}

				response.Value("offset").Equal(0)
				response.Value("total_observations").Equal(137)
				response.Value("unit_of_measure").Equal("Pounds Sterling")
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})

	publishedGraphData.TeardownInstance()
	unpublishedGraphData.TeardownInstance()
}

func TestFailureToGetObservationsForVersion(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, observationTestData)
	if err != nil {
		log.ErrorC("Unable to connect to a neo4j instance", err, nil)
		os.Exit(1)
	}

	if err = publishedGraphData.Setup(); err != nil {
		log.ErrorC("Unable to setup graph data", err, nil)
		os.Exit(1)
	}

	unpublishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, unpublishedInstanceID, observationTestData)
	if err != nil {
		log.ErrorC("Unable to connect to a neo4j instance", err, nil)
		os.Exit(1)
	}

	if err := unpublishedGraphData.Setup(); err != nil {
		log.ErrorC("Unable to setup graph data", err, nil)
		os.Exit(1)
	}

	Convey("Given the dataset, edition and version do not exist", t, func() {
		Convey("When an authorised request to get an observation for a version of a dataset", func() {
			Convey("Then return status not found (404) with message `Dataset not found`", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Dataset not found")
			})
		})
	})

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	unpublishedVersionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID),
	}

	Convey("Given a published dataset exist", t, func() {
		if err := mongo.Setup(datasetDoc); err != nil {
			log.ErrorC("Unable to set up published dataset doc", err, nil)
			os.Exit(1)
		}

		Convey("but edition and version do not exist", func() {
			Convey("When a request to get an observation for a version of a dataset", func() {
				Convey("Then return status not found (404) with message `Edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
						WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
						WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).
						Body().Contains("Edition not found")
				})
			})
		})

		Convey("and a published edition exist", func() {

			if err := mongo.Setup(publishedEditionDoc); err != nil {
				log.ErrorC("Unable to set up published edition doc", err, nil)
				os.Exit(1)
			}

			Convey("but a version does not exist", func() {
				Convey("When a request to get an observation for a version of a dataset", func() {
					Convey("Then return status not found (404) with message `Version not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
							WithHeader(florenceTokenName, florenceToken).
							Expect().Status(http.StatusNotFound).
							Body().Contains("Version not found")
					})
				})
			})

			Convey("and an unpublished version exist", func() {

				if err := mongo.Setup(unpublishedVersionDoc); err != nil {
					log.ErrorC("Unable to set up unpublished version doc", err, nil)
					os.Exit(1)
				}

				Convey("When a request to get an observation for unpublished version of a dataset with incorrect query parameters", func() {
					Convey("Then return status bad request (400) with a message", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithHeader(florenceTokenName, florenceToken).
							WithQueryString("age=24&gender=male&time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusBadRequest).
							Body().Match(`Incorrect selection of query parameters: \[(age gender|gender age)\], these dimensions do not exist for this version of the dataset`)
					})
				})

				Convey("When a request to get an observation for unpublished version of a dataset with missing query parameters", func() {
					Convey("Then return status bad request (400) with a message", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithHeader(florenceTokenName, florenceToken).
							WithQueryString("geography=K02000001").
							Expect().Status(http.StatusBadRequest).
							Body().Match(`Missing query parameters for the following dimensions: \[(time aggregate|aggregate time)\]`)
					})
				})

				Convey("When a request to get an observation for a version of a dataset with more than one wildcard used", func() {
					Convey("Then return status bad request (400) with a message", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithHeader(florenceTokenName, florenceToken).
							WithQueryString("time=*&geography=*&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusBadRequest).
							Body().Contains("only one wildcard (*) is allowed as a value in selected query parameters")
					})
				})

				Convey("When a request to get an observation for a version of a dataset with more than one value per query parameter", func() {
					Convey("Then return status bad request (400) with a message", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithHeader(florenceTokenName, florenceToken).
							WithQueryString("time=Aug-16&time=Aug-17&geography=K02000001&geography=*&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusBadRequest).
							Body().Match(`Multi-valued query parameters for the following dimensions: \[(time geography|geography time)\]`)
					})
				})

				Convey("When a request to get an observation for an unpublished version of a dataset with the correct query parameters but the values don't exist", func() {
					Convey("Then return status not found (404) with message `No observations found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithHeader(florenceTokenName, florenceToken).
							WithQueryString("time=Aug-17&geography=K02000001&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusNotFound).
							Body().Contains("No observations found")
					})
				})
			})
		})

		if err := mongo.Teardown(datasetDoc, publishedEditionDoc, unpublishedVersionDoc); err != nil {
			log.ErrorC("Unable to teardown docs: datasetDoc, publishedEditionDoc, unpublishedVersionDoc", err, nil)
			os.Exit(1)
		}
	})

	publishedGraphData.TeardownInstance()
	unpublishedGraphData.TeardownInstance()
}

func setupObservationDocs(datasetID, editionID, edition, instanceID, unpublishedInstanceID string) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	publishedEditionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	publishedVersionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	unpublishedVersionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      unpublishedInstanceID,
		Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID),
	}

	docs = append(docs, datasetDoc, publishedEditionDoc, publishedVersionDoc, unpublishedVersionDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
