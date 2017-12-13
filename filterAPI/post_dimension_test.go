package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyPostDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter", t, func() {

		update := GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterBlueprintID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		if err := setupInstance(instanceID, GetValidPublishedInstanceDataBSON(instanceID)); err != nil {
			log.ErrorC("Unable to setup instance test resource", err, nil)
			os.Exit(1)
		}

		if err := setupMultipleDimensionsAndOptions(instanceID); err != nil {
			log.ErrorC("Unable to setup dimension option test resources", err, nil)
			os.Exit(1)
		}

		Convey("Add a dimension to the filter blueprint", func() {

			filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/{dimension}", filterBlueprintID, "Residence Type").
				WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
				Expect().Status(http.StatusCreated)

			// Check data has been updated as expected
			filterBlueprint, err := mongo.GetFilter(database, collection, "filter_id", filterBlueprintID)
			if err != nil {
				log.ErrorC("Unable to retrieve updated document", err, nil)
			}

			// Set these empty objects to nil to be able to compare other fields
			filterBlueprint.Downloads = nil
			filterBlueprint.Events = nil

			So(len(filterBlueprint.Dimensions), ShouldEqual, 5)

			// Check dimension has been added to the end of the array
			So(filterBlueprint.Dimensions[4].Name, ShouldEqual, "Residence Type")

			expectedfilterBlueprint := expectedTestData.ExpectedFilterBlueprint(cfg.FilterAPIURL, instanceID, filterBlueprintID)
			expectedfilterBlueprint.InstanceID = instanceID
			expectedfilterBlueprint.FilterID = filterBlueprintID

			So(filterBlueprint, ShouldResemble, expectedfilterBlueprint)
		})

		Convey("Overwrite a dimension that already exists on a filter blueprint", func() {
			ageOptions := `{"options": ["28"]}`

			filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age", filterBlueprintID).
				WithBytes([]byte(ageOptions)).
				Expect().Status(http.StatusCreated)

			filterBlueprint, err := mongo.GetFilter(database, collection, "filter_id", filterBlueprintID)
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
		if err := teardownDimensionOptions(instanceID); err != nil {
			log.ErrorC("Unable to remove dimension option test resources from mongo db", err, nil)
			os.Exit(1)
		}

		if err := teardownInstance(instanceID); err != nil {
			log.ErrorC("Unable to remove instance test resource from mongo db", err, nil)
			os.Exit(1)
		}

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToPostDimension(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a request is made to add a new dimension to the filter blueprint", func() {
			Convey("Then the response returns a status not found (404)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/Residence Type", filterBlueprintID).
					WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter blueprint not found\n")
			})
		})
	})

	Convey("Given an existing filter blueprint", t, func() {

		update := GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterBlueprintID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When the request body is invalid", func() {
			Convey("Then the request body is invalid and returns status bad request (400)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/Residence Type", filterBlueprintID).
					WithBytes([]byte(GetInvalidPOSTDimensionToFilterBlueprintJSON())).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body")
			})
		})

		Convey("When the instance does not exist for filter blueprint", func() {
			Convey("Then the response returns a status unprocessable entity (422)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/Residence Type", filterBlueprintID).
					WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
					Expect().Status(http.StatusUnprocessableEntity).Body().Contains("Unprocessable entity - instance for filter blueprint no longer exists\n")
			})
		})

		Convey("And the instance exists", func() {
			if err := setupInstance(instanceID, GetValidPublishedInstanceDataBSON(instanceID)); err != nil {
				log.ErrorC("Unable to setup instance test resource", err, nil)
				os.Exit(1)
			}

			Convey("When the dimension does not exist against instance", func() {
				Convey("Then the response returns a status bad request (400)", func() {

					filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/foobar", filterBlueprintID).
						WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
						Expect().Status(http.StatusBadRequest).Body().Contains("Dimension not found\n")
				})
			})

			Convey("When the option for a valid dimension does not exist against instance", func() {
				if err := setupMultipleDimensionsAndOptions(instanceID); err != nil {
					log.ErrorC("Unable to setup dimension option test resources", err, nil)
					os.Exit(1)
				}
				Convey("Then the response returns a status bad request (400)", func() {

					filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age", filterBlueprintID).
						WithBytes([]byte(GetValidPOSTDimensionToFilterBlueprintJSON())).
						Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - incorrect dimension options chosen: [Lives in a communal establishment Lives in a household]\n")
				})
				if err := teardownDimensionOptions(instanceID); err != nil {
					log.ErrorC("Unable to remove dimension option test resources from mongo db", err, nil)
					os.Exit(1)
				}
			})

			if err := teardownInstance(instanceID); err != nil {
				log.ErrorC("Unable to remove instance test resource from mongo db", err, nil)
				os.Exit(1)
			}
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})

}
