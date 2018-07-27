package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/helpers"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

func TestGetInstanceDimensions_ReturnsAllDimensionsFromAnInstance(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      ids.DatasetPublished,
		Update:     ValidPublishedWithUpdatesDatasetData(ids.DatasetPublished),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      ids.EditionPublished,
		Update:     ValidPublishedEditionData(ids.DatasetPublished, ids.EditionPublished, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      ids.InstancePublished,
		Update:     validPublishedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
	}

	dimensionOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData("9811", ids.InstancePublished),
	}

	dimensionTwoDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData("9812", ids.InstancePublished),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a list of dimensions for an instance exists", t, func() {
		docs, err := getInstanceDimensionsSetup(ids.DatasetPublished, ids.EditionPublished, edition, ids.InstancePublished, ids.UniqueTimestamp)
		if err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated user sends a GET request for a list of dimensions for instance", func() {
			Convey("Then a list of dimensions is returned with a status of OK (200)", func() {

				response := datasetAPI.GET("/instances/{instance_id}/dimensions", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("items").Array().Length().Equal(2)

				checkInstanceDimensionsResponse(response)
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Failed to tear down test data", err, nil)
				os.Exit(1)
			}
		}
	})

	Convey("Given no dimensions exist for an existing instance", t, func() {
		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstancePublished,
			Update:     validEditionConfirmedInstanceData(ids.DatasetPublished, edition, ids.InstancePublished, ids.UniqueTimestamp),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated user sends a GET request for a list of dimensions for instance", func() {
			Convey("Then return status OK (200) with an empty items array", func() {

				dimensionsResource := datasetAPI.GET("/instances/{id}/dimensions", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				dimensionsResource.Value("items").Null()
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Failed to tear down test data", err, nil)
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetInstanceDimensions(t *testing.T) {
	ids, err := helpers.GetIDsAndTimestamps()
	if err != nil {
		log.ErrorC("unable to generate mongo timestamp", err, nil)
		t.FailNow()
	}

	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance document does not exist", t, func() {
		Convey("When a user sends a GET request of a list of dimensions for instance without sending a token", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions", ids.InstancePublished).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a user sends a GET request of a list of dimensions for instance with an invalid token", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions", ids.InstancePublished).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an authenticated user sends a GET request of a list of dimensions for instance", func() {
			Convey("Then return status not found (404) with a message `instance not found`", func() {

				datasetAPI.GET("/instances/{id}/dimensions", ids.InstancePublished).
					WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("instance not found\n")
			})
		})
	})

	Convey("Given an instance document does exist", t, func() {

		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      ids.InstanceEditionConfirmed,
			Update:     validEditionConfirmedInstanceData(ids.DatasetPublished, edition, ids.InstanceEditionConfirmed, ids.UniqueTimestamp),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When a user sends a GET request for a list of dimensions for instance without sending a token", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions", ids.InstanceEditionConfirmed).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When a user sends a GET request for a list of dimensions for instance with an invalid token", func() {
			Convey("Then return status unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}/dimensions", ids.InstanceEditionConfirmed).
					WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("Failed to tear down test data", err, nil)
				os.Exit(1)
			}
		}
	})
}

func checkInstanceDimensionsResponse(response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("time")

	response.Value("items").Array().Element(0).Object().Value("label").Equal("")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("id").Equal("202.45")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Empty()

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("option").Equal("202.45")

	response.Value("items").Array().Element(0).Object().Value("node_id").Equal("")

	response.Value("items").Array().Element(1).Object().Value("dimension").Equal("aggregate")

	response.Value("items").Array().Element(1).Object().Value("label").Equal("CPI (Overall Index)")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code").Object().Value("id").Equal("cpi1dimA19")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/cpi1dimA19$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("version").Object().Empty()

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(1).Object().Value("option").Equal("cpi1dimA19")

	response.Value("items").Array().Element(1).Object().Value("node_id").Equal("")

}

func getInstanceDimensionsSetup(datasetID, editionID, edition, instanceID string, uniqueTimestamp bson.MongoTimestamp) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	datasetDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID, uniqueTimestamp),
	}

	dimensionOneDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData("9811", instanceID),
	}

	dimensionTwoDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData("9812", instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
