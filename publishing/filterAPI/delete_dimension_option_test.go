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

func TestSuccessfullyDeleteDimensionOptions(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	docs := []*mongo.Doc{instance, filter}

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("Remove an option to a dimension to filter on and verify options are removed", func() {

			filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1T60000", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK)
			filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/time/options/April 1997", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK)

			// TODO call mongo directly instead of using API to get dimension options
			filterAPI.GET("/filters/{filter_blueprint_id}/dimensions", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK).JSON().Array().Length().Equal(4)

			filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK)
			sexDimResponse := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK).JSON().Array()
			goodsAndServicesDimResponse := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK).JSON().Array()
			timeDimResponse := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options", filterBlueprintID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				Expect().Status(http.StatusOK).JSON().Array()

			sexDimResponse.Element(0).Object().Value("option").NotEqual("male").Equal("female")

			goodsAndServicesDimResponse.Element(0).Object().Value("option").Equal("cpi1dim1S10201")
			goodsAndServicesDimResponse.Element(1).Object().Value("option").Equal("cpi1dim1S10105")
			timeDimResponse.Element(0).Object().Value("option").NotEqual("April 1997").Equal("March 1997")
			timeDimResponse.Element(1).Object().Value("option").NotEqual("April 1997").Equal("June 1997")
			timeDimResponse.Element(2).Object().Value("option").NotEqual("April 1997").Equal("September 1997")
			timeDimResponse.Element(3).Object().Value("option").NotEqual("April 1997").Equal("December 1997")
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToDeleteDimensionOptions(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	Convey("Given filter job does not exist", t, func() {
		Convey("When requesting to delete an option from the filter job", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/wages/options/27000", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("Filter blueprint not found")
			})
		})
	})

	Convey("Given a filter job", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to delete an option from a dimension that does not exist against the filter job", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/wages/options/27000", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).
						Body().Contains("Dimension not found")
			})
		})

		Convey("When requesting to delete an option that does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/age/options/44", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).
						Body().Contains("Option not found")
			})
		})

		if err := mongo.Teardown(filter); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
