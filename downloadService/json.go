package downloadService

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"gopkg.in/mgo.v2/bson"
)

func validPublishedDataset(datasetID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":   datasetID,
			"state": "published",
		},
	}
}

func validPublishedEdition(datasetID, edition string) bson.M {
	return bson.M{
		"$set": bson.M{
			"edition":          edition,
			"links.dataset.id": datasetID,
			"state":            "published",
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
