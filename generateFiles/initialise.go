package generateFiles

import (
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	codeListDatabase = "codelists"
	datasetDatabase  = "datasets"
	filterDatabase   = "filters"
	importDatabase   = "imports"
	importTracker    = "dp-import-tracker"

	region     = "eu-west-1"
	bucketName = "ons-dp-cmd-test"

	datasetName = "cpih01"

	internalTokenHeader    = "Internal-Token"
	internalTokenID        = "FD0108EA-825D-411C-9B1D-41EF7727F465"
	invalidInternalTokenID = "FD0108EA-825D-411C-9B1D-41EF7727F465A"
)

func init() {
	var err error

	cfg, err = config.Get()
	if err != nil {
		log.ErrorC("Unable to access configurations", err, nil)
		os.Exit(1)
	}

	if err = mongo.NewDatastore(cfg.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		os.Exit(1)
	}

	// Remove test data that is left in mongo from previous test run
	if success := deleteMongoTestData(datasetName); !success {
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}

func deleteMongoTestData(datasetID string) bool {
	successfullyRemovedMongoTestData := true

	oldInstanceResource, err := mongo.GetInstance("datasets", "instances", "links.dataset.id", datasetID)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove test resources", err, log.Data{"links.dataset.id": datasetID})
			return false
		}

		log.Info("instance does not exist, nothing to delete carry on", log.Data{"links.dataset.id": datasetID})
		return successfullyRemovedMongoTestData
	}

	// delete dataset
	if err = mongo.Teardown("datasets", "datasets", "_id", datasetID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove dataset resource", err, log.Data{"dataset_id": datasetName})
			successfullyRemovedMongoTestData = false
		}
		log.Trace("delete dataset not found", nil)
	}

	instanceID := oldInstanceResource.InstanceID
	log.Info("Removing test data associated to instance", log.Data{"instance_id": instanceID})

	// remove job
	if err = mongo.Teardown("imports", "imports", "links.instances.id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove job resource", err, log.Data{"links.instances[0].id": instanceID})
			successfullyRemovedMongoTestData = false
		}
		log.Trace("delete job not found", nil)
	}

	// remove instance/versions
	if err = mongo.Teardown("datasets", "instances", "id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove instance resource", err, log.Data{"instance_id": instanceID})
			successfullyRemovedMongoTestData = false
		}
		log.Trace("delete instance not found", nil)
	}

	// remove dimension options
	if err = mongo.Teardown("datasets", "dimension.options", "instance_id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove dimension option resources", err, log.Data{"instance_id": instanceID})
			successfullyRemovedMongoTestData = false
		}
		log.Trace("delete dimension options not found", nil)
	}

	// remove edition if exists
	if oldInstanceResource.Links.Edition != nil {
		if err = mongo.Teardown("datasets", "editions", "links.self.href", oldInstanceResource.Links.Edition.HRef); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("failed to remove edition resource", err, log.Data{"links.self.href": oldInstanceResource.Links.Edition.HRef})
				successfullyRemovedMongoTestData = false
			}
			log.Trace("delete edition not found", nil)
		}
	}

	// remove filter blueprint
	if err = mongo.Teardown("filters", "filters", "instance_id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove filter blueprint resource", err, log.Data{"instance_id": instanceID})
			successfullyRemovedMongoTestData = false
		}
	}

	// remove filter output
	if err = mongo.Teardown("filters", "filterOutputs", "instance_id", instanceID); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove filter output resource", err, log.Data{"instance_id": instanceID})
			successfullyRemovedMongoTestData = false
		}
	}

	log.Info("removed mongo test data", log.Data{"success": successfullyRemovedMongoTestData})

	return successfullyRemovedMongoTestData
}
