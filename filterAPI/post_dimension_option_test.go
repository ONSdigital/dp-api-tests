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

func TestSuccessfulPostDimensionOptions(t *testing.T) {

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

		filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/28", filterBlueprintID).
			Expect().Status(http.StatusCreated)
		// Add a duplicate option, this should not be added into the dimension option
		filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/28", filterBlueprintID).
			Expect().Status(http.StatusCreated)
		filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/sex/options/unknown", filterBlueprintID).
			Expect().Status(http.StatusCreated)
		filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/Goods and services/options/welfare", filterBlueprintID).
			Expect().Status(http.StatusCreated)
		filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/time/options/February 2007", filterBlueprintID).
			Expect().Status(http.StatusCreated)

		filterJob, err := mongo.GetFilter(database, collection, "filter_id", filterBlueprintID)
		if err != nil {
			log.ErrorC("Unable to retrieve updated document", err, nil)
		}

		// Set downloads empty object to nil to be able to compare other fields
		filterJob.Downloads = nil

		expectedFilterJob := expectedTestData.ExpectedFilterBlueprintUpdated(cfg.FilterAPIURL, instanceID, filterBlueprintID)
		expectedFilterJob.InstanceID = instanceID
		expectedFilterJob.FilterID = filterBlueprintID

		So(filterJob, ShouldResemble, expectedFilterJob)
	})

	if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostDimensionOptions(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter blueprint does not exist", t, func() {
		invalidfilterBlueprintID := "12345678"

		Convey("When a post request to add an option to a dimension for that filter blueprint", func() {
			Convey("Then return status bad request (400)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/age/options/30", invalidfilterBlueprintID).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - filter blueprint not found")
			})
		})
	})

	Convey("Given a filter blueprint exists", t, func() {

		update := GetValidCreatedFilterBlueprintBSON(cfg.FilterAPIURL, filterID, instanceID, filterBlueprintID)

		if err := mongo.Setup(database, collection, "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a post request to add an option for a dimension that does not exist against that filter blueprint", func() {
			Convey("Then return status not found (404)", func() {

				filterAPI.POST("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
