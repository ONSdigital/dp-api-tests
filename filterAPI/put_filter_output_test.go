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

func TestSuccessfulPutFilterOutput(t *testing.T) {

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

		Convey("When an authorised request to update the filter output with csv download without a public link", func() {

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
				Expect().Status(http.StatusOK)

			Convey("Then the filter output resource contains a non empty csv download object", func() {

				// Check data has been updated as expected
				filterOutput, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput.Downloads.CSV, ShouldNotBeNil)
				So(filterOutput.Downloads.CSV.HRef, ShouldEqual, "download-service-url.csv")
				So(filterOutput.Downloads.CSV.Private, ShouldEqual, "private-s3-csv-location")
				So(filterOutput.Downloads.CSV.Public, ShouldBeEmpty)
				So(filterOutput.Downloads.CSV.Size, ShouldEqual, "12mb")
				So(filterOutput.Downloads.XLS, ShouldBeNil)
				So(filterOutput.State, ShouldEqual, "created")
			})
		})

		Convey("When an authorised request to update the filter output with xls download without a public link", func() {

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				WithBytes([]byte(GetValidPUTFilterOutputWithXLSDownloadJSON())).
				Expect().Status(http.StatusOK)

			Convey("Then the filter output resource contains a non empty xls download object", func() {

				// Check data has been updated as expected
				filterOutput, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput.Downloads.CSV, ShouldNotBeNil)
				So(filterOutput.Downloads.CSV.HRef, ShouldEqual, "download-service-url.csv")
				So(filterOutput.Downloads.CSV.Private, ShouldEqual, "private-s3-csv-location")
				So(filterOutput.Downloads.CSV.Size, ShouldEqual, "12mb")
				So(filterOutput.Downloads.XLS, ShouldNotBeNil)
				So(filterOutput.Downloads.XLS.HRef, ShouldEqual, "download-service-url.xlsx")
				So(filterOutput.Downloads.XLS.Private, ShouldEqual, "private-s3-xls-location")
				So(filterOutput.Downloads.XLS.Size, ShouldEqual, "24mb")
				So(filterOutput.State, ShouldEqual, "completed")
			})
		})

		Convey("When an authorised request to update the filter output with csv public download link and"+
			"a second authorised request to update the filter output with xls public download link", func() {

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				WithBytes([]byte(GetValidPUTFilterOutputWithCSVPublicLinkJSON())).
				Expect().Status(http.StatusOK)

			filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
				WithHeader(serviceAuthTokenName, serviceAuthToken).
				WithBytes([]byte(GetValidPUTFilterOutputWithXLSPublicLinkJSON())).
				Expect().Status(http.StatusOK)

			Convey("Then the filter output resource contains a non empty csv and xls download objects and the state is set to `completed`", func() {

				// Check data has been updated as expected
				filterOutput, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}
				log.Debug("filter output", log.Data{"filter_output": filterOutput.Downloads})

				So(filterOutput.Downloads.CSV, ShouldNotBeNil)
				So(filterOutput.Downloads.CSV.HRef, ShouldEqual, "download-service-url.csv")
				So(filterOutput.Downloads.CSV.Private, ShouldEqual, "private-s3-csv-location")
				So(filterOutput.Downloads.CSV.Public, ShouldEqual, "public-s3-csv-location")
				So(filterOutput.Downloads.CSV.Size, ShouldEqual, "12mb")
				So(filterOutput.Downloads.XLS, ShouldNotBeNil)
				So(filterOutput.Downloads.XLS.HRef, ShouldEqual, "download-service-url.xlsx")
				So(filterOutput.Downloads.XLS.Private, ShouldEqual, "private-s3-xls-location")
				So(filterOutput.Downloads.XLS.Public, ShouldEqual, "public-s3-xls-location")
				So(filterOutput.Downloads.XLS.Size, ShouldEqual, "24mb")
				So(filterOutput.State, ShouldEqual, "completed")
			})
		})
	})

	if err := mongo.Teardown(output); err != nil {
		log.ErrorC("Unable to remove test data from mongo db", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPutFilterOutput(t *testing.T) {

	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
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

	dimensions := goodsAndServicesDimension(cfg.FilterAPIURL, filterBlueprintID)
	csvPublicLink := "public-s3-csv-location"
	xlsPublicLink := "public-s3-xls-location"

	publishedFilterOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, csvPublicLink, xlsPublicLink, version, dimensions),
	}

	unpublishedFilterOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     GetValidFilterOutputBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, "", "", version, dimensions),
	}

	Convey("Given a filter output does not exist", t, func() {
		Convey("When an authorised request is made to update filter output", func() {
			Convey("Then the request fails and returns status not found (404)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithHeader(serviceAuthTokenName, serviceAuthToken).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusNotFound).Body().Contains("Filter output not found\n")
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
			Convey("Then fail to update filter output and return status not found (404)", func() {

				filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
					WithBytes([]byte(GetValidPUTFilterOutputWithCSVDownloadJSON())).
					Expect().Status(http.StatusNotFound).Body().Contains("resource not found\n")
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

	Convey("Given an existing filter output with downloads object", t, func() {
		Convey("But the document is unpublished and does not have a csv or xls publish link", func() {
			if err := mongo.Setup(unpublishedFilterOutput); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an authorised request is made to update a filter output downloads", func() {
				Convey("Then fail to update filter output and return status forbidden (403)", func() {

					filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
						WithHeader(serviceAuthTokenName, serviceAuthToken).
						WithBytes([]byte(`{"downloads":{"csv":{"private":"private-s3-csv-link"}}}`)).
						Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [downloads.csv.private]\n")

					filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
						WithHeader(serviceAuthTokenName, serviceAuthToken).
						WithBytes([]byte(`{"downloads":{"xls":{"private":"private-s3-xls-link"}}}`)).
						Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [downloads.xls.private]\n")
				})
			})

			if err := mongo.Teardown(unpublishedFilterOutput); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})

		Convey("And the document is published", func() {
			if err := mongo.Setup(publishedFilterOutput); err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("When an authorised request is made to update a filter output downloads", func() {
				Convey("Then fail to update filter output and return status forbidden (403)", func() {

					filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
						WithHeader(serviceAuthTokenName, serviceAuthToken).
						WithBytes([]byte(`{"downloads":{"csv":{"private":"private-s3-csv-link"}}}`)).
						Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [downloads.csv]\n")

					filterAPI.PUT("/filter-outputs/{filter_output_id}", filterOutputID).
						WithHeader(serviceAuthTokenName, serviceAuthToken).
						WithBytes([]byte(`{"downloads":{"xls":{"private":"private-s3-xls-link"}}}`)).
						Expect().Status(http.StatusForbidden).Body().Contains("Forbidden from updating the following fields: [downloads.xls]\n")
				})
			})

			if err := mongo.Teardown(publishedFilterOutput); err != nil {
				log.ErrorC("Unable to remove test data from mongo db", err, nil)
				os.Exit(1)
			}
		})
	})
}
