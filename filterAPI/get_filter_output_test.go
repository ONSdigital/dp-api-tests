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

func TestSuccessfullyGetFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter output with downloads", t, func() {

		update := GetValidFilterOutputWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID)

		if err := mongo.Setup(database, "filterOutputs", "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get filter output", func() {
			Convey("Then filter output is returned in the response body", func() {

				response := filterAPI.GET("/filter-outputs/{filter_output_id}", filterOutputID).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("dimensions").Array().Length().Equal(4)
				response.Value("dimensions").Array().Element(0).Object().NotContainsKey("dimension_url") // Check dimension url is not set
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("dimensions").Array().Element(0).Object().Value("options").Equal([]string{"27"})
				response.Value("downloads").Object().Value("csv").Object().Value("url").Equal("s3-csv-location")
				response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("12mb")
				response.Value("downloads").Object().Value("json").Object().Value("url").Equal("s3-json-location")
				response.Value("downloads").Object().Value("json").Object().Value("size").Equal("6mb")
				response.Value("downloads").Object().Value("xls").Object().Value("url").Equal("s3-xls-location")
				response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
				response.Value("filter_id").Equal(filterOutputID)
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("filter_blueprint").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "$")
				response.Value("links").Object().Value("filter_blueprint").Object().Value("id").Equal(filterBlueprintID)
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filter-outputs/" + filterOutputID + "$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/123/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
			})
		})

		if err := mongo.Teardown(database, "filterOutputs", "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToGetFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given filter output does not exist", t, func() {
		Convey("When requesting to get filter output", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.GET("/filter-outputs/{filter_output_id}", filterID).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter output not found\n")
			})
		})
	})
}
