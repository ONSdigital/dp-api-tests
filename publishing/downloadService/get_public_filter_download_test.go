package downloadService

import (
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

func TestRedirectToPublicFilterDownload(t *testing.T) {

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
		Update:     filterAPI.GetValidFilterWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, datasetID, edition, filterBlueprintID, version, publishedTrue),
	}

	filterDoc := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "_id",
		Value:      filterID,
		Update:     filterAPI.GetValidFilterOutputWithMultipleDimensionsBSON(cfg.FilterAPIURL, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, version, publishedTrue),
	}

	if err := mongo.Setup(filterBlueprintDoc, filterDoc); err != nil {
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

				response := downloadService.GET("/downloads/filter-outputs/{filterOutputID}.csv", filterOutputID).
					WithHeader(authHeader, serviceToken).
					Expect().Status(http.StatusMovedPermanently)

				response.Header("Location").Equal(publicLink)
			})
		})
	})

	Convey("Given a version does not exist", t, func() {
		Convey("When a request is made to the download service for that resource", func() {
			Convey("Then the download service returns a not found http status code", func() {

				downloadService.GET("/downloads/filter-outputs/12334534232_wut.csv", filterOutputID).
					WithHeader(authHeader, serviceToken).
					Expect().Status(http.StatusNotFound)
			})
		})
	})

	if err := mongo.Teardown(filterBlueprintDoc, filterDoc); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("Failed to tear down test data", err, nil)
			os.Exit(1)
		}
	}
}
