package filterAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func setupInstance(instanceID string, update bson.M) error {

	if err := teardownInstance(instanceID); err != nil {
		return err
	}

	if err := mongo.Setup("datasets", "instances", "instance_id", instanceID, update); err != nil {
		return err
	}

	return nil
}

func teardownInstance(instanceID string) error {
	if err := mongo.Teardown("datasets", "instances", "instance_id", instanceID); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}
