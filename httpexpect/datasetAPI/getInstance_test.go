package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetInstance(t *testing.T) {

	if err := mongo.Teardown(database, "instances", "_id", "799"); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.Setup(database, "instances", "_id", "799", validUnpublishedInstanceData); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Get an instance", t, func() {
		Convey("When user is authenticated", func() {
			response := datasetAPI.GET("/instances/{id}", "799").WithHeader(internalToken, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			response.Value("id").Equal("799")
			response.Value("collection_id").Equal("208064B3-A808-449B-9041-EA3A2F72CFAB")
			response.Value("downloads").Object().Value("csv").Object().Value("url").String().Match("(.+)/aws/census-2017-2-csv$")
			response.Value("downloads").Object().Value("csv").Object().Value("size").Equal("10mb")
			response.Value("downloads").Object().Value("xls").Object().Value("url").String().Match("(.+)/aws/census-2017-2-xls$")
			response.Value("downloads").Object().Value("xls").Object().Value("size").Equal("24mb")
			response.Value("edition").Equal(edition)
			response.Value("headers").Array().Element(0).String().Equal("time")
			response.Value("headers").Array().Element(1).String().Equal("geography")
			response.Value("links").Object().Value("job").Object().Value("href").String().Match("(.+)/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35$")
			response.Value("links").Object().Value("job").Object().Value("id").Equal("042e216a-7822-4fa0-a3d6-e3f5248ffc35")
			response.Value("links").Object().Value("dataset").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "$")
			response.Value("links").Object().Value("dataset").Object().Value("id").Equal(datasetID)
			response.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2/dimensions$")
			response.Value("links").Object().Value("edition").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "$")
			response.Value("links").Object().Value("edition").Object().Value("id").Equal(edition)
			response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/" + datasetID + "/editions/" + edition + "/versions/2$")
			response.Value("links").Object().Value("version").Object().Value("id").Equal("2")
			response.Value("links").Object().Value("self").Object().Value("href").String().Match("(.+)/instances/799$")
			response.Value("release_date").Equal("2017-12-12")
			response.Value("state").Equal("associated")
			response.Value("total_inserted_observations").Equal(1000)
			response.Value("total_observations").Equal(1000)
			response.Value("version").Equal(2)
		})
	})

	if err := mongo.Teardown(database, "instances", "_id", "799"); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

func TestFailureToGetInstance(t *testing.T) {

	if err := mongo.Teardown(database, "instances", "_id", "799"); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	Convey("Fail to get instance document", t, func() {
		Convey("and return status not found", func() {
			Convey("when instance document does not exist", func() {
				datasetAPI.GET("/instances/{id}", "799").WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusNotFound)
			})
		})
		Convey("and return status unauthorised", func() {
			Convey("when an unauthorised user sends a GET request", func() {
				if err := mongo.Setup(database, "instances", "_id", "799", validUnpublishedInstanceData); err != nil {
					log.ErrorC("Was unable to run test", err, nil)
					os.Exit(1)
				}

				datasetAPI.GET("/instances/{id}", "799").
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

	if err := mongo.Teardown(database, "instances", "_id", "799"); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}
