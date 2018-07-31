package searchAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/globalsign/mgo/bson"
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

var dimensionTwo = mongo.CodeList{
	Description: "An aggregate of the data",
	HRef:        "http://localhost:8080/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD",
	ID:          "508064B3-A808-449B-9041-EA3A2F72CFAD",
	Name:        "aggregate",
}

var dimensionThree = mongo.CodeList{
	Description: "The time in which this dataset spans",
	HRef:        "http://localhost:8080/codelists/508064B3-A808-449B-9041-EA3A2F72CFAD",
	ID:          "508064B3-A808-449B-9041-EA3A2F72CFAD",
	Name:        "time",
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

func validPublishedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"alerts":                      []mongo.Alert{alert},
			"collection_id":               "108064B3-A808-449B-9041-EA3A2F72CFAA",
			"dimensions":                  []mongo.CodeList{dimension, dimensionTwo, dimensionThree},
			"downloads.csv.url":           cfg.DatasetAPIURL + "/aws/census-2017-1-csv",
			"downloads.csv.size":          "10",
			"downloads.xls.url":           cfg.DatasetAPIURL + "/aws/census-2017-1-xls",
			"downloads.xls.size":          "24",
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
			"links.version.id":            "1",
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
