package filterAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	mgo "gopkg.in/mgo.v2"
)

func setupInstance() error {

	if err := teardownInstance(); err != nil {
		return err
	}

	if err := mongo.Setup("datasets", "instances", "instance_id", instanceID, ValidPublishedInstanceData); err != nil {
		return err
	}

	return nil
}

func teardownInstance() error {
	if err := mongo.Teardown("datasets", "instances", "instance_id", instanceID); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}
