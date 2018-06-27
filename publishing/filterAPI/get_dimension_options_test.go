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

func TestSuccessfullyGetListOfDimensionOptions(t *testing.T) {

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

	Convey("Given an existing filter blueprint with dimensions and options", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting a list of options for a dimension", func() {
			Convey("Then return a list of options for `age` dimension", func() {

				dimension := "age"
				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/"+dimension+"/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				checkDimensionOptionJson(response.Element(0), filterBlueprintID, dimension, "27")
			})

			Convey("Then return a list of options for `sex` dimension", func() {

				dimension := "sex"

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/"+dimension+"/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				checkDimensionOptionJson(response.Element(0), filterBlueprintID, dimension, "male")
				checkDimensionOptionJson(response.Element(1), filterBlueprintID, dimension, "female")
			})

			Convey("Then return a list of options for `goods and services` dimension", func() {

				dimension := "aggregate"

				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/"+dimension+"/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				checkDimensionOptionJson(response.Element(0), filterBlueprintID, dimension, "cpi1dim1T60000")
				checkDimensionOptionJson(response.Element(1), filterBlueprintID, dimension, "cpi1dim1S10201")
				checkDimensionOptionJson(response.Element(2), filterBlueprintID, dimension, "cpi1dim1S10105")
			})

			Convey("Then return a list of options for `time` dimension", func() {

				dimension := "time"
				response := filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/"+dimension+"/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).JSON().Array()

				checkDimensionOptionJson(response.Element(0), filterBlueprintID, dimension, "March 1997")
				checkDimensionOptionJson(response.Element(1), filterBlueprintID, dimension, "April 1997")
				checkDimensionOptionJson(response.Element(2), filterBlueprintID, dimension, "June 1997")
				checkDimensionOptionJson(response.Element(3), filterBlueprintID, dimension, "September 1997")
				checkDimensionOptionJson(response.Element(4), filterBlueprintID, dimension, "December 1997")
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func checkDimensionOptionJson(json *httpexpect.Value, filterBlueprintID string, dimension, option string) {

	json.Object().Value("links").Object().Value("filter").Object().Value("href").Equal("http://localhost:22100/filters/" + filterBlueprintID)
	json.Object().Value("links").Object().Value("filter").Object().Value("id").Equal(filterBlueprintID)
	json.Object().Value("links").Object().Value("dimension").Object().Value("href").Equal("http://localhost:22100/filters/" + filterBlueprintID + "/dimensions/" + dimension)
	json.Object().Value("links").Object().Value("dimension").Object().Value("id").Equal(dimension)
	json.Object().Value("links").Object().Value("self").Object().Value("href").Equal("http://localhost:22100/filters/" + filterBlueprintID + "/dimensions/" + dimension + "/options/" + option)
	json.Object().Value("links").Object().Value("self").Object().Value("id").Equal(option)
	json.Object().Value("option").Equal(option)
}

func TestFailureToGetListOfDimensionOptions(t *testing.T) {

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

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a request to get a dimension option against filter blueprint", func() {
			Convey("Then return status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(filterNotFoundResponse)
			})
		})
	})

	Convey("Given a filter blueprint", t, func() {

		if err := mongo.Setup(filter); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When a request to get a dimension option against filter blueprint where the dimension does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/wages/options", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(dimensionNotFoundResponse)
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}
