package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetObservationForVersion(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, neo4j.ObservationTestData)
	if err != nil {
		log.ErrorC("Unable to connect to a neo4j instance", err, nil)
		os.Exit(1)
	}

	if err = publishedGraphData.Setup(); err != nil {
		log.ErrorC("Unable to setup graph data", err, nil)
		os.Exit(1)
	}

	unpublishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, unpublishedInstanceID, neo4j.ObservationTestData)
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

		Convey("When an authenticated request is made to get an observation resource for an unpublished version", func() {
			Convey("Then the response body contains the expected observation data", func() {

				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimensions").Object().Value("aggregate").Object().Value("options").Array().Length().Equal(1)
				response.Value("dimensions").Object().Value("aggregate").Object().Value("options").Array().Element(0).Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD/codes/cpi1dim1S40403$")
				response.Value("dimensions").Object().Value("aggregate").Object().Value("options").Array().Element(0).Object().Value("id").Equal("cpi1dim1S40403")
				response.Value("dimensions").Object().Value("geography").Object().Value("options").Array().Length().Equal(1)
				response.Value("dimensions").Object().Value("geography").Object().Value("options").Array().Element(0).Object().Value("href").String().Match("/codelists/708064B3-A808-449B-9041-EA3A2F72CFAF/codes/K02000001$")
				response.Value("dimensions").Object().Value("geography").Object().Value("options").Array().Element(0).Object().Value("id").Equal("K02000001")
				response.Value("dimensions").Object().Value("time").Object().Value("options").Array().Length().Equal(1)
				response.Value("dimensions").Object().Value("time").Object().Value("options").Array().Element(0).Object().Value("href").String().Match("/codelists/608064B3-A808-449B-9041-EA3A2F72CFAE/codes/Aug-16$")
				response.Value("dimensions").Object().Value("time").Object().Value("options").Array().Element(0).Object().Value("id").Equal("Aug-16")
				response.Value("observation").Equal("154.6")
				response.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2/metadata$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetID + "/editions/" + edition + "/versions/2/observations\\?aggregate=cpi1dim1S40403&geography=K02000001&time=Aug-16$") // ?aggregate=cpi1dim1S40403&geography=K02000001&time=Aug-16")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("2")
				response.Value("unit_of_measure").Equal("Pounds Sterling")
			})
		})

		Convey("When an unauthenticated request is made to get an observation resource for a published version", func() {
			Convey("Then the response body contains the expected observation data", func() {
				response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/observations", datasetID, edition).
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1G50100").Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimensions").Object().Value("aggregate").Object().Value("options").Array().Length().Equal(1)
				response.Value("dimensions").Object().Value("aggregate").Object().Value("options").Array().Element(0).Object().Value("href").String().Match("/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD/codes/cpi1dim1G50100$")
				response.Value("dimensions").Object().Value("aggregate").Object().Value("options").Array().Element(0).Object().Value("id").Equal("cpi1dim1G50100")
				response.Value("dimensions").Object().Value("geography").Object().Value("options").Array().Length().Equal(1)
				response.Value("dimensions").Object().Value("geography").Object().Value("options").Array().Element(0).Object().Value("href").String().Match("/codelists/708064B3-A808-449B-9041-EA3A2F72CFAF/codes/K02000001$")
				response.Value("dimensions").Object().Value("geography").Object().Value("options").Array().Element(0).Object().Value("id").Equal("K02000001")
				response.Value("dimensions").Object().Value("time").Object().Value("options").Array().Length().Equal(1)
				response.Value("dimensions").Object().Value("time").Object().Value("options").Array().Element(0).Object().Value("href").String().Match("/codelists/608064B3-A808-449B-9041-EA3A2F72CFAE/codes/Aug-16$")
				response.Value("dimensions").Object().Value("time").Object().Value("options").Array().Element(0).Object().Value("id").Equal("Aug-16")
				response.Value("observation").Equal("117.9")
				response.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1/metadata$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetID + "/editions/" + edition + "/versions/1/observations\\?aggregate=cpi1dim1G50100&geography=K02000001&time=Aug-16$") // ?aggregate=cpi1dim1S40403&geography=K02000001&time=Aug-16")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
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

func TestFailureToGetObservationForVersion(t *testing.T) {

	instanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, neo4j.ObservationTestData)
	if err != nil {
		log.ErrorC("Unable to connect to a neo4j instance", err, nil)
		os.Exit(1)
	}

	if err = publishedGraphData.Setup(); err != nil {
		log.ErrorC("Unable to setup graph data", err, nil)
		os.Exit(1)
	}

	unpublishedGraphData, err := neo4j.NewDatastore(cfg.Neo4jAddr, unpublishedInstanceID, neo4j.ObservationTestData)
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
					WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("Dataset not found")
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
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("but edition and version do not exist", func() {
			Convey("When a request to get an observation for a version of a dataset", func() {
				Convey("Then return status not found (404) with message `Edition not found`", func() {
					datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
						WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").WithHeader(florenceTokenName, florenceToken).
						Expect().Status(http.StatusNotFound).Body().Contains("Edition not found")
				})
			})
		})

		Convey("and a published edition exist", func() {

			if err := mongo.Setup(publishedEditionDoc); err != nil {
				log.ErrorC("Was unable to run test", err, nil)
				os.Exit(1)
			}

			Convey("but a version does not exist", func() {
				Convey("When a request to get an observation for a version of a dataset", func() {
					Convey("Then return status not found (404) with message `Version not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").WithHeader(florenceTokenName, florenceToken).
							Expect().Status(http.StatusNotFound).Body().Contains("Version not found")
					})
				})
			})

			Convey("and an unpublished version exist", func() {

				if err := mongo.Setup(unpublishedVersionDoc); err != nil {
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}

				Convey("When a request to get an observation for a version of a dataset", func() {
					Convey("Then return status not found (404) with message `Version not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).
							WithQueryString("time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusNotFound).Body().Contains("Version not found")
					})
				})

				Convey("When a request to get an observation for a version of a dataset with incorrect query parameters", func() {
					Convey("Then return status bad request (400) with a message", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).WithHeader(florenceTokenName, florenceToken).
							WithQueryString("age=24&time=Aug-16&geography=K02000001&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusBadRequest).Body().Contains("Incorrect selection of query parameters: [age], these dimensions do not exist for this version of the dataset")
					})
				})

				Convey("When a request to get an observation for a version of a dataset with missing query parameters", func() {
					Convey("Then return status bad request (400) with a message", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).WithHeader(florenceTokenName, florenceToken).
							WithQueryString("geography=K02000001&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusBadRequest).Body().Contains("Missing query parameters for the following dimensions: [time]")
					})
				})

				Convey("When a request to get an observation for a version of a dataset with the correct query parameters but the values don't exist", func() {
					Convey("Then return status not found (404) with message `Observation not found`", func() {
						datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/2/observations", datasetID, edition).WithHeader(florenceTokenName, florenceToken).
							WithQueryString("time=Aug-17&geography=K02000001&aggregate=cpi1dim1S40403").
							Expect().Status(http.StatusNotFound).Body().Contains("Observation not found")
					})
				})
			})
		})

		if err := mongo.Teardown(datasetDoc, publishedEditionDoc, unpublishedVersionDoc); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
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
