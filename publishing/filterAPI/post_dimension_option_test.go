package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/publishing/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulPostDimensionOptions(t *testing.T) {
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
	docs := setupMultipleDimensionsAndOptions(instanceID)
	docs = append(docs, instance, filter)

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Unable to setup instance test resource", err, nil)
			os.Exit(1)
		}

		response := filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/28", filterBlueprintID).
			WithHeader(serviceAuthTokenName, serviceAuthToken).
			Expect().Status(http.StatusCreated).JSON().Object()
		validateOptionResponse(*response, filterBlueprintID, "age", "28")

		response = filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/sex/options/unknown", filterBlueprintID).
			WithHeader(serviceAuthTokenName, serviceAuthToken).
			Expect().Status(http.StatusCreated).JSON().Object()
		validateOptionResponse(*response, filterBlueprintID, "sex", "unknown")
		response = filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1T60000", filterBlueprintID).
			WithHeader(serviceAuthTokenName, serviceAuthToken).
			Expect().Status(http.StatusCreated).JSON().Object()
		validateOptionResponse(*response, filterBlueprintID, "aggregate", "cpi1dim1T60000")
		response = filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/time/options/February 2007", filterBlueprintID).
			WithHeader(serviceAuthTokenName, serviceAuthToken).
			Expect().Status(http.StatusCreated).JSON().Object()
		validateOptionResponse(*response, filterBlueprintID, "time", "February 2007")

		// Re-adding an existing option should not change the dimension option
		response = filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/28", filterBlueprintID).
			WithHeader(serviceAuthTokenName, serviceAuthToken).
			Expect().Status(http.StatusCreated).JSON().Object()
		validateOptionResponse(*response, filterBlueprintID, "age", "28")

		filterJob, err := mongo.GetFilter(cfg.MongoFiltersDB, collection, "filter_id", filterBlueprintID)
		if err != nil {
			log.ErrorC("Unable to retrieve updated document", err, nil)
		}

		// Set downloads empty object to nil to be able to compare other fields
		filterJob.Downloads = nil

		expectedFilterJob := expectedTestData.ExpectedFilterBlueprintUpdated(cfg.FilterAPIURL, filterBlueprintID, datasetID)
		expectedFilterJob.InstanceID = instanceID
		expectedFilterJob.FilterID = filterBlueprintID

		So(filterJob.UniqueTimestamp, ShouldNotBeEmpty)
		filterJob.UniqueTimestamp = 0

		So(filterJob, ShouldResemble, expectedFilterJob)
	})

	if err := mongo.Teardown(docs...); err != nil {
		log.ErrorC("Unable to remove instance test resource from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostDimensionOptions(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	var docs []*mongo.Doc

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidCreatedFilterBlueprintBSON(cfg.FilterAPIURL, filterID, instanceID, filterBlueprintID, datasetID, edition, version),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	option := setupDimensionOptions(uuid.NewV4().String(), GetValidAgeDimensionData(instanceID, "27"))
	docs = append(docs, filter, instance, option)

	Convey("Given filter blueprint does not exist", t, func() {
		invalidfilterBlueprintID := uuid.NewV4().String()

		Convey("When a post request to add an option to a dimension for that filter blueprint", func() {
			Convey("Then return status bad request", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/30", invalidfilterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given a filter blueprint exists", t, func() {

		if err := mongo.Setup(filter, option); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a post request to add an option for a dimension for a version of a dataset that does not exist", func() {
			Convey("Then return status unprocessable entity (422)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusUnprocessableEntity).
					Body().Contains("version for filter blueprint no longer exists\n")
			})
		})

		Convey("And the version that is associated with this filter blueprint does exist", func() {
			if err := mongo.Setup(instance); err != nil {
				log.ErrorC("Unable to setup instance test resource", err, nil)
				os.Exit(1)
			}

			Convey("When a post request to add an option for a dimension that does not exist", func() {
				Convey("Then return status not found (400)", func() {

					filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
						WithHeader(serviceAuthTokenName, serviceAuthToken).
						Expect().Status(http.StatusBadRequest).
						Body().Contains("incorrect dimensions chosen: [sex]\n")
				})
			})

			Convey("When a post request to add an option that does not exist for a dimension", func() {
				if err := mongo.Setup(option); err != nil {
					log.ErrorC("Unable to setup instance test resource", err, nil)
					os.Exit(1)
				}

				Convey("Then return status not found (400)", func() {

					filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/75", filterBlueprintID).
						WithHeader(serviceAuthTokenName, serviceAuthToken).
						Expect().Status(http.StatusBadRequest).Body().Contains("incorrect dimension options chosen: [75]\n")
				})
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			log.ErrorC("Unable to remove instance test resource from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
