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

func TestSuccessfullyPostfilterBlueprint(t *testing.T) {

	instanceID := uuid.NewV4().String()

	update := GetValidPublishedInstanceDataBSON(instanceID)

	if err := setupInstance(instanceID, update); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a valid json input to create a filter", t, func() {
		Convey("Then the response returns a status of created (201)", func() {

			response := filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(instanceID))).
				Expect().Status(http.StatusCreated).JSON().Object()

			response.Value("filter_id").NotNull()
			response.Value("instance_id").Equal(instanceID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/(.+)/dimensions$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/(.+)$")
			response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/123/editions/2017/versions/1$")
			response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
		})
	})

	Convey("Given a valid json input to create a filter", t, func() {
		Convey("When request contains query parameter `submitted` set to `true`", func() {

			response := filterAPI.POST("/filters").
				WithQuery("submitted", "true").
				WithBytes([]byte(GetValidPOSTCreateFilterJSON(instanceID))).
				Expect().Status(http.StatusCreated).JSON().Object()

			filterBlueprintID := response.Value("filter_id").String().Raw()

			response.Value("filter_id").NotNull()
			response.Value("instance_id").Equal(instanceID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "/dimensions$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "$")
			response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/123/editions/2017/versions/1$")
			response.Value("links").Object().Value("version").Object().Value("id").Equal("1")

			Convey("Then filter blueprint is created, a filter output document is created and in the response is a link to the output resource", func() {

				filterOutputLinkObject := response.Value("links").Object()
				filterOutputLinkObject.Value("filter_output").Object().Value("id").NotNull()

				filterOutputID := filterOutputLinkObject.Value("filter_output").Object().Value("id").String().Raw()
				filterOutputLinkObject.Value("filter_output").Object().Value("href").String().Match("/filter-outputs/" + filterOutputID + "$")

				filterOutput, err := mongo.GetFilter(database, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput, ShouldResemble, expectedTestData.ExpectedFilterOutputOnPost(cfg.FilterAPIURL, instanceID, filterOutputID, filterBlueprintID))

				if err := mongo.Teardown(database, "filterOutputs", "filter_id", filterOutputID); err != nil {
					log.ErrorC("Unable to remove test data from mongo db", err, nil)
					os.Exit(1)
				}
			})

			if err := mongo.Teardown(database, collection, "filter_id", filterBlueprintID); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})
	})

	if err := mongo.Teardown(database, collection, "instance_id", instanceID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}

	if err := teardownInstance(instanceID); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostfilterBlueprint(t *testing.T) {

	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given invalid json input to create a filter", t, func() {
		Convey("Then the response returns status bad request (400)", func() {

			filterAPI.POST("/filters").WithBytes([]byte(GetInvalidJSON(instanceID))).
				Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
		})
	})

	Convey("Given a request to create a filter", t, func() {
		Convey("When the request body contains an instance id which does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(instanceID))).
					Expect().Status(http.StatusNotFound).Body().Contains("Instance not found\n")
			})
		})
	})
}
