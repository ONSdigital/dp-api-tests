package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"strconv"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetInstance(t *testing.T) {

	publishedInstanceID := uuid.NewV4().String()
	unpublishedInstanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given a published instance resource", t, func() {
		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      publishedInstanceID,
			Update:     validPublishedInstanceData(datasetID, edition, publishedInstanceID),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated request to get instance", func() {
			Convey("Then response contains the expected json object and a status ok (200)", func() {
				response := datasetAPI.GET("/instances/{id}", publishedInstanceID).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				response.Value("alerts").Array().Element(0).Object().Value("date").String().Equal("2017-12-10")
				response.Value("alerts").Array().Element(0).Object().Value("description").String().Equal("A correction to an observation for males of age 25, previously 11 now changed to 12")
				response.Value("alerts").Array().Element(0).Object().Value("type").String().Equal("Correction")

				response.Value("state").Equal("published")

				checkResponse(datasetID, edition, publishedInstanceID, "1", response)
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})

	Convey("Given an unpublished instance resource", t, func() {
		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      unpublishedInstanceID,
			Update:     validAssociatedInstanceData(datasetID, edition, unpublishedInstanceID),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an authenticated request to get instance", func() {
			Convey("Then response contains the expected json object and a status ok (200)", func() {
				response := datasetAPI.GET("/instances/{id}", unpublishedInstanceID).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusOK).JSON().Object()

				checkResponse(datasetID, edition, unpublishedInstanceID, "2", response)

				response.Value("state").Equal("associated")

			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToGetInstance(t *testing.T) {

	instanceID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	edition := "2017"

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Given an instance resource does not exist", t, func() {
		Convey("When an authorised request is made to get instance", func() {
			Convey("Then return a status not found (404) with message `Instance not found`", func() {
				datasetAPI.GET("/instances/{id}", instanceID).WithHeader(florenceTokenName, florenceToken).
					Expect().Status(http.StatusNotFound).Body().Contains("Instance not found")
			})
		})
	})

	Convey("Given a published instance resource exists", t, func() {
		instance := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "instances",
			Key:        "_id",
			Value:      instanceID,
			Update:     validPublishedInstanceData(datasetID, edition, instanceID),
		}

		if err := mongo.Setup(instance); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
		Convey("When no authentication header is provided in request to get resource", func() {
			Convey("Then return a status of not found (404) with message `requested resource not found`", func() {

				datasetAPI.GET("/instances/{id}", instanceID).
					Expect().Status(http.StatusNotFound).Body().Contains("requested resource not found")
			})
		})

		Convey("When an unauthorised request is made to get resource", func() {
			Convey("Then return a status of unauthorized (401)", func() {

				datasetAPI.GET("/instances/{id}", instanceID).WithHeader(florenceTokenName, unauthorisedAuthToken).
					Expect().Status(http.StatusUnauthorized)
			})
		})

		if err := mongo.Teardown(instance); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func checkResponse(datasetID, edition, instanceID, version string, response *httpexpect.Object) {
	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	response.Value("id").Equal(instanceID)
	response.Value("dimensions").Array().Element(0).Object().Value("description").Equal("A list of ages between 18 and 75+")
	response.Value("dimensions").Array().Element(0).Object().Value("href").String().Match("/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC$")
	response.Value("dimensions").Array().Element(0).Object().Value("id").Equal("408064B3-A808-449B-9041-EA3A2F72CFAC")
	response.Value("dimensions").Array().Element(0).Object().Value("name").Equal("age")
	response.Value("downloads").Object().Value("csv").Object().Value("href").String().Match("/aws/census-2017-" + version + "-csv$")
	response.Value("downloads").Object().Value("csv").Object().Value("private").String().Match("/private/myfile.csv$")
	response.Value("downloads").Object().Value("csv").Object().Value("public").String().Match("/public/myfile.csv$")
	response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10")
	response.Value("downloads").Object().Value("xls").Object().Value("href").String().Match("/aws/census-2017-" + version + "-xls$")
	response.Value("downloads").Object().Value("xls").Object().Value("private").String().Match("/private/myfile.xls$")
	response.Value("downloads").Object().Value("xls").Object().Value("public").String().Match("/public/myfile.xls$")
	response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24")
	response.Value("edition").Equal(edition)
	response.Value("headers").Array().Length().Equal(7)
	response.Value("headers").Array().Element(0).String().Equal("v4_0")
	response.Value("headers").Array().Element(1).String().Equal("time")
	response.Value("headers").Array().Element(2).String().Equal("time")
	response.Value("headers").Array().Element(3).String().Equal("uk-only")
	response.Value("headers").Array().Element(4).String().Equal("geography")
	response.Value("headers").Array().Element(5).String().Equal("cpi1dim1aggid")
	response.Value("headers").Array().Element(6).String().Equal("aggregate")
	response.Value("latest_changes").Array().Element(0).Object().Value("description").String().Equal("The border of Southampton changed after the south east cliff face fell into the sea.")
	response.Value("latest_changes").Array().Element(0).Object().Value("name").String().Equal("Changes in Classification")
	response.Value("latest_changes").Array().Element(0).Object().Value("type").String().Equal("Summary of Changes")
	response.Value("links").Object().Value("job").Object().Value("href").String().Match("/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35$")
	response.Value("links").Object().Value("job").Object().Value("id").Equal("042e216a-7822-4fa0-a3d6-e3f5248ffc35")
	response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("/datasets/" + datasetID + "$")
	response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
	response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/" + version + "/dimensions$")
	response.Value("links").Object().Value("edition").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "$")
	response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
	response.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetID + "/editions/" + edition + "/versions/" + version + "$")
	response.Value("links").Object().Value("version").Object().Value("id").Equal(version)
	response.Value("links").Object().Value("self").Object().Value("href").String().Match("/instances/" + instanceID + "$")
	response.Value("links").Object().Value("spatial").Object().Value("href").Equal("http://ons.gov.uk/geographylist")
	response.Value("release_date").Equal("2017-12-12")
	response.Value("temporal").Array().Element(0).Object().Value("start_date").Equal("2014-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("end_date").Equal("2017-09-09")
	response.Value("temporal").Array().Element(0).Object().Value("frequency").Equal("monthly")
	response.Value("import_tasks").Object().Value("import_observations").Object().Value("total_inserted_observations").Number().Equal(1000)
	response.Value("total_observations").Number().Equal(1000)
	response.Value("version").Number().Equal(versionNumber)
}
