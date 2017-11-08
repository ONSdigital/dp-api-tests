package importAPI

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

// Get a list of all jobs
// Lists can be filtered by the job state
// 200 - A list of jobs has been returned
func TestSuccessfullyGetListOfImportJobs(t *testing.T) {

	var docs []mongo.Doc

	importCreateJobDoc := mongo.Doc{
		Database:   "imports",
		Collection: "imports",
		Key:        "_id",
		Value:      jobID,
		Update:     validCreatedImportJobData,
	}

	importSubmittedJobDoc := mongo.Doc{
		Database:   "imports",
		Collection: "imports",
		Key:        "_id",
		Value:      "01C24F0D-24BE-479F-962B-C76BCCD0AD00",
		Update:     validSubmittedImportJobData,
	}

	docs = append(docs, importCreateJobDoc, importSubmittedJobDoc)

	d := &mongo.ManyDocs{
		Docs: docs,
	}

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}
	}

	if err := mongo.SetupMany(d); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

	Convey("Given an existing job", t, func() {

		Convey("Get a list of all jobs", func() {

			response := importAPI.GET("/jobs").Expect().Status(http.StatusOK).JSON().Array()

			fmt.Println(response)
			checkImportJobsResponse(response)

		})

	})

	if err := mongo.TeardownMany(d); err != nil {
		if err != mgo.ErrNotFound {
			os.Exit(1)
		}
	}
}

func checkImportJobsResponse(response *httpexpect.Array) {

	response.Element(0).Object().Value("id").Equal(jobID)
	response.Element(0).Object().Value("recipe").Equal("2080CACA-1A82-411E-AA46-F00804968E78")
	response.Element(0).Object().Value("state").Equal("Created")

}

// func TestSuccessfullyGetListOfImportJobs(t *testing.T) {

// 	importAPI := httpexpect.New(t, cfg.ImportAPIURL)

// 	Convey("Given an existing job", t, func() {

// 		importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
// 			Expect().Status(http.StatusCreated)

// 		Convey("Get a list of all jobs", func() {

// 			importAPI.GET("/jobs").Expect().Status(http.StatusOK).JSON().NotNull()

// 		})

// 	})
// }
