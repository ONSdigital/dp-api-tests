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

func TestSuccessfulPutFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter output", t, func() {

		update := GetValidFilterOutputWithoutDownloadsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID)

		if err := mongo.Setup(database, "filterOutputs", "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised request to update the filter output with csv download", func() {

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
				Expect().Status(http.StatusOK)

			Convey("Then the filter output resource contains a non empty csv download object", func() {

				// Check data has been updated as expected
				filterOutput, err := mongo.GetFilter(database, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput.Downloads.CSV.URL, ShouldEqual, "s3-csv-location")
				So(filterOutput.Downloads.CSV.Size, ShouldEqual, "12mb")
				So(filterOutput.Downloads.XLS.URL, ShouldEqual, "")
				So(filterOutput.Downloads.XLS.Size, ShouldEqual, "")
				So(filterOutput.State, ShouldEqual, "created")
			})
		})

		Convey("When an authorised request to update the filter output with csv download and xls download", func() {

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
				Expect().Status(http.StatusOK)

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(GetValidPUTFilterOutputWithXLSDownloadJSON())).
				Expect().Status(http.StatusOK)

			Convey("Then the filter output resource contains a non empty csv and xls download objects and the state is set to `completed`", func() {

				// Check data has been updated as expected
				filterOutput, err := mongo.GetFilter(database, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput.Downloads.CSV.URL, ShouldEqual, "s3-csv-location")
				So(filterOutput.Downloads.CSV.Size, ShouldEqual, "12mb")
				So(filterOutput.Downloads.XLS.URL, ShouldEqual, "s3-xls-location")
				So(filterOutput.Downloads.XLS.Size, ShouldEqual, "24mb")
				So(filterOutput.State, ShouldEqual, "completed")
			})
		})

		if err := teardownInstance(instanceID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}

		if err := mongo.Teardown(database, "filterOutputs", "_id", filterID); err != nil {
			log.ErrorC("Unable to remove test data from mongo db", err, nil)
			os.Exit(1)
		}
	})
}

func TestFailureToPutFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a filter output does not exist", t, func() {
		Convey("When an authorised request is made to update filter output", func() {
			Convey("Then the request fails and returns status not found (404)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(internalTokenHeader, internalTokenID).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter output not found\n")
			})
		})
	})

	Convey("Given an existing filter output", t, func() {

		update := GetValidFilterOutputWithoutDownloadsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID)

		if err := mongo.Setup(database, "filterOutputs", "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}

		Convey("When an authorised request with invalid json body is sent to update filter output", func() {
			Convey("Then fail to update filter output and return status bad request (400)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(internalTokenHeader, internalTokenID).
					WithBytes([]byte(`{`)).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})

		Convey("When an unauthorised request is made to update filter output", func() {
			Convey("Then fail to update filter output and return status unauthorised (401)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusUnauthorized).Body().Contains("Unauthorised, request lacks valid authentication credentials\n")
			})
		})

		Convey("When an authorised request is made to update filter output and json body contains dimensions or an instance id", func() {
			Convey("Then fail to update filter output and return status forbidden (403)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(internalTokenHeader, internalTokenID).
					WithBytes([]byte(GetValidPUTFilterOutputWithDimensionsJSON())).
					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [dimensions]\n")

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(internalTokenHeader, internalTokenID).
					WithBytes([]byte(`{"instance_id": "1234"}`)).
					Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [instance_id]")
			})
		})
	})

	if err := mongo.Teardown(database, "filterOutputs", "_id", filterID); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}
