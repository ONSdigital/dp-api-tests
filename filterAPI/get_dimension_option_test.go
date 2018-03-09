package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetDimensionOption(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterBlueprintID, true),
	}

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When checking the dimension options", func() {
			Convey("Then return status no content (204) for dimension `age` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for dimension `sex` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/female", filterBlueprintID).
					Expect().Status(http.StatusNoContent)

			})

			Convey("Then return status no content (204) for dimension `goods and services` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1S10201", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1S10105", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1T60000", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
			})

			Convey("Then return status no content (204) for dimension `time` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/March 1997", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/April 1997", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/June 1997", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/September 1997", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/December 1997", filterBlueprintID).
					Expect().Status(http.StatusNoContent)
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToGetDimensionOption(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filter := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: collection,
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterBlueprintID, true),
	}

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a request to get a dimension option against filter blueprint", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter blueprint not found")
			})
		})
	})

	Convey("Given a filter blueprint containing dimension options", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a request to get a dimension option where the dimension does not exist", func() {
			Convey("Then return a status bad request (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/ages/options/27", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})

		Convey("When a request to get a dimension option that does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/unknown", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains("Option not found")
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}
