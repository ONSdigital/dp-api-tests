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

var relatedDatasets = GenericObject{
	HRef:  "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices",
	Title: "Consumer Price Inflation time series dataset",
}

var validPublishedDatasetData = bson.M{
	"$set": bson.M{
		"id": datasetID,
		"current.collection_id":             "108064B3-A808-449B-9041-EA3A2F72CFAA",
		"current.contacts":                  []Contact{contact},
		"current.description":               "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
		"current.id":                        datasetID,
		"current.keywords":                  []string{"cpi", "boy"},
		"current.last_updated":              "2017-06-06", // TODO this should be an isodate
		"current.links.editions.href":       "http://localhost:8080/datasets/" + datasetID + "/editions",
		"current.links.latest_version.id":   "1",
		"current.links.latest_version.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/1",
		"current.links.self.href":           "http://localhost:8080/datasets/" + datasetID,
		"current.methodologies":             []GenericObject{methodology},
		"current.national_statistic":        true,
		"current.next_release":              "2017-10-10",
		"current.publications":              []GenericObject{publication},
		"current.publisher.name":            "Automation Tester",
		"current.publisher.type":            "publisher",
		"current.publisher.href":            "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		"current.qmi.description":           "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
		"current.qmi.href":                  "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
		"current.qmi.title":                 "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
		"current.related_datasets":          []GenericObject{relatedDatasets},
		"current.release_frequency":         "Monthly",
		"current.state":                     "published",
		"current.theme":                     "Goods and services",
		"current.title":                     "CPI",
		"current.uri":                       "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
		"next.collection_id":                "208064B3-A808-449B-9041-EA3A2F72CFAB",
		"next.contacts":                     []Contact{contact},
		"next.description":                  "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
		"next.id":                           datasetID,
		"next.keywords":                     []string{"cpi", "boy"},
		"next.last_updated":                 "2017-10-11", // TODO this should be an isodate
		"next.links.editions.href":          "http://localhost:8080/datasets/" + datasetID + "/editions",
		"next.links.latest_version.id":      "1",
		"next.links.latest_version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2018/versions/1",
		"next.links.self.href":              "http://localhost:8080/datasets/" + datasetID,
		"next.methodologies":                []GenericObject{methodology},
		"next.national_statistic":           true,
		"next.next_release":                 "2018-10-10",
		"next.publications":                 []GenericObject{publication},
		"next.publisher.name":               "Automation Tester",
		"next.publisher.type":               "publisher",
		"next.publisher.href":               "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		"next.qmi.description":              "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
		"next.qmi.href":                     "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
		"next.qmi.title":                    "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
		"next.related_datasets":             []GenericObject{relatedDatasets},
		"next.release_frequency":            "Monthly",
		"next.state":                        "created",
		"next.theme":                        "Goods and services",
		"next.title":                        "CPI",
		"next.uri":                          "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
		"test_data":                         "true",
	},
}

var validUnpublishedDatasetData = bson.M{
	"$set": bson.M{
		"id":                             "133",
		"next.collection_id":             "208064B3-A808-449B-9041-EA3A2F72CFAB",
		"next.contacts":                  []Contact{contact},
		"next.description":               "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
		"next.id":                        "133",
		"next.keywords":                  []string{"cpi", "boy"},
		"next.last_updated":              "2017-10-11", // TODO this should be an isodate
		"next.links.editions.href":       "http://localhost:8080/datasets/" + datasetID + "/editions",
		"next.links.latest_version.id":   "1",
		"next.links.latest_version.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2018/versions/1",
		"next.links.self.href":           "http://localhost:8080/datasets/" + datasetID,
		"next.methodologies":             []GenericObject{methodology},
		"next.national_statistic":        true,
		"next.next_release":              "2018-10-10",
		"next.publications":              []GenericObject{publication},
		"next.publisher.name":            "Automation Tester",
		"next.publisher.type":            "publisher",
		"next.publisher.href":            "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
		"next.qmi.description":           "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
		"next.qmi.href":                  "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
		"next.qmi.title":                 "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
		"next.related_datasets":          []GenericObject{relatedDatasets},
		"next.release_frequency":         "Monthly",
		"next.state":                     "created",
		"next.theme":                     "Goods and services",
		"next.title":                     "CPI",
		"next.uri":                       "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
		"test_data":                      "true",
	},
}

var validTimeDimensionsData = bson.M{
	"$set": bson.M{

		"_id":                  "9811",
		"instance_id":          instanceID,
		"name":                 "time",
		"value":                "202.45",
		"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
		"links.code_list.href": "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
		"links.code.id":        "202.45",
		"links.code.href":      "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45",

		"node_id": "",

		"last_updated": "2017-09-09", // TODO Should be isodate
		"test_data":    "true",
	},
}
var validAggregateDimensionsData = bson.M{
	"$set": bson.M{

		"_id":                  "9812",
		"instance_id":          instanceID,
		"name":                 "aggregate",
		"value":                "cpi1dimA19",
		"label":                "CPI (Overall Index)",
		"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
		"links.code_list.href": "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
		"links.code.id":        "cpi1dimA19",
		"links.code.href":      "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/cpi1dimA19",

		"last_updated": "2017-09-08", // TODO Should be isodate
		"test_data":    "true",
	},
}

// var validTimeDimensionsOptionsData = bson.M{
// 	"$set": bson.M{

// 		"_id":         dimensionOptionID,
// 		"instance_id": instanceID,
// 		"name":        "time",
// 		"value":       "2050.56",

// 		"label":         "",
// 		"links.code.id": "2050.56",

// 		"links.code.href":      "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/2050.56",
// 		"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
// 		"links.code_list.href": "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
// 		"node_id":              "90",
// 		"test_data":            "true",
// 	},
// }

// var validSexDimensionsOptionsData = bson.M{
// 	"$set": bson.M{

// 		"_id":         dimensionOptionID,
// 		"instance_id": instanceID,
// 		"name":        "sex",
// 		"value":       "male",

// 		"label":         "male",
// 		"links.code.id": "2050.56",

// 		"links.code.href":      "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/2050.56",
// 		"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
// 		"links.code_list.href": "http://localhost:8080/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
// 		"node_id":              "92",
// 		"test_data":            "true",
// 	},
// }
var validPublishedEditionData = bson.M{
	"$set": bson.M{
		"edition":             "2017",
		"id":                  editionID,
		"last_updated":        "2017-09-08", // TODO Should be isodate
		"links.dataset.id":    datasetID,
		"links.dataset.href":  "http://localhost:8080/datasets/" + datasetID,
		"links.self.href":     "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.versions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions",
		"state":               "published",
		"test_data":           "true",
	},
}

var validUnpublishedEditionData = bson.M{
	"$set": bson.M{
		"edition":             "2018",
		"id":                  "466",
		"last_updated":        "2017-10-08", // TODO Should be isodate
		"links.dataset.id":    datasetID,
		"links.dataset.href":  "http://localhost:8080/datasets/" + datasetID,
		"links.self.href":     "http://localhost:8080/datasets/" + datasetID + "/editions/2018",
		"links.versions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2018/versions",
		"state":               "edition-confirmed",
		"test_data":           "true",
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
		"links.edition.id":      edition,
		"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.self.href":       "http://localhost:8080/instances/" + instanceID,
		"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/1",
		"links.version.id":      "1",
		"release_date":          "2017-12-12", // TODO Should be isodate
		"state":                 "published",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     1,
		"test_data":                   "true",
	},
}

var validUnpublishedInstanceData = bson.M{
	"$set": bson.M{
		"collection_id":         "208064B3-A808-449B-9041-EA3A2F72CFAB",
		"downloads.csv.url":     "http://localhost:8080/aws/census-2017-2-csv",
		"downloads.csv.size":    "10mb",
		"downloads.xls.url":     "http://localhost:8080/aws/census-2017-2-xls",
		"downloads.xls.size":    "24mb",
		"edition":               edition,
		"headers":               []string{"time", "geography"},
		"id":                    "799",
		"last_updated":          "2017-09-08", // TODO Should be isodate
		"license":               "ONS license",
		"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.dataset.id":      datasetID,
		"links.dataset.href":    "http://localhost:8080/datasets/" + datasetID,
		"links.dimensions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
		"links.edition.id":      edition,
		"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.self.href":       "http://localhost:8080/instances/" + "799",
		"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2",
		"links.version.id":      "2",
		"release_date":          "2017-12-12", // TODO Should be isodate
		"state":                 "associated",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     2,
		"test_data":                   "true",
	},
}

var validCompletedInstanceData = bson.M{
	"$set": bson.M{
		"collection_id":         "208064B3-A808-449B-9041-EA3A2F72CFAB",
		"downloads.csv.url":     "http://localhost:8080/aws/census-2017-2-csv",
		"downloads.csv.size":    "10mb",
		"downloads.xls.url":     "http://localhost:8080/aws/census-2017-2-xls",
		"downloads.xls.size":    "24mb",
		"edition":               edition,
		"headers":               []string{"time", "geography"},
		"id":                    "799",
		"last_updated":          "2017-09-08", // TODO Should be isodate
		"license":               "ONS license",
		"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.dataset.id":      datasetID,
		"links.dataset.href":    "http://localhost:8080/datasets/" + datasetID,
		"links.dimensions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
		"links.edition.id":      edition,
		"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.self.href":       "http://localhost:8080/instances/" + "799",
		"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2",
		"links.version.id":      "2",
		"release_date":          "2017-12-12", // TODO Should be isodate
		"state":                 "completed",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     2,
		"test_data":                   "true",
	},
}

var validEditionConfirmedInstanceData = bson.M{
	"$set": bson.M{
		"collection_id":         "208064B3-A808-449B-9041-EA3A2F72CFAB",
		"downloads.csv.url":     "http://localhost:8080/aws/census-2017-2-csv",
		"downloads.csv.size":    "10mb",
		"downloads.xls.url":     "http://localhost:8080/aws/census-2017-2-xls",
		"downloads.xls.size":    "24mb",
		"edition":               edition,
		"headers":               []string{"time", "geography"},
		"id":                    "779",
		"last_updated":          "2017-09-08", // TODO Should be isodate
		"license":               "ONS license",
		"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.dataset.id":      datasetID,
		"links.dataset.href":    "http://localhost:8080/datasets/" + datasetID,
		"links.dimensions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
		"links.edition.id":      edition,
		"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017",
		"links.self.href":       "http://localhost:8080/instances/" + "799",
		"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/2017/versions/2",
		"links.version.id":      "2",
		"release_date":          "2017-12-12", // TODO Should be isodate
		"state":                 "edition-confirmed",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     2,
		"test_data":                   "true",
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

var validPUTUpdateDatasetJSON string = `{

		"collection_id": "308064B3-A808-449B-9041-EA3A2F72CFAC",
		"contacts": [
			{
			"email": "rpi@onstest.gov.uk",
			"name": "Test Automation",
			"telephone": "+44 (0)1833 456123"
			}
		],
		"description": "Producer Price Indices (PPIs) are a series of economic indicators that measure the price movement of goods bought and sold by UK manufacturers",
		"keywords": [
			"rpi"
		],
		"methodologies": [
			{
			"description": "The Producer Price Index (PPI) is a monthly survey that measures the price changes of goods bought and sold by UK manufacturers",
			"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/producerpriceindicesqmi",
			"title": "Producer price indices QMI"
			}
		],
		"national_statistic": false,
		"next_release": "18 September 2017",
		"publications": [
			{
			"description": "Changes in the prices of goods bought and sold by UK manufacturers including price indices of materials and fuels purchased (input prices) and factory gate prices (output prices)",
			"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/producerpriceinflation/september2017",
			"title": "Producer price inflation, UK: September 2017"
			}
		],
		"publisher": {
			"name": "Test Automation Engineer",
			"type": "publisher",
			"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/producerpriceinflation/september2017"
		},
		"qmi": {
			"description": "PPI provides an important measure of inflation",
			"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/producerpriceindicesqmi",
			"title": "The Producer Price Index (PPI) is a monthly survey that measures the price changes"
		},
		"related_datasets": [
			{
			"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex",
			"title": "Producer Price Index time series dataset"
			}
		],
		"release_frequency": "Quaterly",
		"state": "created",
		"theme": "Price movement of goods",
		"title": "RPI",
		"uri": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex"
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

var validPUTCompletedInstanceJSON string = `
{
  "state": "completed"
}`

var validPUTUpdateVersionJSON string = `
{
		"state": "published"
}`
