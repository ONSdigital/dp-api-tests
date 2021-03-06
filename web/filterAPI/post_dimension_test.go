package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/web/filterAPI/expectedTestData"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyPostDimension(t *testing.T) {

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
	docs = append(docs, filter, instance)

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(docs...); err != nil {
			log.ErrorC("Unable to setup dimension option test resources", err, nil)
			os.Exit(1)
		}

		Convey("Add a dimension to the filter blueprint", func() {

			response := filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "Residence Type").
				WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON("Lives in a communal establishment", "Lives in a household"))).
				Expect().Status(http.StatusCreated).JSON().Object()

			validateDimensionResponse(*response, filterBlueprintID, "Residence Type")

			// Check data has been updated as expected
			filterBlueprint, err := mongo.GetFilter(cfg.MongoFiltersDB, collection, "filter_id", filterBlueprintID)
			if err != nil {
				log.ErrorC("Unable to retrieve updated document", err, nil)
			}

			// Set these empty objects to nil to be able to compare other fields
			filterBlueprint.Downloads = nil
			filterBlueprint.Events = nil

			So(len(filterBlueprint.Dimensions), ShouldEqual, 5)

			// Check dimension has been added to the end of the array
			So(filterBlueprint.Dimensions[4].Name, ShouldEqual, "Residence Type")

			expectedfilterBlueprint := expectedTestData.ExpectedFilterBlueprint(cfg.FilterAPIURL, datasetID, filterBlueprintID)
			expectedfilterBlueprint.InstanceID = instanceID
			expectedfilterBlueprint.FilterID = filterBlueprintID

			So(filterBlueprint.UniqueTimestamp, ShouldNotBeEmpty)
			filterBlueprint.UniqueTimestamp = 0

			So(filterBlueprint, ShouldResemble, expectedfilterBlueprint)
		})

		Convey("Add a dimension to the filter blueprint with duplicate options", func() {

			filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "Residence Type").
				WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON("Lives in a communal establishment", "Lives in a communal establishment", "Lives in a household"))).
				Expect().Status(http.StatusCreated)

			// Check data has been updated as expected
			filterBlueprint, err := mongo.GetFilter(cfg.MongoFiltersDB, collection, "filter_id", filterBlueprintID)
			if err != nil {
				log.ErrorC("Unable to retrieve updated document", err, nil)
			}

			// Set these empty objects to nil to be able to compare other fields
			filterBlueprint.Downloads = nil
			filterBlueprint.Events = nil

			So(len(filterBlueprint.Dimensions), ShouldEqual, 5)

			// Check dimension has been added to the end of the array
			So(filterBlueprint.Dimensions[4].Name, ShouldEqual, "Residence Type")

			expectedfilterBlueprint := expectedTestData.ExpectedFilterBlueprint(cfg.FilterAPIURL, datasetID, filterBlueprintID)
			expectedfilterBlueprint.InstanceID = instanceID
			expectedfilterBlueprint.FilterID = filterBlueprintID

			So(filterBlueprint.UniqueTimestamp, ShouldNotBeEmpty)
			filterBlueprint.UniqueTimestamp = 0

			So(filterBlueprint, ShouldResemble, expectedfilterBlueprint)
		})

		Convey("Overwrite a dimension that already exists on a filter blueprint", func() {
			ageOptions := `{"options": ["28"]}`

			filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "age").
				WithBytes([]byte(ageOptions)).
				Expect().Status(http.StatusCreated)

			filterBlueprint, err := mongo.GetFilter(cfg.MongoFiltersDB, collection, "filter_id", filterBlueprintID)
			if err != nil {
				log.ErrorC("Unable to retrieve updated document", err, nil)
			}

			// Check data has been updated as expected
			So(len(filterBlueprint.Dimensions), ShouldEqual, 4)

			for _, dimension := range filterBlueprint.Dimensions {
				if dimension.Name == "age" {
					So(dimension.Options, ShouldResemble, []string{"28"})
					break
				}
			}
		})

		if err := mongo.Teardown(docs...); err != nil {
			log.ErrorC("Unable to remove instance test resource from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToPostDimension(t *testing.T) {

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
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, true),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	dimensions := setupMultipleDimensionsAndOptions(instanceID)

	docs = append(docs, dimensions...)
	docs = append(docs, filter, instance)

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a request is made to add a new dimension to the filter blueprint", func() {
			Convey("Then the response returns a status not found (404)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "Residence Type").
					WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
					Expect().Status(http.StatusNotFound).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given an existing filter blueprint", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When the request body is invalid", func() {
			Convey("Then the request body is invalid and returns status bad request (400)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "Residence Type").
					WithBytes([]byte(GetInvalidPOSTDimensionToFilterBlueprintJSON())).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body")
			})
		})

		Convey("When the version associated with filter blueprint no longer exists", func() {
			Convey("Then the response returns a status unprocessable entity (422)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "Residence Type").
					WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
					Expect().Status(http.StatusUnprocessableEntity).Body().Contains("version for filter blueprint no longer exists\n")
			})
		})

		Convey("And the version associated with filter blueprint exists", func() {

			if err := mongo.Setup(instance); err != nil {
				log.ErrorC("Unable to setup instance test resource", err, nil)
				os.Exit(1)
			}

			if err := mongo.Setup(dimensions...); err != nil {
				log.ErrorC("Unable to setup dimension option test resources", err, nil)
				os.Exit(1)
			}

			Convey("When the dimension does not exist against version", func() {
				Convey("Then the response returns a status bad request (400)", func() {

					filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "foobar").
						WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
						Expect().Status(http.StatusBadRequest).
						Body().Contains("incorrect dimensions chosen: [foobar]")
				})
			})

			Convey("When the option for a valid dimension does not exist against version", func() {

				if err := mongo.Setup(dimensions...); err != nil {
					log.ErrorC("Unable to setup dimension option test resources", err, nil)
					os.Exit(1)
				}

				Convey("Then the response returns a status bad request (400)", func() {

					filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "age").
						WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON("Lives in a communal establishment", "Lives in a household"))).
						Expect().Status(http.StatusBadRequest).Body().Contains("incorrect dimension options chosen: [Lives in a communal establishment Lives in a household]\n")
				})
			})
		})

		if err := mongo.Teardown(docs...); err != nil {
			log.ErrorC("Unable to remove test resources from mongo db", err, nil)
			os.Exit(1)
		}
	})

}
