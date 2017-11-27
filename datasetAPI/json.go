package datasetAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"gopkg.in/mgo.v2/bson"
)

var alert = mongo.Alert{
	Date:        "2017-12-10",
	Description: "A correction to an observation for males of age 25, previously 11 now changed to 12",
	Type:        "Correction",
}

var contact = mongo.ContactDetails{
	Email:     "cpi@onstest.gov.uk",
	Name:      "Automation Tester",
	Telephone: "+44 (0)1633 123456",
}

var latestChanges = mongo.LatestChange{
	Description: "The border of Southampton changed after the south east cliff face fell into the sea.",
	Name:        "Changes in Classification",
	Type:        "Summary of Changes",
}

var methodology = mongo.GeneralDetails{
	Description: "Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.",
	HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
	Title:       "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
}

var publication = mongo.GeneralDetails{
	Description: "Price indices, percentage changes and weights for the different measures of consumer price inflation.",
	HRef:        "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
	Title:       "UK consumer price inflation: August 2017",
}

var relatedDatasets = mongo.GeneralDetails{
	HRef:  "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices",
	Title: "Consumer Price Inflation time series dataset",
}

var dimension = mongo.CodeList{
	Description: "A list of ages between 18 and 75+",
	HRef:        "http://localhost:8080/codelists/408064B3-A808-449B-9041-EA3A2F72CFAC",
	ID:          "408064B3-A808-449B-9041-EA3A2F72CFAC",
	Name:        "age",
}

var temporal = mongo.TemporalFrequency{
	EndDate:   "2017-09-09",
	Frequency: "monthly",
	StartDate: "2014-09-09",
}

func validPublishedDatasetData(datasetID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"id": datasetID,
			"current.collection_id":             "108064B3-A808-449B-9041-EA3A2F72CFAA",
			"current.contacts":                  []mongo.ContactDetails{contact},
			"current.description":               "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
			"current.id":                        datasetID,
			"current.keywords":                  []string{"cpi", "boy"},
			"current.last_updated":              "2017-06-06", // TODO this should be an isodate
			"current.license":                   "ONS license",
			"current.links.access_rights.href":  "http://ons.gov.uk/accessrights",
			"current.links.editions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			"current.links.latest_version.id":   "1",
			"current.links.latest_version.href": cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/1",
			"current.links.self.href":           cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"current.methodologies":             []mongo.GeneralDetails{methodology},
			"current.national_statistic":        true,
			"current.next_release":              "2017-10-10",
			"current.publications":              []mongo.GeneralDetails{publication},
			"current.publisher.name":            "Automation Tester",
			"current.publisher.type":            "publisher",
			"current.publisher.href":            "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
			"current.qmi.description":           "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
			"current.qmi.href":                  "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
			"current.qmi.title":                 "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
			"current.related_datasets":          []mongo.GeneralDetails{relatedDatasets},
			"current.release_frequency":         "Monthly",
			"current.state":                     "published",
			"current.theme":                     "Goods and services",
			"current.title":                     "CPI",
			"current.unit_of_measure":           "Pounds Sterling",
			"current.uri":                       "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
			"next.collection_id":                "208064B3-A808-449B-9041-EA3A2F72CFAB",
			"next.contacts":                     []mongo.ContactDetails{contact},
			"next.description":                  "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
			"next.id":                           datasetID,
			"next.keywords":                     []string{"cpi", "boy"},
			"next.last_updated":                 "2017-10-11", // TODO this should be an isodate
			"next.license":                      "ONS license",
			"next.links.access_rights.href":     "http://ons.gov.uk/accessrights",
			"next.links.editions.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			"next.links.latest_version.id":      "1",
			"next.links.latest_version.href":    cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2018/versions/1",
			"next.links.self.href":              cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"next.methodologies":                []mongo.GeneralDetails{methodology},
			"next.national_statistic":           true,
			"next.next_release":                 "2018-10-10",
			"next.publications":                 []mongo.GeneralDetails{publication},
			"next.publisher.name":               "Automation Tester",
			"next.publisher.type":               "publisher",
			"next.publisher.href":               "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
			"next.qmi.description":              "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
			"next.qmi.href":                     "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
			"next.qmi.title":                    "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
			"next.related_datasets":             []mongo.GeneralDetails{relatedDatasets},
			"next.release_frequency":            "Monthly",
			"next.state":                        "created",
			"next.theme":                        "Goods and services",
			"next.title":                        "CPI",
			"next.unit_of_measure":              "Pounds Sterling",
			"next.uri":                          "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
			"test_data":                         "true",
		},
	}
}

func validAssociatedDatasetData(datasetID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"id":                             datasetID,
			"next.collection_id":             "208064B3-A808-449B-9041-EA3A2F72CFAB",
			"next.contacts":                  []mongo.ContactDetails{contact},
			"next.description":               "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
			"next.id":                        datasetID,
			"next.keywords":                  []string{"cpi", "boy"},
			"next.last_updated":              "2017-10-11", // TODO this should be an isodate
			"next.links.access_rights.href":  "http://ons.gov.uk/accessrights",
			"next.links.editions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			"next.links.latest_version.id":   "1",
			"next.links.latest_version.href": cfg.DatasetAPIURL + "/datasets" + datasetID + "/editions/2018/versions/1",
			"next.links.self.href":           cfg.DatasetAPIURL + "/datasets" + datasetID,
			"next.methodologies":             []mongo.GeneralDetails{methodology},
			"next.national_statistic":        true,
			"next.next_release":              "2018-10-10",
			"next.publications":              []mongo.GeneralDetails{publication},
			"next.publisher.name":            "Automation Tester",
			"next.publisher.type":            "publisher",
			"next.publisher.href":            "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
			"next.qmi.description":           "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
			"next.qmi.href":                  "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
			"next.qmi.title":                 "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
			"next.related_datasets":          []mongo.GeneralDetails{relatedDatasets},
			"next.release_frequency":         "Monthly",
			"next.state":                     "associated",
			"next.theme":                     "Goods and services",
			"next.title":                     "CPI",
			"next.unit_of_measure":           "Pounds Sterling",
			"next.uri":                       "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
			"test_data":                      "true",
		},
	}
}

func validCreatedDatasetData(datasetID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"id":                             datasetID,
			"next.contacts":                  []mongo.ContactDetails{contact},
			"next.description":               "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
			"next.id":                        datasetID,
			"next.keywords":                  []string{"cpi", "boy"},
			"next.last_updated":              "2017-10-11", // TODO this should be an isodate
			"next.links.access_rights.href":  "http://ons.gov.uk/accessrights",
			"next.links.editions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions",
			"next.links.latest_version.id":   "1",
			"next.links.latest_version.href": cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2018/versions/1",
			"next.links.self.href":           cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"next.methodologies":             []mongo.GeneralDetails{methodology},
			"next.national_statistic":        true,
			"next.next_release":              "2018-10-10",
			"next.publications":              []mongo.GeneralDetails{publication},
			"next.publisher.name":            "Automation Tester",
			"next.publisher.type":            "publisher",
			"next.publisher.href":            "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
			"next.qmi.description":           "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
			"next.qmi.href":                  "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
			"next.qmi.title":                 "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)",
			"next.related_datasets":          []mongo.GeneralDetails{relatedDatasets},
			"next.release_frequency":         "Monthly",
			"next.state":                     "created",
			"next.theme":                     "Goods and services",
			"next.title":                     "CPI",
			"next.unit_of_measure":           "Pounds Sterling",
			"next.uri":                       "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation",
			"test_data":                      "true",
		},
	}
}

func validTimeDimensionsData(instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                  "9811",
			"instance_id":          instanceID,
			"name":                 "time",
			"option":               "202.45",
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code.id":        "202.45",
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45",
			"node_id":              "",
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func validTimeDimensionsDataWithOutOptions(instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                  "9811",
			"instance_id":          instanceID,
			"name":                 "time",
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code.id":        "202.45",
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/202.45",
			"node_id":              "",
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func validAggregateDimensionsData(instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                  "9812",
			"instance_id":          instanceID,
			"name":                 "aggregate",
			"option":               "cpi1dimA19",
			"label":                "CPI (Overall Index)",
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code.id":        "cpi1dimA19",
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/cpi1dimA19",
			"last_updated":         "2017-09-08", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func validPublishedEditionData(datasetID, editionID, edition string) bson.M {
	return bson.M{
		"$set": bson.M{
			"edition":                   edition,
			"id":                        editionID,
			"last_updated":              "2017-09-08", // TODO Should be isodate
			"links.dataset.id":          datasetID,
			"links.dataset.href":        cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.latest_version.id":   "1",
			"links.latest_version.href": cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1",
			"links.self.href":           cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition,
			"links.versions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions",
			"state":                     "published",
			"test_data":                 "true",
		},
	}
}

func validUnpublishedEditionData(datasetID, editionID, edition string) bson.M {
	return bson.M{
		"$set": bson.M{
			"edition":             edition,
			"id":                  editionID,
			"last_updated":        "2017-10-08", // TODO Should be isodate
			"links.dataset.id":    datasetID,
			"links.dataset.href":  cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.self.href":     cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition,
			"links.versions.href": cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions",
			"state":               "edition-confirmed",
			"test_data":           "true",
		},
	}
}

func validPublishedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"alerts":                      []mongo.Alert{alert},
			"collection_id":               "108064B3-A808-449B-9041-EA3A2F72CFAA",
			"dimensions":                  []mongo.CodeList{dimension},
			"downloads.csv.url":           cfg.DatasetAPIURL + "/aws/census-2017-1-csv",
			"downloads.csv.size":          "10mb",
			"downloads.xls.url":           cfg.DatasetAPIURL + "/aws/census-2017-1-xls",
			"downloads.xls.size":          "24mb",
			"edition":                     edition,
			"headers":                     []string{"time", "geography"},
			"id":                          instanceID,
			"latest_changes":              []mongo.LatestChange{latestChanges},
			"last_updated":                "2017-09-08", // TODO Should be isodate
			"license":                     "ONS License",
			"links.job.id":                "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":              cfg.DatasetAPIURL + "/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":            datasetID,
			"links.dataset.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.dimensions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions",
			"links.edition.id":            edition,
			"links.edition.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition,
			"links.self.href":             cfg.DatasetAPIURL + "/instances/" + instanceID,
			"links.spatial.href":          "http://ons.gov.uk/geographylist",
			"links.version.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1",
			"links.version.id":            "2",
			"release_date":                "2017-12-12", // TODO Should be isodate
			"state":                       "published",
			"temporal":                    []mongo.TemporalFrequency{temporal},
			"total_inserted_observations": 1000,
			"total_observations":          1000,
			"version":                     1,
			"test_data":                   "true",
		},
	}
}

func validAssociatedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"collection_id":               "208064B3-A808-449B-9041-EA3A2F72CFAB",
			"dimensions":                  []mongo.CodeList{dimension},
			"downloads.csv.url":           cfg.DatasetAPIURL + "/aws/census-2017-2-csv",
			"downloads.csv.size":          "10mb",
			"downloads.xls.url":           cfg.DatasetAPIURL + "/aws/census-2017-2-xls",
			"downloads.xls.size":          "24mb",
			"edition":                     edition,
			"headers":                     []string{"time", "geography"},
			"id":                          instanceID,
			"last_updated":                "2017-09-08", // TODO Should be isodate
			"latest_changes":              []mongo.LatestChange{latestChanges},
			"license":                     "ONS license",
			"links.job.id":                "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":              cfg.DatasetAPIURL + "/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":            datasetID,
			"links.dataset.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.dimensions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/2/dimensions",
			"links.edition.id":            edition,
			"links.edition.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition,
			"links.self.href":             cfg.DatasetAPIURL + "/instances/" + instanceID,
			"links.spatial.href":          "http://ons.gov.uk/geographylist",
			"links.version.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/2",
			"links.version.id":            "2",
			"release_date":                "2017-12-12", // TODO Should be isodate
			"state":                       "associated",
			"temporal":                    []mongo.TemporalFrequency{temporal},
			"total_inserted_observations": 1000,
			"total_observations":          1000,
			"version":                     2,
			"test_data":                   "true",
		},
	}
}

func validEditionConfirmedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"dimensions":                  []mongo.CodeList{dimension},
			"downloads.csv.url":           cfg.DatasetAPIURL + "/aws/census-2017-2-csv",
			"downloads.csv.size":          "10mb",
			"downloads.xls.url":           cfg.DatasetAPIURL + "/aws/census-2017-2-xls",
			"downloads.xls.size":          "24mb",
			"edition":                     edition,
			"headers":                     []string{"time", "geography"},
			"id":                          instanceID,
			"last_updated":                "2017-09-08", // TODO Should be isodate
			"license":                     "ONS license",
			"links.job.id":                "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":              cfg.DatasetAPIURL + "/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":            datasetID,
			"links.dataset.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.dimensions.href":       cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/2/dimensions",
			"links.edition.id":            edition,
			"links.edition.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition,
			"links.self.href":             cfg.DatasetAPIURL + "/instances/" + instanceID,
			"links.spatial.href":          "http://ons.gov.uk/geographylist",
			"links.version.href":          cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/2",
			"links.version.id":            "2",
			"release_date":                "2017-12-12", // TODO Should be isodate
			"state":                       "edition-confirmed",
			"temporal":                    []mongo.TemporalFrequency{temporal},
			"total_inserted_observations": 1000,
			"total_observations":          1000,
			"version":                     2,
			"test_data":                   "true",
		},
	}
}

func validCompletedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"collection_id":         "208064B3-A808-449B-9041-EA3A2F72CFAB",
			"downloads.csv.url":     cfg.DatasetAPIURL + "/aws/census-2017-2-csv",
			"downloads.csv.size":    "10mb",
			"downloads.xls.url":     cfg.DatasetAPIURL + "/aws/census-2017-2-xls",
			"downloads.xls.size":    "24mb",
			"edition":               edition,
			"headers":               []string{"time", "geography"},
			"id":                    instanceID,
			"last_updated":          "2017-09-08", // TODO Should be isodate
			"latest_changes":        []mongo.LatestChange{latestChanges},
			"license":               "ONS license",
			"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":        cfg.DatasetAPIURL + "/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":      datasetID,
			"links.dataset.href":    cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.dimensions.href": cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/2/dimensions",
			"links.edition.id":      edition,
			"links.edition.href":    cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017",
			"links.self.href":       cfg.DatasetAPIURL + "/instances/" + instanceID,
			"links.version.href":    cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/2017/versions/2",
			"links.version.id":      "2",
			"release_date":          "2017-12-12", // TODO Should be isodate
			"state":                 "completed",
			"total_inserted_observations": 1000,
			"total_observations":          1000,
			"version":                     2,
			"test_data":                   "true",
		},
	}
}

func validCreatedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"edition":            edition,
			"headers":            []string{"time", "geography"},
			"id":                 instanceID,
			"last_updated":       "2017-09-08", // TODO Should be isodate
			"links.job.id":       "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":     cfg.DatasetAPIURL + "/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":   datasetID,
			"links.dataset.href": cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.self.href":    cfg.DatasetAPIURL + "/instances/" + instanceID,
			"state":              "created",
			"total_observations": 1000,
			"test_data":          "true",
		},
	}
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
	"license": "ONS license",
	"links": {
		"access_rights": {
			"href": "http://ons.gov.uk/accessrights"
		}
	},
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
	"unit_of_measure": "Pounds Sterling",
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
		"release_frequency": "Quarterly",
		"state": "associated",
		"theme": "Price movement of goods",
		"title": "RPI",
		"unit_of_measure": "Pounds",
		"uri": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/producerpriceindex"
}`

var validPOSTCreateInstanceJSON string = `
{
  "links": {
    "job": {
      "id": "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
      "href": "http://localhost:21800/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35"
    }
  }
}`

var validPOSTCreateFullInstanceJSON string = `
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
	"dimensions": [
	  {
			"description": "The age ranging from 16 to 75+",
			"href": "http://localhost:22400//code-lists/43513D18-B4D8-4227-9820-492B2971E7T5",
			"id": "43513D18-B4D8-4227-9820-492B2971E7T5",
			"name": "age"
	  }
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

var validPUTFullInstanceJSON string = `
{
	"alerts": [
	  {
		  "date": "2017-04-05",
		  "description": "All data entries (observations) for Plymouth have been updated",
			"type": "Correction"
	  }
	],
	"dimensions": [
		{
			"description": "The age ranging from 16 to 75+",
			"href": "http://localhost:22400//code-lists/43513D18-B4D8-4227-9820-492B2971E7T5",
			"id": "43513D18-B4D8-4227-9820-492B2971E7T5",
			"name": "age"
		}
	],
	"latest_changes": [
	  {
		  "description": "change to the period frequency from quarterly to monthly",
			"name": "Changes to the period frequency",
			"type": "Summary of Changes"
	  }
	],
	"links": {
		"spatial": {
			"href": "http://ons.gov.uk/geography-list"
		}
	},
	"release_date": "2017-11-11",
  "state": "completed",
	"temporal": [
		{
			"start_date": "2014-10-10",
			"end_date": "2016-10-10",
			"frequency": "monthly"
		}
	],
	"total_inserted_observations": 1000
}`

var validPUTEditionConfirmedInstanceJSON string = `
{
  "alerts": [
	  {
		  "date": "2017-04-05",
		  "description": "All data entries (observations) for Plymouth have been updated",
		  "type": "Correction"
	  }
  ],
	"dimensions": [
		{
			"description": "The age ranging from 16 to 75+",
			"href": "http://localhost:22400//code-lists/43513D18-B4D8-4227-9820-492B2971E7T5",
			"id": "43513D18-B4D8-4227-9820-492B2971E7T5",
			"name": "age"
		}
	],
	"latest_changes": [
	  {
		  "description": "change to the period frequency from quarterly to monthly",
			"name": "Changes to the period frequency",
			"type": "Summary of Changes"
	  }
	],
	"links": {
		"spatial": {
			"href": "http://ons.gov.uk/geography-list"
		}
	},
	"release_date": "2017-11-11",
  "state": "edition-confirmed",
	"temporal": [
		{
			"start_date": "2014-10-10",
			"end_date": "2016-10-10",
			"frequency": "monthly"
		}
	],
	"total_inserted_observations": 1000
}`

var validPUTUpdateVersionMetaDataJSON string = `
{
"alerts": [
	{
		"date": "2017-04-05",
		"description": "All data entries (observations) for Plymouth have been updated",
		"type": "Correction"
	}
],
"latest_changes": [
	{
		"description": "change to the period frequency from quarterly to monthly",
		"name": "Changes to the period frequency",
		"type": "Summary of Changes"
	}
],
"links": {
  "spatial": {
	  "href": "http://ons.gov.uk/new-geography-list"
	},
  "self": {
	  "href": "http://bogus/bad-link"
	}
},
"release_date": "2018-11-11",
"temporal": [
	{
		"start_date": "2014-11-11",
		"end_date": "2017-11-11",
		"frequency": "monthly"
	}
]
}`

var validPUTUpdateVersionAlertsJSON string = `
{
"alerts": [
	{
		"date": "2017-04-05",
		"description": "All data entries (observations) for Plymouth have been updated",
		"type": "Correction"
	}
],
}`

var validPUTUpdateVersionToAssociatedJSON string = `
{
	"state": "associated",
	"collection_id": "45454545"
}`

var validPUTUpdateVersionFromAssociatedToEditionConfirmedJSON string = `
{
	"collection_id": ""
}`

var validPUTUpdateVersionToPublishedWithCollectionIDJSON string = `
{
	"collection_id": "33333333",
	"state": "published"
}`

var validPUTUpdateVersionToPublishedJSON string = `
{
	"state": "published"
}`

var invalidPOSTCreateInstanceJSON string = `
{
  "links": {
    "dataset": {
    	"id": "34B13D18-B4D8-4227-9820-492B2971E221",
      "href": "http://localhost:21800/datasets/34B13D18-B4D8-4227-9820-492B2971E221"
    }
  }
}`
