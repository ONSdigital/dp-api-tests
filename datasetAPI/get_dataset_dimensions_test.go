package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestGetDimensions_ReturnsAllDimensionsFromADataset(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	edition := "2017"

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	dimensionOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData(instanceID),
	}
	dimensionTwoDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimension.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData(instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.SetupMany(d); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of all dimensions of a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", datasetID, edition, 1).WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			checkDimensionsResponse(datasetID, edition, instanceID, response)
		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", datasetID, edition, 1).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)

			checkDimensionsResponse(datasetID, edition, instanceID, response)
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

// TODO Remove skipped tests when code has been refactored (and hence fixed)
func TestGetDimensions_Failed(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	unpublishedInstancID := uuid.NewV4().String()

	edition := "2017"

	var docs []mongo.Doc

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	editionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData(datasetID, editionID, edition),
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData(datasetID, edition, instanceID),
	}

	docs = append(docs, datasetDoc, editionDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.SetupMany(d); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Fail to get a list of Dimensions for a dataset", t, func() {
		Convey("When authenticated", func() {

			// TODO Remove skip on test once endpoint fixed
			SkipConvey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", "1234", "2018").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Dataset not found\n")
			})

			// TODO Remove skip on test once endpoint fixed
			SkipConvey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, "2018").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Edition not found\n")
			})

			Convey("When there are no versions", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/3/dimensions", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Version not found\n")
			})
		})

		Convey("When unauthenticated", func() {

			// TODO Remove skip on test once endpoint fixed
			SkipConvey("When the dataset does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", "1234", "2018").
					Expect().Status(http.StatusBadRequest).Body().Contains("Dataset not found\n")
			})

			// TODO Remove skip on test once endpoint fixed
			SkipConvey("When the edition does not exist", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, "2018").
					Expect().Status(http.StatusBadRequest).Body().Contains("Edition not found\n")
			})

			Convey("When there are no published versions", func() {
				// Create an unpublished instance document
				mongo.Setup(database, "instances", "_id", unpublishedInstancID, validEditionConfirmedInstanceData(datasetID, edition, unpublishedInstancID))
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/3/dimensions", datasetID, edition).
					Expect().Status(http.StatusBadRequest).Body().Contains("Version not found\n")

				mongo.Teardown(database, "instances", "_id", unpublishedInstancID)
			})
		})
	})

	Convey("Given a valid dataset id, edition and version with no dimensions", t, func() {
		Convey("When authenticated and get the dimensions", func() {
			Convey("Then the error code should be 404", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, edition).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimensions not found\n")
			})
		})

		Convey("When unauthenticated and get the dimensions", func() {
			Convey("Then the error code should be 404", func() {
				datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/1/dimensions", datasetID, edition).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimensions not found\n")
			})
		})
	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}

func checkDimensionsResponse(datasetID, edition, instanceID string, response *httpexpect.Object) {

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("id").Equal("aggregate")
	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/aggregate/options$")

	response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("items").Array().Element(0).Object().Value("dimension").Equal("aggregate")
	response.Value("items").Array().Element(0).Object().Value("description").Equal("An aggregate of the data")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("options").Object().Value("id").Equal("time")
	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/time/options$")

	response.Value("items").Array().Element(1).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

	response.Value("items").Array().Element(1).Object().Value("dimension").Equal("time")
	response.Value("items").Array().Element(1).Object().Value("description").Equal("The time in which this dataset spans")
}
