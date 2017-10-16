package datasetAPI

import "gopkg.in/mgo.v2/bson"

// Contact represents an object containing contact information
type Contact struct {
	Email     string
	Name      string
	Telephone string
}

// GenericObject represents a generic object structure
type GenericObject struct {
	Description string
	HRef        string
	Title       string
}

var contact = Contact{
	Email:     "cpi@onstest.gov.uk",
	Name:      "Automation Tester",
	Telephone: "+44 (0)1633 123456",
}

var methodology = GenericObject{
	Description: "Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.",
	HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
	Title:       "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
}

var publication = GenericObject{
	Description: "Price indices, percentage changes and weights for the different measures of consumer price inflation.",
	HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
	Title:       "UK consumer price inflation: August 2017",
}

var validDatasetData = bson.M{
	"$set": bson.M{
		"collection_id":             "108064B3-A808-449B-9041-EA3A2F72CFAA",
		"contacts":                  []Contact{contact},
		"description":               "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
		"id":                        datasetID,
		"keywords":                  []string{"cpi", "boy"},
		"last_updated":              "2017-06-06", // TODO this should be an isodate
		"links.editions.href":       "http://localhost:8080/datasets/" + datasetID + "/editions",
		"links.latest_version.id":   "1",
		"links.latest_version.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/1",
		"links.self.href":           "http://localhost:8080/datasets/" + datasetID,
		"methodologies":             []GenericObject{methodology},
		"national_statistic":        true,
		"next_release":              "17 October 2017",
		"publications":              []GenericObject{publication},
		"publisher.name":            "Automation Tester",
		"publisher.type":            "publisher",
		"publisher.href":            "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		"qmi.description":           "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
		"qmi.href":                  "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
		"qmi.title":                 "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
		"related_datasets.href":     "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices",
		"related_datasets.title":    "Consumer Price Inflation time series dataset",
		"release_frequency":         "Monthly",
		"state":                     "created",
		"theme":                     "Goods and services",
		"title":                     "CPI",
		"uri":                       "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
	},
}

var validPublishedEditionData = bson.M{
	"$set": bson.M{
		"edition":             "2017",
		"id":                  "208064B3-A808-449B-9041-EA3A2F72CFAB-2017",
		"last_updated":        "2017-09-08", // TODO Should be isodate
		"links.dataset.id":    datasetID,
		"links.self.href":     "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.versions.href": "http://localhost:8080/datasets/208064B3-A808-449B-9041-EA3A2F72CFAB/editions/2017/versions",
		"state":               "published",
	},
}

var validUnpublishedEditionData = bson.M{
	"$set": bson.M{
		"edition":             "2017",
		"id":                  "208064B3-A808-449B-9041-EA3A2F72CFAB-2017",
		"last_updated":        "2017-10-08", // TODO Should be isodate
		"links.dataset.id":    datasetID,
		"links.dataset.href":  "http://localhost:8080/datasets/" + datasetID,
		"links.self.href":     "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.versions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions",
		"state":               "edition-confirmed",
	},
}

var validPublishedInstanceData = bson.M{
	"$set": bson.M{
		"collection_id":         "108064B3-A808-449B-9041-EA3A2F72CFAA",
		"downloads.csv.url":     "http://localhost:8080/aws/census-2017-1-csv",
		"downloads.csv.size":    "10mb",
		"downloads.xls.url":     "http://localhost:8080/aws/census-2017-1-xls",
		"downloads.xls.size":    "24mb",
		"edition":               "2017",
		"headers":               []string{"time", "geography"},
		"id":                    instanceID,
		"last_updated":          "2017-09-08", // TODO Should be isodate
		"license":               "ONS License",
		"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.dataset.id":      datasetID,
		"links.dataset.href":    "http://localhost:8080/datasets/" + datasetID,
		"links.dimensions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/1/dimensions",
		"links.edition.id":      "2017",
		"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.self.href":       "http://localhost:8080/instances/" + instanceID,
		"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/1",
		"links.version.id":      "1",
		"release_date":          "2017-12-12", // TODO Should be isodate
		"state":                 "published",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     1,
	},
}

var validUnpublishedInstanceData = bson.M{
	"$set": bson.M{
		"collection_id":         "208064B3-A808-449B-9041-EA3A2F72CFAB",
		"downloads.csv.url":     "http://localhost:8080/aws/census-2017-1-csv",
		"downloads.csv.size":    "10mb",
		"downloads.xls.url":     "http://localhost:8080/aws/census-2017-1-xls",
		"downloads.xls.size":    "24mb",
		"edition":               "2017",
		"headers":               []string{"time", "geography"},
		"id":                    "799",
		"last_updated":          "2017-09-08", // TODO Should be isodate
		"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.dataset.id":      datasetID,
		"links.dataset.href":    "http://localhost:8080/datasets/" + datasetID,
		"links.dimensions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
		"links.edition.id":      "2017",
		"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.self.href":       "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2",
		"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2",
		"links.version.id":      "2",
		"state":                 "associated",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     2,
	},
}

var validPOSTCreateDatasetJSON string = `
{
	"collection_id": "108064B3-A808-449B-9041-EA3A2F72CFAA",
	"contacts": [
	  {
		"email": "cpi@onstest.gov.uk",
		"name": "Automation Tester",
		"telephone": "+44 (0)1633 123456"
	  }
	],
	"description": "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
	"keywords": [
	  "cpi"
	],
	"methodologies": [
	  {
		"description": "Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.",
		"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
		"title": "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)"
	  }
	],
	"national_statistic": true,
	"next_release": "17 October 2017",
	"publications": [
	  {
		"description": "Price indices, percentage changes and weights for the different measures of consumer price inflation.",
		"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		"title": "UK consumer price inflation: August 2017"
	  }
	],
	"publisher": {
	  "name": "Automation Tester",
	  "type": "publisher",
	  "href": "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017"
	},
	"qmi": {
	  "description": "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
	  "href": "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
	  "title": "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)"
	},
	"related_datasets": [
	  {
		"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices",
		"title": "Consumer Price Inflation time series dataset"
	  }
	],
	"release_frequency": "Monthly",
	"state": "created",
	"theme": "Goods and services",
	"title": "CPI",
	"uri": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation"
  }`

var validPOSTCreateInstanceJSON string = `
{
  "links": {
    "job": {
      "id": "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
      "href": "http://localhost:21800/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35"
		},
		"dataset": {
			"id": "34B13D18-B4D8-4227-9820-492B2971E221",
      "href": "http://localhost:21800/datasets/34B13D18-B4D8-4227-9820-492B2971E221"
		}
  },
  "state": "completed",
	"edition": "2017",
	"total_inserted_observations": 1000,
  "total_observations": 1000,
  "headers": [
		"time",
		"geography"
  ]
}`

var validPUTUpdateInstanceJSON string = `
{

  "state": "edition-confirmed"

}`

var validPUTUpdateVersionJSON string = `
{

		"state": "associated"
	}`
