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

func TestSuccessfullyDeleteDimension(t *testing.T) {

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

	Convey("Given an existing filter blueprint", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When sending a delete request to remove an existing dimension on the filter blueprint", func() {
			Convey("Then the filter blueprint should not contain that dimension", func() {
				filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/aggregate", filterBlueprintID).
					Expect().Status(http.StatusOK)

				var expectedDimensions []mongo.Dimension

				dimensionAge := mongo.Dimension{
					URL:     cfg.FilterAPIURL + "/filters/" + filterBlueprintID + "/dimensions/age",
					Name:    "age",
					Options: []string{"27"},
				}

				dimensionSex := mongo.Dimension{
					URL:     cfg.FilterAPIURL + "/filters/" + filterBlueprintID + "/dimensions/sex",
					Name:    "sex",
					Options: []string{"male", "female"},
				}

				dimensionTime := mongo.Dimension{
					URL:     cfg.FilterAPIURL + "/filters/" + filterBlueprintID + "/dimensions/time",
					Name:    "time",
					Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
				}

				expectedDimensions = append(expectedDimensions, dimensionAge, dimensionSex, dimensionTime)

				// Check dimension has been removed from filter blueprint
				filterBlueprint, err := mongo.GetFilter(cfg.MongoFiltersDB, collection, "filter_id", filterBlueprintID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
					os.Exit(1)
				}

				So(filterBlueprint.Dimensions, ShouldResemble, expectedDimensions)
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToDeleteDimension(t *testing.T) {

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

	Convey("Given filter blueprint does not exist", t, func() {
		Convey("When requesting to delete a dimension from filter blueprint", func() {
			Convey("Then response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/age", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter blueprint not found")
			})
		})
	})

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to delete a dimension from filter blueprint where the dimension does not exist", func() {
			Convey("Then response returns status not found (404)", func() {

				filterAPI.DELETE("/filters/{filter_blueprint_id}/dimensions/wage", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains("Dimension not found")
			})
		})

		if err := mongo.Teardown(filter); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
