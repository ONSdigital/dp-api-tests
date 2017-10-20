package filterAPI

import "github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"

func setupDatastores() {

	mongo.RemoveAll(database, collection)
	mongo.Teardown("datasets", "instances", "instance_id", instanceID)
	mongo.Setup("datasets", "instances", "instance_id", instanceID, validPublishedInstanceData)
}
