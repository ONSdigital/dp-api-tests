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

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions", filterBlueprintID).
					Expect().Status(http.StatusOK).JSON().Array()

				expectedDimensions := []string{"age", "sex", "aggregate", "time"}
				for i, dim := range expectedDimensions {

					filterURL := cfg.FilterAPIURL + "/filters/" + filterBlueprintID
					selfURL := filterURL + "/dimensions/" + dim
					optionsURL := selfURL + "/options"

					links := response.Element(i).Object().Value("links").Object()

					So(response.Element(i).Object().Value("name").Raw(), ShouldEqual, dim)
					So(links.Value("filter").Object().Value("href").Raw(), ShouldEqual, filterURL)
					So(links.Value("filter").Object().Value("id").Raw(), ShouldEqual, filterBlueprintID)
					So(links.Value("self").Object().Value("href").Raw(), ShouldEqual, selfURL)
					So(links.Value("self").Object().Value("id").Raw(), ShouldEqual, dim)
					So(links.Value("options").Object().Value("href").Raw(), ShouldEqual, optionsURL)
				}
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
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
