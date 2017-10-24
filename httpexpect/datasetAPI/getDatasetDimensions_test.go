package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDimensions_ReturnsAllDimensionsFromADataset(t *testing.T) {
	var docs []mongo.Doc

	// datasetOneDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "datasets",
	// 	Key:        "_id",
	// 	Value:      datasetID,
	// 	Update:     validPublishedDatasetData,
	// }

	// datasetTwoDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "datasets",
	// 	Key:        "_id",
	// 	Value:      datasetID,
	// 	Update:     validUnpublishedDatasetData,
	// }

	// editionOneDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "editions",
	// 	Key:        "_id",
	// 	Value:      editionID,
	// 	Update:     validPublishedEditionData,
	// }

	// editionTwoDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "editions",
	// 	Key:        "_id",
	// 	Value:      editionID,
	// 	Update:     validUnpublishedEditionData,
	// }

	// dimensionOneDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimensions",
	// 	Key:        "_id",
	// 	Value:      dimensionID,
	// 	Update:     validTimeDimensionsData,
	// }
	// dimensionTwoDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimensions",
	// 	Key:        "_id",
	// 	Value:      dimensionID,
	// 	Update:     validSexDimensionsData,
	// }

	// dimensionTimeOptionsDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimension.options",
	// 	Key:        "_id",
	// 	Value:      dimensionOptionID,
	// 	Update:     validTimeDimensionsOptionsData,
	// }

	// dimensionSexOptionsDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimension.options",
	// 	Key:        "_id",
	// 	Value:      dimensionOptionID,
	// 	Update:     validSexDimensionsOptionsData,
	// }

	// instanceOneDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "instances",
	// 	Key:        "_id",
	// 	Value:      instanceID,
	// 	Update:     validPublishedInstanceData,
	// }

	// instanceTwoDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "instances",
	// 	Key:        "_id",
	// 	Value:      "799",
	// 	Update:     validUnpublishedInstanceData,
	// }

	datasetDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData,
	}

	editionDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEditionData,
	}

	instanceOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     validPublishedInstanceData,
	}

	dimensionOneDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimensions.options",
		Key:        "_id",
		Value:      "9811",
		Update:     validTimeDimensionsData,
	}
	dimensionTwoDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "dimensions.options",
		Key:        "_id",
		Value:      "9812",
		Update:     validAggregateDimensionsData,
	}

	// dimensionTimeOptionsDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimension.options",
	// 	Key:        "_id",
	// 	Value:      dimensionOptionID,
	// 	Update:     validTimeDimensionsOptionsData,
	// }

	// dimensionSexOptionsDoc := mongo.Doc{
	// 	Database:   "datasets",
	// 	Collection: "dimension.options",
	// 	Key:        "_id",
	// 	Value:      dimensionOptionID,
	// 	Update:     validSexDimensionsOptionsData,
	// }

	docs = append(docs, datasetDoc, editionDoc, dimensionOneDoc, dimensionTwoDoc, instanceOneDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	mongo.TeardownMany(d)

	if err := mongo.SetupMany(d); err != nil {
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get a list of all dimensions of a dataset", t, func() {
		Convey("When user is authenticated", func() {

			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", datasetID, edition, 1).WithHeader("internal-token", "FD0108EA-825D-411C-9B1D-41EF7727F465").
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)

			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("id").Equal("sex")
			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/sex/options$")

			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

			// response.Value("items").Array().Element(0).Object().Value("dimension_id").Equal("sex")
		})

		Convey("When a user is not authenticated", func() {
			response := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/dimensions", datasetID, edition, 1).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("items").Array().Length().Equal(2)
			// fmt.Println(response.Raw())

			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("id").Equal("64d384f1-ea3b-445c-8fb8-aa453f96e58a")
			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("code_list").Object().Value("href").String().Match("(.+)/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a$")

			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("id").Equal("sex")
			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("options").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions/sex/options$")

			// response.Value("items").Array().Element(0).Object().Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/instances/" + instanceID + "$")

			// response.Value("items").Array().Element(0).Object().Value("dimension_id").Equal("sex")
			//	checkVersionResponse(response, 0)
		})
	})

	//	mongo.TeardownMany(d)
}
