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

func TestSuccessfullyGetListOfDimensions(t *testing.T) {

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

	Convey("Given an existing filter", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting a list of dimensions in a filter blueprint", func() {
			Convey("Then return a list of all dimensions for filter blueprint", func() {

				actual := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Array()

				checkDimensionJson(actual.Element(0), filterBlueprintID, "age")
				checkDimensionJson(actual.Element(1), filterBlueprintID, "sex")
				checkDimensionJson(actual.Element(2), filterBlueprintID, "aggregate")
				checkDimensionJson(actual.Element(3), filterBlueprintID, "time")
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func checkDimensionJson(dimensionJSON *httpexpect.Value, filterBlueprintID string, dimension string) {

	dimensionJSON.Object().Value("links").Object().Value("filter").Object().Value("href").Equal("http://localhost:22100/filters/" + filterBlueprintID)
	dimensionJSON.Object().Value("links").Object().Value("filter").Object().Value("id").Equal(filterBlueprintID)
	dimensionJSON.Object().Value("links").Object().Value("options").Object().Value("href").Equal("http://localhost:22100/filters/" + filterBlueprintID + "/dimensions/" + dimension + "/options")
	dimensionJSON.Object().Value("links").Object().Value("self").Object().Value("href").Equal("http://localhost:22100/filters/" + filterBlueprintID + "/dimensions/" + dimension)
	dimensionJSON.Object().Value("links").Object().Value("self").Object().Value("id").Equal(dimension)
	dimensionJSON.Object().Value("name").Equal(dimension)
}

func TestFailureToGetListOfDimensions(t *testing.T) {

	filterBlueprintID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When requesting a list of dimensions for filter blueprint", func() {
			Convey("Then return status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions", filterBlueprintID).
					Expect().Status(http.StatusNotFound).Body().Contains(filterNotFoundResponse)
			})
		})
	})
}
