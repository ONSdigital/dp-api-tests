package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetFilterBlueprint(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	unpublishedFilterBlueprintID := uuid.NewV4().String()
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

	unpublishedFilter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, unpublishedFilterBlueprintID, version, false),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetUnpublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get filter blueprint", func() {
			Convey("Then filter blueprint is returned in the response body", func() {

				response := filterAPI.GET("/filters/{filter_blueprint_id}", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("filter_id").Equal(filterBlueprintID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "/dimensions$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/(.+)$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("published").Equal(true)
			})
		})

		if err := mongo.Teardown(filter); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given an existing filter for an unpublished instance", t, func() {

		if err := mongo.Setup(unpublishedFilter, instance); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get filter blueprint with authentication", func() {
			Convey("Then filter blueprint is returned in the response body", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}", unpublishedFilterBlueprintID).
					Expect().Status(http.StatusNotFound)
			})
		})

		Convey("When the instance is published and a request is made to get filter blueprint", func() {
			instance.Update = GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version)

			if err := mongo.Setup(instance); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("Then filter blueprint is returned in the response body", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}", unpublishedFilterBlueprintID).
					Expect().Status(http.StatusNotFound)
			})
		})

		if err := mongo.Teardown(unpublishedFilter, instance); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToGetFilterBlueprint(t *testing.T) {

	filterID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	unpublishedFilterBlueprintID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1
	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	unpublishedFilter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, unpublishedFilterBlueprintID, version, false),
	}

	Convey("Given filter blueprint does not exist", t, func() {
		Convey("When requesting to get filter blueprint", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}", filterID).
					Expect().Status(http.StatusNotFound).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given an existing filter for an unpublished instance", t, func() {

		if err := mongo.Setup(unpublishedFilter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get filter blueprint without authentication", func() {
			Convey("Then the response returns status not found (404)", func() {
				filterAPI.GET("/filters/{filter_blueprint_id}", unpublishedFilterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains(versionNotFoundResponse)
			})
		})

		if err := mongo.Teardown(unpublishedFilter); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
