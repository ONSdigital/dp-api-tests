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
	"time"
)

func TestSuccessfulPostFilterOutputEvent(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	output := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputWithoutDownloadsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, datasetID, edition, version),
	}

	if err := mongo.Setup(output); err != nil {
		log.ErrorC("Unable to setup test data", err, nil)
		os.Exit(1)
	}

	Convey("Given an existing filter output", t, func() {

		Convey("When an authenticated POST request is made to add an event to a filter output", func() {

			filterAPI.POST("/filter-outputs/{filter_output_id}/events", filterOutputID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				WithBytes([]byte(`{"type":"CSVCreated","time":"2018-06-10T05:59:05.893+01:00"}`)).
				Expect().Status(http.StatusCreated)

			Convey("Then the filter output resource the event", func() {

				// Check data has been updated as expected
				filterOutput, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput.Events, ShouldNotBeNil)
				So(len(filterOutput.Events), ShouldEqual, 1)
				So(filterOutput.Events[0].Type, ShouldEqual, "CSVCreated")

				eventTime := filterOutput.Events[0].Time.Format(time.RFC3339Nano)
				So(eventTime, ShouldEqual, "2018-06-10T05:59:05.893+01:00")
			})
		})
	})

	if err := mongo.Teardown(output); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostFilterOutputEvent(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	filterOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputWithoutDownloadsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, datasetID, edition, version),
	}

	Convey("Given a filter output does not exist", t, func() {
		Convey("When an authorised request is made to update filter output", func() {
			Convey("Then the request fails and returns status not found (404)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusNotFound).Body().Contains(filterOutputNotFoundResponse)
			})
		})
	})

	Convey("Given an existing filter output without downloads object", t, func() {

		if err := mongo.Setup(filterOutput); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised request with invalid json body is sent to update filter output", func() {
			Convey("Then fail to update filter output and return status bad request (400)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(`{`)).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When an unauthorised request is made to update filter output", func() {
			Convey("Then fail to update filter output and return status unauthorized (401)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an invalid authorization header is set on request to update filter output", func() {
			Convey("Then fail to update filter output and return status unauthorized (401)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(serviceAuthTokenName, invalidServiceAuthToken).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When an authorised request is made to update filter output and json body contains dimensions or an instance id", func() {
			Convey("Then fail to update filter output and return status forbidden (403)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(GetValidPUTFilterOutputWithDimensionsJSON())).
					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [dimensions]\n")

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(`{"instance_id": "1234"}`)).
					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [instance_id]")
			})
		})

		if err := mongo.Teardown(filterOutput); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}
