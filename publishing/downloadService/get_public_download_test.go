package downloadService

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestRedirectToPublicDownload(t *testing.T) {
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	versionID := 1

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDataset(datasetID),
	}

	edition := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "editions",
		Key:        "_id",
		Value:      editionID,
		Update:     validPublishedEdition(datasetID, editionID),
	}

	version := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "_id",
		Value:      strconv.Itoa(versionID),
		Update:     validPublishedVersionWithPublicLink(datasetID, editionID, versionID),
	}

	if err := mongo.Setup(dataset, edition, version); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	downloadService := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewRequireReporter(t),
		BaseURL:  cfg.DownloadServiceURL,
		Client:   cli,
	})

	Convey("Given a published version exists with a public link", t, func() {
		Convey("When a request is made for the public document", func() {
			Convey("Then the response redirects to the the public file", func() {

				response := downloadService.GET("/downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv", datasetID, editionID, versionID).
					WithHeader(authHeader, serviceToken).
					Expect().Status(http.StatusMovedPermanently)

				response.Header("Location").Equal(publicLink)
			})
		})
	})

	Convey("Given a version does not exist", t, func() {
		Convey("When a request is made to the download service for that resource", func() {
			Convey("Then the download service returns a not found http status code", func() {
				downloadService.GET("/downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv", datasetID, editionID, 2).
					WithHeader(authHeader, serviceToken).
					Expect().Status(http.StatusNotFound).
					Body().Contains("resource not found")
			})
		})
	})

	if err := mongo.Teardown(dataset, edition, version); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

}
