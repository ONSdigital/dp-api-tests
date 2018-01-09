package datasetAPI

import (
	"net/http"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulyUpdateDataset(t *testing.T) {

	datasetID := uuid.NewV4().String()

	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	Convey("Given a published dataset already exists", t, func() {

		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		originalDataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", datasetID)
		if err != nil {
			log.ErrorC("Unable to retrieve original dataset document", err, nil)
			os.Exit(1)
		}

		// Check dataset current subdocument
		expectedCurrentSubDoc := expectedCurrentSubDoc(datasetID, "2017")
		So(originalDataset.Current, ShouldResemble, expectedCurrentSubDoc)

		Convey("When a Put request is made to update the dataset", func() {
			Convey("Then the dataset resource is updated and response contains a status ok (200)", func() {
				datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTUpdateDatasetJSON)).
					Expect().Status(http.StatusOK)

				expectedNextSubDoc := expectedNextSubDoc(datasetID, "2018")

				dataset, err := mongo.GetDataset(cfg.MongoDB, collection, "_id", datasetID)
				if err != nil {
					log.ErrorC("Unable to retrieve updated dataset document", err, nil)
					os.Exit(1)
				}

				// Check dataset current subdocument has not changed
				So(dataset.Current, ShouldResemble, expectedCurrentSubDoc)

				// Check dataset next subdocument does not match the original dataset next subdocument
				So(originalDataset.Next, ShouldNotResemble, expectedNextSubDoc)

				So(dataset.Next, ShouldResemble, expectedNextSubDoc)
			})
		})

		if err := mongo.Teardown(publishedDataset); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func TestFailureToUpdateDataset(t *testing.T) {

	datasetID := uuid.NewV4().String()
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)

	publishedDataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
		Update:     validPublishedDatasetData(datasetID),
	}

	Convey("Given a published dataset does not exist", t, func() {
		Convey("When an authorised PUT request is made to update dataset resource", func() {
			Convey("Then fail to update resource and return a status of not found (404) with a message `Dataset not found`", func() {

				datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader(internalToken, internalTokenID).WithBytes([]byte(validPUTUpdateDatasetJSON)).
					Expect().Status(http.StatusNotFound).Body().Contains("Dataset not found\n")
			})
		})
	})

	Convey("Given a published dataset exists", t, func() {

		if err := mongo.Setup(publishedDataset); err != nil {
			log.ErrorC("Was unable to run test", err, nil)
			os.Exit(1)
		}

		Convey("When an unauthorised PUT request is made to update a dataset resource with an invalid authentication header", func() {
			Convey("Then fail to update resource and return a status unauthorized (401) with a message `Unauthorised access to API`", func() {

				datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader(internalToken, invalidInternalTokenID).WithBytes([]byte(validPUTUpdateDatasetJSON)).
					Expect().Status(http.StatusUnauthorized).Body().Contains("Unauthorised access to API\n")
			})
		})

		Convey("When no authentication header is provided in PUT request to update dataset resource", func() {
			Convey("Then fail to update resource and return a status of unauthorized (401) with a message `No authentication header provided`", func() {

				datasetAPI.POST("/datasets/{id}", datasetID).WithBytes([]byte(validPUTUpdateDatasetJSON)).
					Expect().Status(http.StatusUnauthorized).Body().Contains("No authentication header provided\n")
			})
		})

		Convey("When an authorised PUT request is made to update dataset resource with an invalid body", func() {
			Convey("Then fail to update resource and return a status of bad request (400) with a message `Failed to parse json body`", func() {

				datasetAPI.PUT("/datasets/{id}", datasetID).WithHeader(internalToken, internalTokenID).WithBytes([]byte("{")).
					Expect().Status(http.StatusBadRequest).Body().Contains("Failed to parse json body\n")
			})
		})

		if err := mongo.Teardown(publishedDataset); err != nil {
			if err != mgo.ErrNotFound {
				os.Exit(1)
			}
		}
	})
}

func expectedCurrentSubDoc(datasetID, edition string) *mongo.Dataset {
	contactDetails := mongo.ContactDetails{
		Email:     "cpi@onstest.gov.uk",
		Name:      "Automation Tester",
		Telephone: "+44 (0)1633 123456",
	}

	expectedCurrentMethodology := mongo.GeneralDetails{
		Description: "Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.",
		HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
		Title:       "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
	}

	nationalStatistic := true

	expectedCurrentPublication := mongo.GeneralDetails{
		Description: "Price indices, percentage changes and weights for the different measures of consumer price inflation.",
		HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		Title:       "UK consumer price inflation: August 2017",
	}

	relatedDataset := mongo.GeneralDetails{
		HRef:  "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices",
		Title: "Consumer Price Inflation time series dataset",
	}

	currentSubDoc := &mongo.Dataset{
		CollectionID: "108064B3-A808-449B-9041-EA3A2F72CFAA",
		Contacts:     []mongo.ContactDetails{contactDetails},
		Description:  "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
		Keywords:     []string{"cpi", "boy"},
		ID:           "",
		License:      "ONS license",
		Links: &mongo.DatasetLinks{
			AccessRights: &mongo.LinkObject{
				HRef: "http://ons.gov.uk/accessrights",
			},
			Editions: &mongo.LinkObject{
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			},
			LatestVersion: &mongo.LinkObject{
				ID:   "1",
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1",
			},
			Self: &mongo.LinkObject{
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID,
			},
		},
		Methodologies:     []mongo.GeneralDetails{expectedCurrentMethodology},
		NationalStatistic: &nationalStatistic,
		NextRelease:       "2017-10-10",
		Publications:      []mongo.GeneralDetails{expectedCurrentPublication},
		Publisher: &mongo.Publisher{
			Name: "Automation Tester",
			Type: "publisher",
			HRef: "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		},
		QMI: &mongo.GeneralDetails{
			Description: "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
			HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
			Title:       "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
		},
		RelatedDatasets:  []mongo.GeneralDetails{relatedDataset},
		ReleaseFrequency: "Monthly",
		State:            "published",
		Theme:            "Goods and services",
		Title:            "CPI",
		UnitOfMeasure:    "Pounds Sterling",
		URI:              "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
	}

	return currentSubDoc
}

func expectedNextSubDoc(datasetID, edition string) *mongo.Dataset {
	contactDetails := mongo.ContactDetails{
		Email:     "rpi@onstest.gov.uk",
		Name:      "Test Automation",
		Telephone: "+44 (0)1833 456123",
	}

	expectedMethodology := mongo.GeneralDetails{
		Description: "The Producer Price Index (PPI) is a monthly survey that measures the price changes of goods bought and sold by UK manufacturers",
		HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/producerpriceindicesqmi",
		Title:       "Producer price indices QMI",
	}

	nationalStatistic := false

	expectedPublication := mongo.GeneralDetails{
		Description: "Changes in the prices of goods bought and sold by UK manufacturers including price indices of materials and fuels purchased (input prices) and factory gate prices (output prices)",
		HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/producerpriceinflation/september2017",
		Title:       "Producer price inflation, UK: September 2017",
	}

	relatedDataset := mongo.GeneralDetails{
		HRef:  "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex",
		Title: "Producer Price Index time series dataset",
	}

	nextSubDoc := &mongo.Dataset{
		CollectionID: "308064B3-A808-449B-9041-EA3A2F72CFAC",
		Contacts:     []mongo.ContactDetails{contactDetails},
		Description:  "Producer Price Indices (PPIs) are a series of economic indicators that measure the price movement of goods bought and sold by UK manufacturers",
		Keywords:     []string{"rpi"},
		ID:           "",
		License:      "ONS license",
		Links: &mongo.DatasetLinks{
			AccessRights: &mongo.LinkObject{
				HRef: "http://ons.gov.uk/accessrights",
			},
			Editions: &mongo.LinkObject{
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			},
			LatestVersion: &mongo.LinkObject{
				ID:   "1",
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1",
			},
			Self: &mongo.LinkObject{
				HRef: cfg.DatasetAPIURL + "/datasets/" + datasetID,
			},
		},
		Methodologies:     []mongo.GeneralDetails{expectedMethodology},
		NationalStatistic: &nationalStatistic,
		NextRelease:       "18 September 2017",
		Publications:      []mongo.GeneralDetails{expectedPublication},
		Publisher: &mongo.Publisher{
			Name: "Test Automation Engineer",
			Type: "publisher",
			HRef: "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/producerpriceinflation/september2017",
		},
		QMI: &mongo.GeneralDetails{
			Description: "PPI provides an important measure of inflation",
			HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/producerpriceindicesqmi",
			Title:       "The Producer Price Index (PPI) is a monthly survey that measures the price changes",
		},
		RelatedDatasets:  []mongo.GeneralDetails{relatedDataset},
		ReleaseFrequency: "Quarterly",
		State:            "associated",
		Theme:            "Price movement of goods",
		Title:            "RPI",
		UnitOfMeasure:    "Pounds",
		URI:              "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex",
	}

	return nextSubDoc
}
