package filterAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetFilterOutputPreview(t *testing.T) {
	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter output exists", t, func() {

		update := GetValidFilterOutputWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID)

		if err := mongo.Setup(database, "filterOutputs", "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}
		graphData := neo4j.NewDatastore("bolt://localhost:7687", instanceID, neo4j.ObservationTestData)

		if err := graphData.Setup(); err != nil {
			log.ErrorC("Unable to setup graph data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get a preview for the filter output", func() {
			Convey("Then the filtered preview is returned in the response body", func() {
				response := filterAPI.GET("/filter-outputs/{filter_output_id}/preview", filterOutputID).
					Expect().Status(http.StatusOK).JSON().Object()
				response.Value("rows").Array().Length().Equal(3)
				response.Value("headers").Array().Length().Equal(7)
			})
		})

		mongo.Teardown(database, "filterOutputs", "_id", filterID)
		graphData.TeardownInstance()
	})
}

func TestErrorCasesGetFilterOutputPreview(t *testing.T) {
	filterID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given an existing filter output exists", t, func() {

		update := GetValidFilterOutputBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, nil)

		if err := mongo.Setup(database, "filterOutputs", "_id", filterID, update); err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}
		graphData := neo4j.NewDatastore("bolt://localhost:7687", instanceID, neo4j.ObservationTestData)

		if err := graphData.Setup(); err != nil {
			log.ErrorC("Unable to setup graph data", err, nil)
			os.Exit(1)
		}

		Convey("When requesting to get a preview with no dimensions", func() {
			Convey("Then the filtered preview is returned in the response body", func() {
				filterAPI.GET("/filter-outputs/{filter_output_id}/preview", filterOutputID).
					Expect().Status(http.StatusBadRequest)
			})
		})

		Convey("When requesting to get a preview with limit query as letters", func() {
			Convey("Then the filtered preview is returned in the response body", func() {
				filterAPI.GET("/filter-outputs/{filter_output_id}/preview", filterOutputID).WithQuery("limit", "abc").
					Expect().Status(http.StatusBadRequest)
			})
		})

		Convey("When requesting to get a preview with invalid filter ouput id", func() {
			Convey("Then the filtered preview is returned in the response body", func() {
				filterAPI.GET("/filter-outputs/{filter_output_id}/preview", "123-321").
					Expect().Status(http.StatusNotFound)
			})
		})

		mongo.Teardown(database, "filterOutputs", "_id", filterID)
		graphData.TeardownInstance()
	})
}
