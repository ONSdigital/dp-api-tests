package datasetAPI

import "github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"

const (
	created   = "created"
	submitted = "submitted"
	invalid   = "invalid"
)

func setupInstances(datasetID, edition string, instances map[string]string) ([]*mongo.Doc, error) {
	var docs []*mongo.Doc

	for instanceType, instanceID := range instances {

		switch instanceType {
		case created:
			createdInstanceDoc := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validCreatedInstanceData(datasetID, edition, instanceID, "created"),
			}

			docs = append(docs, createdInstanceDoc)
		case submitted:
			submittedInstanceDoc := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validSubmittedInstanceData(datasetID, edition, instanceID),
			}

			docs = append(docs, submittedInstanceDoc)
		case invalid:
			invalidInstanceDoc := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "_id",
				Value:      instanceID,
				Update:     validCreatedInstanceData(datasetID, edition, instanceID, "gobbledygook"),
			}

			docs = append(docs, invalidInstanceDoc)
		}
	}

	if err := mongo.Setup(docs...); err != nil {
		return nil, err
	}

	return docs, nil
}
