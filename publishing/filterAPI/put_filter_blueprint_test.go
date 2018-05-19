package filterAPI

import (
	"net/http"
	"os"
	"testing"
	"time"

	datasetJSON "github.com/ONSdigital/dp-api-tests/web/datasetAPI"
	"github.com/ONSdigital/dp-api-tests/publishing/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulPutFilterBlueprint(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	newInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	newInstance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      newInstanceID,
		Update:     GetValidPublishedInstanceDataBSON(newInstanceID, datasetID, edition, 2),
	}

	docs := setupMultipleDimensionsAndOptions(instanceID)
	newInstanceDimensionDocs := setupMultipleDimensionsAndOptions(newInstanceID)
	docs = append(docs, newInstanceDimensionDocs...)
	docs = append(docs, filter, dataset, editionDoc, instance, newInstance)

	Convey("Given an existing filter blueprint", t, func() {

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Unable to setup dimension option test resources", err, nil)
			os.Exit(1)
		}

		Convey("When a request to update the filter blueprint with an info event and new version of the same edition and dataset", func() {

			time := time.Now()

			response := filterAPI.PUT("/filters/{filter_id}", filterBlueprintID).
				WithBytes([]byte(GetValidPUTFilterBlueprintJSON(2, time))).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then the response contains the updated filter blueprint", func() {

				// check response contains the correct data
				response.Value("instance_id").Equal(newInstanceID)
				response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/filters/" + filterBlueprintID + "/dimensions$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/filters/" + filterBlueprintID + "$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.NotContainsKey("last_updated")
				response.NotContainsKey("downloads")

				info := response.Value("events").Object().Value("info").Array()
				info.Length().Equal(1)
				info.Element(0).Object().Value("message").Equal("blueprint has created filter output resource")
				info.Element(0).Object().Value("time").Equal(time.String())
				info.Element(0).Object().Value("type").Equal("info")
			})
		})

		Convey("When a request to submit the filter blueprint occurs", func() {

			response := filterAPI.PUT("/filters/{filter_id}", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				WithQuery("submitted", "true").
				WithBytes([]byte(`{}`)).
				Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then filter blueprint creates a filter output document and in the response is a link to this resource", func() {

				filterOutputLinkObject := response.Value("links").Object()
				filterOutputLinkObject.Value("filter_output").Object().Value("href").String().Match("(.+)/filter-outputs/(.+)$")
				filterOutputLinkObject.Value("filter_output").Object().Value("id").NotNull()

				filterOutputID := filterOutputLinkObject.Value("filter_output").Object().Value("id").String().Raw()

				filterOutput, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput, ShouldResemble, expectedTestData.ExpectedFilterOutput(cfg.FilterAPIURL, instanceID, filterOutputID, filterBlueprintID))

				//enable teardown of resources created during test
				docs = append(docs, &mongo.Doc{
					Database:   cfg.MongoFiltersDB,
					Collection: "filterOutputs",
					Key:        "filter_id",
					Value:      filterOutputID,
				})
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToPutFilterBlueprint(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a post request is made to update filter blueprint", func() {
			Convey("Then the request fails and returns status not found (404)", func() {

				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetValidPUTUpdateFilterBlueprintJSON(instanceID))).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given an existing filter blueprint", t, func() {

		if err := mongo.Setup(filter, instance); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an invalid json body is sent to update filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {

				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetInvalidSyntaxJSON(instanceID))).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When a put request to change the dataset against a filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {
				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetInValidPUTFilterBlueprintJSON(datasetID, "", version, time.Now()))).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When a put request to change the edition against a filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {
				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetInValidPUTFilterBlueprintJSON("", edition, version, time.Now()))).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When a put request to change the dataset and edition against a filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {
				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetInValidPUTFilterBlueprintJSON(datasetID, edition, version, time.Now()))).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When a put request to change the version to a non existing one against a filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {
				newVersion := 2
				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetValidPUTFilterBlueprintJSON(newVersion, time.Now()))).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains(versionNotFoundResponse)
			})
		})

		if err := mongo.Teardown(filter, instance); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
