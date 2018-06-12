package downloadService

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gedge/mgo/bson"
)

func validPublishedDataset(datasetID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":           datasetID,
			"current.state": "published",
		},
	}
}

func validPublishedEdition(datasetID, edition string) bson.M {
	return bson.M{
		"$set": bson.M{
			"current.edition":          edition,
			"current.links.dataset.id": datasetID,
			"current.state":            "published",
			"next.edition":             edition,
			"next.links.dataset.id":    datasetID,
			"next.state":               "published",
		},
	}
}

func validPublishedVersionWithPublicLink(datasetID, edition string, version int) bson.M {
	return bson.M{
		"$set": bson.M{
			"downloads":          mongo.DownloadList{CSV: &mongo.DownloadObject{Public: publicLink}},
			"edition":            edition,
			"links.dataset.id":   datasetID,
			"links.version.href": "version-link",
			"links.self.href":    "self-link",
			"version":            version,
			"state":              "published",
		},
	}
}

func validVersionWithPrivateLink(datasetID, edition string, version int, state string) bson.M {
	return bson.M{
		"$set": bson.M{
			"downloads":          mongo.DownloadList{CSV: &mongo.DownloadObject{Private: privateLink}},
			"edition":            edition,
			"links.dataset.id":   datasetID,
			"links.version.href": "version-link",
			"links.self.href":    "self-link",
			"version":            version,
			"state":              state,
		},
	}
}

func validPublishedInstanceData(datasetID, edition, instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"collection_id":                                                "",
			"downloads.csv.url":                                            cfg.DatasetAPIURL + "/aws/census-2017-1-csv",
			"downloads.csv.size":                                           "10",
			"downloads.csv.public":                                         "https://s3-eu-west-1.amazon.com/csv-exported/myfile.csv",
			"downloads.csv.private":                                        "s3://csv-exported/myfile.csv",
			"downloads.xls.url":                                            cfg.DatasetAPIURL + "/aws/census-2017-1-xls",
			"downloads.xls.size":                                           "24",
			"downloads.xls.public":                                         "https://s3-eu-west-1.amazon.com/csv-exported/myfile.xls",
			"downloads.xls.private":                                        "s3://csv-exported/myfile.xls",
			"edition":                                                      edition,
			"headers":                                                      []string{"time", "geography"},
			"id":                                                           instanceID,
			"last_updated":                                                 "2017-09-08", // TODO Should be isodate
			"license":                                                      "ONS License",
			"links.job.id":                                                 "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":                                               cfg.DatasetAPIURL + "/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":                                             datasetID,
			"links.dataset.href":                                           cfg.DatasetAPIURL + "/datasets/" + datasetID,
			"links.dimensions.href":                                        cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1/dimensions",
			"links.edition.id":                                             edition,
			"links.edition.href":                                           cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition,
			"links.self.href":                                              cfg.DatasetAPIURL + "/instances/" + instanceID,
			"links.spatial.href":                                           "http://ons.gov.uk/geographylist",
			"links.version.href":                                           cfg.DatasetAPIURL + "/datasets/" + datasetID + "/editions/" + edition + "/versions/1",
			"links.version.id":                                             "1",
			"release_date":                                                 "2017-12-12", // TODO Should be isodate
			"state":                                                        "published",
			"total_inserted_observations":                                  1000,
			"total_observations":                                           1000,
			"version":                                                      1,
			"test_data":                                                    "true",
			"import_tasks.import_observations.state":                       "completed",
			"import_tasks.import_observations.total_inserted_observations": 1000,
		},
	}
}
