package downloadService

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/web/filterAPI"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestPrivateFilterDownloadDecryptedAndStreamedWithoutError(t *testing.T) {
	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Info("failing test as no vault token or address set - use make test", nil)
		t.FailNow()
	}

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	editionID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	log.Debug("TestPrivateFilterDownloadDecryptedAndStreamed", log.Data{
		"filterID":          filterID,
		"filterBlueprintID": filterBlueprintID,
		"filterOutputID":    filterOutputID,
		"datasetID":         datasetID,
		"editionID":         editionID,
		"instanceID":        instanceID,
	})

	filterBlueprintDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filters",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, publishedTrue),
	}

	filterDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterOutputWithPrivateDownloads(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, version, publishedTrue),
	}

	origFileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := sendV4FileToAWS(region, bucketName, fileName); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := mongo.Setup(filterDoc, filterBlueprintDoc); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given a published version exists with a private link", t, func() {
		Convey("When a request is made for the private document", func() {
			Convey("Then the response streams the decrypted private file", func() {

				response := downloadService.GET("/downloads/filter-outputs/{filterOutputID}.csv", filterOutputID).
					Expect().Status(http.StatusOK)

				response.Body().Equal(string(origFileContent))
			})
		})
	})

	if err := mongo.Teardown(filterDoc, filterBlueprintDoc); err != nil {
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

func TestPrivateFilterDownloadDecryptedAndStreamedUnpublishedWithoutAuthentication(t *testing.T) {
	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Info("failing test as no vault token or address set - use make test", nil)
		t.FailNow()
	}

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterBlueprintDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filters",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, publishedFalse),
	}

	filterDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterOutputWithPrivateDownloads(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, version, publishedFalse),
	}

	if err := sendV4FileToAWS(region, bucketName, fileName); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	if err := mongo.Setup(filterDoc, filterBlueprintDoc); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given an associated version exists with a private link", t, func() {
		Convey("When a request is made for the private document without authentication", func() {
			Convey("Then the response returns a 404 http status", func() {

				downloadService.GET("/downloads/filter-outputs/{filterOutputID}.csv", filterOutputID).
					Expect().Status(http.StatusNotFound)

			})
		})
	})

	if err := mongo.Teardown(filterDoc, filterBlueprintDoc); err != nil {
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

func TestPrivateFilterDownloadDecryptedAndStreamedFailure(t *testing.T) {
	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Info("failing test as no vault token or address set - use make test", nil)
		t.FailNow()
	}

	filterID := uuid.NewV4().String()
	filterBlueprintID := uuid.NewV4().String()
	filterOutputID := uuid.NewV4().String()
	datasetID := uuid.NewV4().String()
	instanceID := uuid.NewV4().String()
	edition := "2017"
	version := 1

	filterBlueprintDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filters",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, publishedFalse),
	}

	filterDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterOutputWithPrivateDownloads(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, version, publishedFalse),
	}

	if err := mongo.Setup(filterDoc, filterBlueprintDoc); err != nil {
		log.ErrorC("Was unable to run test", err, nil)
		os.Exit(1)
	}

	downloadService := httpexpect.New(t, cfg.DownloadServiceURL)

	Convey("Given a public version exists with a private link, but the file is missing from Amazon S3", t, func() {
		Convey("When a request is made for the private document", func() {
			Convey("Then the download service returns a not found status code", func() {
				response := downloadService.GET("/downloads/filter-outputs/{filterOutputID}.csv", filterOutputID).
					Expect().Status(http.StatusNotFound)

				response.Body().Contains("resource not found")
			})
		})
	})

	if err := mongo.Teardown(filterDoc, filterBlueprintDoc); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}
