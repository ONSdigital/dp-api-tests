package datasetAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/globalsign/mgo/bson"
)

const (
	created          = "created"
	submitted        = "submitted"
	completed        = "completed"
	editionConfirmed = "edition-confirmed"
	invalid          = "invalid"
)

func setupInstances(datasetID, edition string, uniqueTimestamp bson.MongoTimestamp, instances map[string]string) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	for instanceType, instanceID := range instances {

		switch instanceType {
		case created:
			createdInstanceDoc := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validCreatedInstanceData(datasetID, edition, instanceID, "created", uniqueTimestamp),
			}

			docs = append(docs, createdInstanceDoc)
		case submitted:
			submittedInstanceDoc := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validSubmittedInstanceData(datasetID, edition, instanceID, "submitted", uniqueTimestamp),
			}

			docs = append(docs, submittedInstanceDoc)
		case invalid:
			invalidInstanceDoc := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validCreatedInstanceData(datasetID, edition, instanceID, "gobbledygook", uniqueTimestamp),
			}

			docs = append(docs, invalidInstanceDoc)
		}
	}

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
