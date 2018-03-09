package downloadService

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
)

func TestPrivateDownloadDecryptedAndStreamed(t *testing.T) {
	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Info("skipping private download tests, as no vault token or address set - use make test", nil)
		t.Skip()
	}

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
		Update:     validVersionWithPrivateLink(datasetID, editionID, versionID, "published"),
	}

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := sendV4FileToAWS(region, bucketName, fileName); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := mongo.Setup(dataset, edition, version); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given a published version exists with a private link", t, func() {
		Convey("When a request is made for the private document", func() {
			Convey("Then the response streams the decrypted private file", func() {

				response := downloadService.GET("/downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv", datasetID, editionID, versionID).
					Expect().Status(http.StatusOK)

				response.Body().Equal(string(f))
			})
		})
	})

	if err := mongo.Teardown(dataset, edition, version); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

	if err := deleteS3File(region, bucketName, fileName); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

}

func TestPrivateDownloadDecryptedAndStreamedWithAuthentication(t *testing.T) {
	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Info("skipping private download tests, as no vault token or address set - use make test", nil)
		t.Skip()
	}

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
		Update:     validVersionWithPrivateLink(datasetID, editionID, versionID, "associated"),
	}

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := sendV4FileToAWS(region, bucketName, fileName); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := mongo.Setup(dataset, edition, version); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given an associated version exists with a private link", t, func() {
		Convey("When a request is made for the private document without authentication", func() {
			Convey("Then the response returns a not found http status", func() {

				response := downloadService.GET("/downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv", datasetID, editionID, versionID).
					Expect().Status(http.StatusNotFound)

				response.Body().Contains("resource not found")
			})
		})
	})

	Convey("Given an associated version exists with a private link", t, func() {
		Convey("When a request is made for the private document with authentication", func() {
			Convey("Then the response streams the decrypted private file", func() {

				response := downloadService.GET("/downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv", datasetID, editionID, versionID).WithHeader(internalToken, internalTokenID).
					Expect().Status(http.StatusOK)

				response.Body().Equal(string(f))
			})
		})
	})

	if err := mongo.Teardown(dataset, edition, version); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}

	if err := deleteS3File(region, bucketName, fileName); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}
}

func TestPrivateDownloadDecryptedAndStreamedFailure(t *testing.T) {

	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Info("skipping private download tests, as no vault token or address set - use make test", nil)
		t.Skip()
	}

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
		Update:     validVersionWithPrivateLink(datasetID, editionID, versionID, "published"),
	}

	if err := mongo.Setup(dataset, edition, version); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given a public version exists with a private link, but the file is missing from Amazon S3", t, func() {
		Convey("When a request is made for the private document", func() {
			Convey("Then the download service returns an internal server error status code", func() {
				response := downloadService.GET("/downloads/datasets/{datasetID}/editions/{edition}/versions/{version}.csv", datasetID, editionID, versionID).
					Expect().Status(http.StatusInternalServerError)

				response.Body().Contains("internal server error")
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
