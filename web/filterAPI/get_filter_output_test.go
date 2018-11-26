package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/davecgh/go-spew/spew"
)

func TestSuccessfullyGetFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()
	publishedFilterOutputID := uuid.NewV4().String()
	unpublishedFilterOutputID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	publishedOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, publishedFilterOutputID, filterBlueprintID, datasetID, edition, version, true),
	}

	unpublishedOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, unpublishedFilterOutputID, filterBlueprintID, datasetID, edition, version, false),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetUnpublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	Convey("Given an existing public filter output with downloads", t, func() {

		if err := mongo.Setup(publishedOutput); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get filter output", func() {
			Convey("Then filter output is returned in the response body without private or public download links", func() {

				response := filterAPI.GET("/filter-outputs/{filter_output_id}", publishedFilterOutputID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dataset").Object().Value("id").Equal(datasetID)
				response.Value("dataset").Object().Value("edition").Equal(edition)
				response.Value("dataset").Object().Value("version").Equal(version)
				response.Value("dimensions").Array().Length().Equal(4)
				response.Value("dimensions").Array().Element(0).Object().NotContainsKey("dimension_url") // Check dimension url is not set
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("dimensions").Array().Element(0).Object().Value("options").Equal([]string{"27"})
				response.Value("downloads").Object().Value("csv").Object().Value("href").Equal("download-service-url.csv")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("12mb")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("xls").Object().Value("href").Equal("download-service-url.xlsx")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
				response.Value("filter_id").Equal(publishedFilterOutputID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("filter_blueprint").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "$")
				response.Value("links").Object().Value("filter_blueprint").Object().Value("id").Equal(filterBlueprintID)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filter-outputs/" + publishedFilterOutputID + "$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("state").Equal("completed")
			})
		})

		if err := mongo.Teardown(publishedOutput); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

	Convey("Given an unpublished instance, and an existing pre-publish filter output with downloads", t, func() {

		if err := mongo.Setup(instance, unpublishedOutput); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When making an unauthenticated request to get filter output", func() {
			Convey("Then filter output is not found", func() {

				filterAPI.GET("/filter-outputs/{filter_output_id}", unpublishedFilterOutputID).
					Expect().Status(http.StatusNotFound).Body().Contains(filterOutputNotFoundResponse)

			})
		})

		Convey("When making a request to get filter output with a download service token header", func() {
			Convey("Then filter output is returned in the response body", func() {

				response := filterAPI.GET("/filter-outputs/{filter_output_id}", unpublishedFilterOutputID).
					WithHeader(common.DownloadServiceHeaderKey, downloadServiceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dataset").Object().Value("id").Equal(datasetID)
				response.Value("dataset").Object().Value("edition").Equal(edition)
				response.Value("dataset").Object().Value("version").Equal(version)
				response.Value("dimensions").Array().Length().Equal(4)
				response.Value("dimensions").Array().Element(0).Object().NotContainsKey("dimension_url") // Check dimension url is not set
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("dimensions").Array().Element(0).Object().Value("options").Equal([]string{"27"})
				response.Value("downloads").Object().Value("csv").Object().Value("href").Equal("download-service-url.csv")
				response.Value("downloads").Object().Value("csv").Object().Value("public").Equal("https://s3-eu-west-1.amazonaws.com/dp-frontend-florence-file-uploads/2470609-cpicoicoptestcsv")
				response.Value("downloads").Object().Value("csv").Object().Value("private").Equal("private-s3-csv-location")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("12mb")
				response.Value("downloads").Object().Value("xls").Object().Value("href").Equal("download-service-url.xlsx")
				response.Value("downloads").Object().Value("xls").Object().Value("public").Equal("public-s3-xls-location")
				response.Value("downloads").Object().Value("xls").Object().Value("private").Equal("private-s3-xls-location")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
				response.Value("filter_id").Equal(unpublishedFilterOutputID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("filter_blueprint").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "$")
				response.Value("links").Object().Value("filter_blueprint").Object().Value("id").Equal(filterBlueprintID)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filter-outputs/" + unpublishedFilterOutputID + "$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("state").Equal("completed")
			})
		})

		Convey("When the instance has been published and a request is made with no authentication to get filter output", func() {

			instance.Update = GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version)

			if err := mongo.Setup(instance); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("Then filter output is returned in the response body", func() {

				log.Info("\n\n ------------------------ ", nil)
				spew.Dump(mongo.GetInstance(	cfg.MongoDB, "instances", "instance_id", instanceID))
				log.Info(" ------------------------\n\n ", nil)


				response := filterAPI.GET("/filter-outputs/{filter_output_id}", unpublishedFilterOutputID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dataset").Object().Value("id").Equal(datasetID)
				response.Value("dataset").Object().Value("edition").Equal(edition)
				response.Value("dataset").Object().Value("version").Equal(version)
				response.Value("dimensions").Array().Length().Equal(4)
				response.Value("dimensions").Array().Element(0).Object().NotContainsKey("dimension_url") // Check dimension url is not set
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("dimensions").Array().Element(0).Object().Value("options").Equal([]string{"27"})
				response.Value("downloads").Object().Value("csv").Object().Value("href").Equal("download-service-url.csv")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("csv").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("12mb")
				response.Value("downloads").Object().Value("xls").Object().Value("href").Equal("download-service-url.xlsx")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("private")
				response.Value("downloads").Object().Value("xls").Object().NotContainsKey("public")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
				response.Value("filter_id").Equal(unpublishedFilterOutputID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("filter_blueprint").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "$")
				response.Value("links").Object().Value("filter_blueprint").Object().Value("id").Equal(filterBlueprintID)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filter-outputs/" + unpublishedFilterOutputID + "$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.Value("state").Equal("completed")
			})
		})

		if err := mongo.Teardown(instance, unpublishedOutput); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToGetFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)
	datasetID := "test-cpih01"
	edition := "2017"
	version := 1

	unpublishedOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, version, false),
	}

		filterBlueprint := &mongo.Doc{
			Database:   cfg.MongoFiltersDB,
			Collection: collection,
			Key:        "_id",
			Value:      filterID,
			Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, false),
		}

		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "instance_id",
			Value:      instanceID,
			Update:     GetUnpublishedInstanceDataBSON(instanceID, datasetID, edition, version),
		}


	Convey("Given filter output does not exist", t, func() {
		Convey("When requesting to get filter output", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filter-outputs/{filter_output_id}", filterID).
					Expect().Status(http.StatusNotFound).Body().Contains(filterOutputNotFoundResponse)
			})
		})
	})

	Convey("Given an unpublished instance, and an existing pre-publish filter output with downloads", t, func() {

		if err := mongo.Setup(instance, filterBlueprint, unpublishedOutput); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When making an unauthenticated request to get filter output", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filter-outputs/{filter_output_id}", filterOutputID).
					Expect().Status(http.StatusNotFound).Body().Contains(filterOutputNotFoundResponse)
			})
		})

		if err := mongo.Teardown(instance, filterBlueprint, unpublishedOutput); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
