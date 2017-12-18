package filterAPI

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-tests/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulPutFilterBlueprint(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	update := GetValidPublishedInstanceDataBSON(instanceID)

	if err := setupInstance(instanceID, update); err != nil {
		log.ErrorC("Unable to setup instance", err, nil)
		os.Exit(1)
	}

	Convey("Given an existing filter blueprint", t, func() {

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

		Convey("When a request to update the filter blueprint with an info event and new instance id", func() {

			newInstanceID := uuid.NewV4().String()
			time := time.Now()

			response := filterAPI.PUT("/filters/{filter_id}", filterBlueprintID).
				WithBytes([]byte(GetValidPUTFilterBlueprintJSON(newInstanceID, time))).
				Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then the response contains the updated filter blueprint", func() {

				// check response contains the correct data
				response.Value("instance_id").Equal(newInstanceID)
				response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/filters/" + filterBlueprintID + "/dimensions$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/filters/" + filterBlueprintID + "$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/123/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				response.NotContainsKey("last_updated")
				response.NotContainsKey("downloads")

				info := response.Value("events").Object().Value("info").Array()
				info.Length().Equal(1)
				info.Element(0).Object().Value("message").Equal("blueprint has created filter output resource")
				info.Element(0).Object().Value("time").Equal(time.String())
				info.Element(0).Object().Value("type").Equal("info")
			})
		})

		Convey("When a request to submit the filter blueprint occurs", func() {

			response := filterAPI.PUT("/filters/{filter_id}", filterBlueprintID).
				WithQuery("submitted", "true").
				WithBytes([]byte(`{}`)).
				Expect().Status(http.StatusOK).JSON().Object()

			Convey("Then filter blueprint creates a filter output document and in the response is a link to this resource", func() {

				filterOutputLinkObject := response.Value("links").Object()
				filterOutputLinkObject.Value("filter_output").Object().Value("href").String().Match("(.+)/filter-outputs/(.+)$")
				filterOutputLinkObject.Value("filter_output").Object().Value("id").NotNull()

				filterOutputID := filterOutputLinkObject.Value("filter_output").Object().Value("id").String().Raw()

				filterOutput, err := mongo.GetFilter(database, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput, ShouldResemble, expectedTestData.ExpectedFilterOutput(cfg.FilterAPIURL, instanceID, filterOutputID, filterBlueprintID))

				if err := mongo.Teardown(database, "filterOutputs", "filter_id", filterOutputID); err != nil {
					log.ErrorC("Unable to remove test data from mongo db", err, nil)
					os.Exit(1)
				}
			})
		})

		if err := teardownInstance(instanceID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}

		if err := teardownDimensionOptions(instanceID); err != nil {
			log.ErrorC("Unable to remove dimension option test resources from mongo db", err, nil)
			os.Exit(1)
		}

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToPutFilterBlueprint(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter blueprint does not exist", t, func() {
		Convey("When a post request is made to update filter blueprint", func() {
			Convey("Then the request fails and returns status not found (404)", func() {

				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetValidPUTUpdateFilterBlueprintJSON(instanceID))).
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

		Convey("When an invalid json body is sent to update filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {

				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetInvalidSyntaxJSON(instanceID))).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When a put request to change the instance id to a non existing one against a filter blueprint", func() {
			Convey("Then fail to update filter blueprint and return status bad request (400)", func() {

				filterAPI.PUT("/filters/{filter_blueprint_id}", filterBlueprintID).WithBytes([]byte(GetValidPUTFilterBlueprintJSON(instanceID, time.Now()))).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - instance not found\n")
			})
		})

		if err := mongo.Teardown(database, collection, "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
