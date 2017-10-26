package filterAPI

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

func TestSuccessfullyPostFilterJob(t *testing.T) {

	err := teardownCreateFilterTestData()
	if err != nil {
		log.ErrorC("Failed to tear down test data", err, nil)
		os.Exit(1)
	}

	if err = mongo.Setup("datasets", "instances", "_id", instanceID, ValidPublishedInstanceData); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given a valid json input to create a filter", t, func() {
		Convey("Then the response returns a status of created (201)", func() {

			response := filterAPI.POST("/filters").WithBytes([]byte(ValidPOSTCreateFilterJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()

			// TODO Check all fields in response
			response.Value("filter_job_id").NotNull()
			response.Value("dimension_list_url").String().Match("(.+)/filters/(.+)/dimensions$")
			response.Value("instance_id").Equal(instanceID)
			response.Value("links").Object().Value("version").Object().Value("href").String().Match("(.+)/datasets/123/editions/2017/versions/1$")
			response.Value("links").Object().Value("version").Object().Value("id").Equal("1")
			response.Value("state").Equal("created")
		})
	})

	err = teardownCreateFilterTestData()
	if err != nil {
		log.ErrorC("Failed to tear down test data", err, nil)
		os.Exit(1)
	}
}

func TestFailureToPostFilterJob(t *testing.T) {

	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	Convey("Given invalid json input to create a filter", t, func() {
		Convey("Then the response returns status bad request (400)", func() {

			filterAPI.POST("/filters").WithBytes([]byte(InvalidJSON)).
				Expect().Status(http.StatusBadRequest).Body().Contains("Bad request - Invalid request body\n")
		})
	})

	err := teardownCreateFilterTestData()
	if err != nil {
		log.ErrorC("Failed to tear down test data", err, nil)
		os.Exit(1)
	}
}

func teardownCreateFilterTestData() error {
	var docs []mongo.Doc

	instanceDoc := mongo.Doc{
		Database:   "datasets",
		Collection: "instances",
		Key:        "_id",
		Value:      instanceID,
		Update:     ValidPublishedInstanceData,
	}

	filterJobDoc := mongo.Doc{
		Database:   database,
		Collection: collection,
		Key:        "instance_id",
		Value:      instanceID,
		Update:     ValidCreatedFilterJob,
	}

	docs = append(docs, instanceDoc, filterJobDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}

	return nil
}
