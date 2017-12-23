package filterAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func setupDimensionOptions(id string, update bson.M) *mongo.Doc {
	return &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "_id",
		Value:      id,
		Update:     update,
	}
}

func setupMultipleDimensionsAndOptions(instanceID string) []*mongo.Doc {
	var docs []*mongo.Doc

	options := []bson.M{
		GetValidAgeDimensionData(instanceID, "27"),
		GetValidAgeDimensionData(instanceID, "28"),
		GetValidSexDimensionData(instanceID, "male"),
		GetValidSexDimensionData(instanceID, "female"),
		GetValidSexDimensionData(instanceID, "unknown"),
		GetValidGoodsAndServicesDimensionData(instanceID, "Education"),
		GetValidGoodsAndServicesDimensionData(instanceID, "health"),
		GetValidGoodsAndServicesDimensionData(instanceID, "communication"),
		GetValidGoodsAndServicesDimensionData(instanceID, "welfare"),
		GetValidTimeDimensionData(instanceID, "March 1997"),
		GetValidTimeDimensionData(instanceID, "April 1997"),
		GetValidTimeDimensionData(instanceID, "June 1997"),
		GetValidTimeDimensionData(instanceID, "September 1997"),
		GetValidTimeDimensionData(instanceID, "December 1997"),
		GetValidTimeDimensionData(instanceID, "February 2007"),
		GetValidResidenceTypeDimensionData(instanceID, "Lives in a communal establishment"),
		GetValidResidenceTypeDimensionData(instanceID, "Lives in a household"),
	}

	for _, o := range options {
		docs = append(docs, setupDimensionOptions(uuid.NewV4().String(), o))
	}

	return docs
}
