package datasetAPI

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

// NOTE If endpoint is only available on publishing, remember to add a test to
// web/datasetAPI/hidden_endpoints_test.go to check request returns 404

func TestSuccessfullyPostInstance(t *testing.T) {
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an authorised user wants to create an instance", t, func() {
		Convey("When a valid authorised POST request is made with a job properties", func() {
			Convey("Then the expected response body is returned and a status of created (201)", func() {

				timeStart := time.Now().Truncate(time.Second).UTC()

				response := datasetAPI.POST("/instances").
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTCreateInstanceJSON)).
					Expect().Status(http.StatusCreated).JSON().Object()

				instanceUniqueID := response.Value("id").String().Raw()

				response.Value("id").NotNull()
				response.Value("links").Object().Value("job").Object().Value("id").Equal("042e216a-7822-4fa0-a3d6-e3f5248ffc35")
				response.Value("links").Object().Value("job").Object().Value("href").String().Match("/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/instances/" + instanceUniqueID + "$")
				response.Value("state").Equal("created")
				response.Value("last_updated").NotNull()
				response.NotContainsKey("unique_timestamp")

				// ensure DB has a unique_timestamp which is current
				instanceFromDB, err := mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceUniqueID)
				if err != nil {
					if err != mgo.ErrNotFound {
						log.ErrorC("Was unable to retrieve test data", err, nil)
						os.Exit(1)
					}
				}
				So(instanceFromDB.InstanceID, ShouldEqual, instanceUniqueID)
				So(instanceFromDB.UniqueTimestamp.Time().UTC(), ShouldHappenOnOrBetween, timeStart, time.Now().UTC())

				instance := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: "instances",
					Key:        "_id",
					Value:      instanceUniqueID,
				}

				if err := mongo.Teardown(instance); err != nil {
					if err != mgo.ErrNotFound {
						os.Exit(1)
					}
				}
			})
		})

		Convey("When a valid authorised POST request is made with dimensions and dataset as well as a job properties", func() {
			Convey("Then the expected response body is returned and a status of created (201)", func() {

				response := datasetAPI.POST("/instances").
					WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(validPOSTCreateFullInstanceJSON)).
					Expect().Status(http.StatusCreated).JSON().Object()

				instanceUniqueID := response.Value("id").String().Raw()

				response.Value("id").NotNull()
				response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("The age ranging from 16 to 75+")
				response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/code-lists/43513D18-B4D8-4227-9820-492B2971E7T5$")
				response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("43513D18-B4D8-4227-9820-492B2971E7T5")
				response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
				response.Value("links").Object().Value("job").Object().Value("id").Equal("042e216a-7822-4fa0-a3d6-e3f5248ffc35")
				response.Value("links").Object().Value("job").Object().Value("href").String().Match("/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35$")
				response.Value("links").Object().Value("dataset").Object().Value("id").Equal("34B13D18-B4D8-4227-9820-492B2971E221")
				response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/34B13D18-B4D8-4227-9820-492B2971E221$")
				response.Value("links").Object().Value("self").Object().Value("href").String().Match("/instances/" + instanceUniqueID + "$")
				response.Value("state").Equal("created")
				response.Value("last_updated").NotNull()
				response.NotContainsKey("unique_timestamp")

				instance := &mongo.Doc{
					Database:   cfg.MongoDB,
					Collection: "instances",
					Key:        "id",
					Value:      instanceUniqueID,
				}

				if err := mongo.Teardown(instance); err != nil {
					if err != mgo.ErrNotFound {
						os.Exit(1)
					}
				}
			})
		})
	})
}

func TestFailureToPostInstance(t *testing.T) {
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an authorised user wants to create an instance", t, func() {
		Convey("When an authorised POST request is made to create an instance resource with invalid json", func() {
			Convey("Then fail to create resource and return a status bad request (400) with a message `failed to parse json body: unexpected end of JSON input`", func() {

				datasetAPI.POST("/instances").WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(`{`)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("failed to parse json body")
			})
		})
	})

	Convey("Given an authorised user wants to create an instance", t, func() {
		Convey("When an authorised POST request is made to create an instance resource with missing job properties", func() {
			Convey("Then fail to create resource and return a status bad request (400) with a message `missing job properties`", func() {

				datasetAPI.POST("/instances").WithHeader(florenceTokenName, florenceToken).
					WithBytes([]byte(invalidPOSTCreateInstanceJSON)).
					Expect().Status(http.StatusBadRequest).
					Body().Contains("missing job properties")
			})
		})
	})

	Convey("Given an unauthorised user wants to create an instance", t, func() {
		Convey("When an unauthorised POST request is made to create an instance resource with an invalid authentication header", func() {
			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.POST("/instances").WithHeader(florenceTokenName, unauthorisedAuthToken).
					WithBytes([]byte(validPOSTCreateInstanceJSON)).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		Convey("When no authentication header is provided in POST request to create an instance resource", func() {
			Convey("Then fail to create resource and return a status unauthorized (401)", func() {

				datasetAPI.POST("/instances").WithBytes([]byte(validPOSTCreateInstanceJSON)).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})
}
