package filterAPI

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	datasetJSON "github.com/ONSdigital/dp-api-tests/datasetAPI"
	"github.com/ONSdigital/dp-api-tests/filterAPI/expectedTestData"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyPostFilterBlueprintForPublishedInstance(t *testing.T) {

	instanceID := uuid.NewV4().String()
	dimensionOptionOneID := uuid.NewV4().String()
	dimensionOptionTwoID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	var docs []*mongo.Doc

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	docs = append(docs, dataset, editionDoc, instance)
	docs = append(docs, setupDimensionOptions(dimensionOptionOneID, GetValidAgeDimensionData(instanceID, "27")))
	docs = append(docs, setupDimensionOptions(dimensionOptionTwoID, GetValidAgeDimensionData(instanceID, "42")))

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Unable to setup test resources", err, nil)
		os.Exit(1)
	}

	Convey("Given a valid json input to create a filter", t, func() {
		Convey("Then the response returns a status of created (201)", func() {
			response := filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(datasetID, edition, version))).
				Expect().Status(http.StatusCreated).JSON().Object()

			response.Value("filter_id").NotNull()
			response.Value("instance_id").Equal(instanceID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/(.+)/dimensions$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/(.+)$")
			response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/" + strconv.Itoa(version) + "$")
			response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
		})
	})

	Convey("Given a valid json input to create a filter", t, func() {
		Convey("When request contains query parameter `submitted` set to `true`", func() {

			response := filterAPI.POST("/filters").
				WithQuery("submitted", "true").
				WithBytes([]byte(GetValidPOSTCreateFilterJSON(datasetID, edition, version))).
				Expect().Status(http.StatusCreated).JSON().Object()

			filterBlueprintID := response.Value("filter_id").String().Raw()

			response.Value("filter_id").NotNull()
			response.Value("instance_id").Equal(instanceID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "/dimensions$")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "$")
			response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/" + strconv.Itoa(version) + "$")
			response.Value("links").Object().Value("version").Object().Value("id").Equal("1")

			Convey("Then filter blueprint is created, a filter output document is created and in the response is a link to the output resource", func() {

				filterOutputLinkObject := response.Value("links").Object()
				filterOutputLinkObject.Value("filter_output").Object().Value("id").NotNull()

				filterOutputID := filterOutputLinkObject.Value("filter_output").Object().Value("id").String().Raw()
				filterOutputLinkObject.Value("filter_output").Object().Value("href").String().Match("/filter-outputs/" + filterOutputID + "$")

				filterOutput, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated document", err, nil)
				}

				So(filterOutput, ShouldResemble, expectedTestData.ExpectedFilterOutputOnPost(cfg.FilterAPIURL, datasetID, edition, instanceID, filterOutputID, filterBlueprintID, version))

				//enable teardown of resources created during test
				docs = append(docs, &mongo.Doc{
					Database:   cfg.MongoFiltersDB,
					Collection: "filterOutputs",
					Key:        "filter_id",
					Value:      filterOutputID,
				})

				docs = append(docs, &mongo.Doc{
					Database:   cfg.MongoFiltersDB,
					Collection: collection,
					Key:        "filter_id",
					Value:      filterBlueprintID,
				})
			})
		})
	})

	if err := mongo.Teardown(docs...); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostfilterBlueprintForPublishedInstance(t *testing.T) {

	instanceID := uuid.NewV4().String()
	dimensionOptionID := uuid.NewV4().String()
	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editonDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	dimension := setupDimensionOptions(dimensionOptionID, GetValidAgeDimensionData(instanceID, "27"))

	Convey("Given invalid json input to create a filter", t, func() {
		Convey("When the request body does not contain the dataset object details", func() {
			Convey("Then the response returns status bad request (400)", func() {

				filterAPI.POST("/filters").WithBytes([]byte(GetInvalidJSON())).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
			})
		})
	})

	Convey("Given a request to create a filter", t, func() {
		Convey("When the request body contains a dataset version which does not exist", func() {
			Convey("Then the response returns status not found (404)", func() {

				filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(datasetID, edition, version))).
					Expect().Status(http.StatusNotFound).Body().Contains("Version not found\n")
			})
		})
	})

	Convey("Given that a dataset version is published", t, func() {

		if err := mongo.Setup(dataset, editonDoc, instance, dimension); err != nil {
			log.ErrorC("Unable to setup dimension option", err, nil)
			os.Exit(1)
		}

		Convey("When the request contains a valid version of an editon for a dataset but a dimension that does not exist", func() {
			Convey("Then the response returns status bad request (400)", func() {

				filterAPI.POST("/filters").WithBytes([]byte(GetInvalidDimensionJSON(datasetID, edition, version))).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - incorrect dimensions chosen: [weight]\n")
			})
		})

		Convey("When the request contains a valid version of an editon for a dataset and dimension but dimension options is invalid", func() {
			Convey("Then the response returns status bad request (400)", func() {

				filterAPI.POST("/filters").WithBytes([]byte(GetInvalidDimensionOptionJSON(datasetID, edition, version))).
					Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - incorrect dimension options chosen: [33]\n")
			})
		})

		if err := mongo.Teardown(instance, dimension); err != nil {
			log.ErrorC("Unable to teardown instance", err, nil)
			os.Exit(1)
		}
	})
}

func TestPostFilterBlueprintForUnpublishedInstance(t *testing.T) {

	instanceID := uuid.NewV4().String()
	dimensionOptionOneID := uuid.NewV4().String()
	dimensionOptionTwoID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	var docs []*mongo.Doc

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedWithUpdatesDatasetData(datasetID),
	}

	editionDoc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      datasetID,
		Update:     datasetJSON.ValidPublishedEditionData(datasetID, editionID, edition),
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "instance_id",
		Value:      instanceID,
		Update:     GetUnpublishedInstanceDataBSON(instanceID, datasetID, edition, version),
	}

	docs = append(docs, dataset, editionDoc, instance)
	docs = append(docs, setupDimensionOptions(dimensionOptionOneID, GetValidAgeDimensionData(instanceID, "27")))
	docs = append(docs, setupDimensionOptions(dimensionOptionTwoID, GetValidAgeDimensionData(instanceID, "42")))

	if err := mongo.Setup(docs...); err != nil {
		log.ErrorC("Unable to setup instance test resources", err, nil)
		os.Exit(1)
	}

	Convey("Given an unpublished version", t, func() {
		Convey("When no authentication is provided on the POST request", func() {
			Convey("Then the response returns a status of not found (404)", func() {
				filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(datasetID, edition, version))).
					Expect().Status(http.StatusNotFound)
			})
		})

		Convey("When invalid authentication is provided on the POST request", func() {
			Convey("Then the response returns a status of not found (404)", func() {
				filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(datasetID, edition, version))).
					WithHeader(internalTokenHeader, "failure").Expect().Status(http.StatusNotFound)
			})
		})

		Convey("When valid authentication is provided on the POST request", func() {
			Convey("Then the response returns a status of created (201)", func() {
				response := filterAPI.POST("/filters").WithBytes([]byte(GetValidPOSTCreateFilterJSON(datasetID, edition, version))).
					WithHeader(internalTokenHeader, internalTokenID).Expect().Status(http.StatusCreated).JSON().Object()
				response.Value("filter_id").NotNull()
				response.Value("instance_id").Equal(instanceID)
				response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/(.+)/dimensions$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/(.+)$")
				response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/123/editions/2017/versions/1$")
				response.Value("links").Object().Value("version").Object().Value("id").Equal("1")

				//enable teardown of resources created during test
				docs = append(docs, &mongo.Doc{
					Database:   cfg.MongoFiltersDB,
					Collection: collection,
					Key:        "filter_id",
					Value:      response.Value("filter_id").String().Raw(),
				})
			})
		})
	})
	if err := mongo.Teardown(docs...); err != nil {
		log.ErrorC("Unable to teardown instance", err, nil)
		os.Exit(1)
	}
}
