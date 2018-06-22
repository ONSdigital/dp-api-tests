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
	"fmt"
)

func TestSuccessfullyGetDimensionOption(t *testing.T) {

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

		Convey("When checking the dimension options", func() {
			Convey("Then return status ok (200) and expected response body for dimension `age` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "27", filterBlueprintID, "age", "27"))
			})

			Convey("Then return status ok (200) and expected response body for dimension `sex` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/male", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "male", filterBlueprintID, "sex", "male"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/female", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "female", filterBlueprintID, "sex", "female"))
			})

			Convey("Then return status ok (200) and expected response body for dimension `aggregate` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1S10201", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "cpi1dim1S10201", filterBlueprintID, "aggregate", "cpi1dim1S10201"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1S10105", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "cpi1dim1S10105", filterBlueprintID, "aggregate", "cpi1dim1S10105"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/aggregate/options/cpi1dim1T60000", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "cpi1dim1T60000", filterBlueprintID, "aggregate", "cpi1dim1T60000"))

			})

			Convey("Then return status no content (204) for dimension `time` options", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/March 1997", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "March 1997", filterBlueprintID, "time", "March 1997"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/April 1997", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "April 1997", filterBlueprintID, "time", "April 1997"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/June 1997", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "June 1997", filterBlueprintID, "time", "June 1997"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/September 1997", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "September 1997", filterBlueprintID, "time", "September 1997"))

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/time/options/December 1997", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusOK).Body().
					Contains(fmt.Sprintf(`self":{"id":"%s","href":"http://localhost:22100/filter/%s/dimensions/%s/options/%s"}`, "December 1997", filterBlueprintID, "time", "December 1997"))

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
			Convey("Then return a status 400 bad request", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/age/options/27", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusBadRequest).Body().Contains(filterNotFoundResponse)
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
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(dimensionNotFoundResponse)
			})
		})

		Convey("When a request to get a dimension option that does not exist", func() {
			Convey("Then return a status not found (404)", func() {

				filterAPI.GET("/filters/{filter_blueprint_id}/dimensions/sex/options/unknown", filterBlueprintID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					Expect().Status(http.StatusNotFound).Body().Contains(optionNotFoundResponse)
			})
		})
	})

	if err := mongo.Teardown(filter); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}
