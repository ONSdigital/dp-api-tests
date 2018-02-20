package generateFiles

import (
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	importTracker = "dp-import-tracker"

	region     = "eu-west-1"
	bucketName = "ons-dp-cmd-test"

	datasetName      = "cpih01"
	genericHierarchy = "cpih1dim1aggid"

	internalTokenHeader    = "Internal-Token"
	internalTokenID        = "FD0108EA-825D-411C-9B1D-41EF7727F465"
	invalidInternalTokenID = "FD0108EA-825D-411C-9B1D-41EF7727F465A"
	importAPIInternalToken = "0C30662F-6CF6-43B0-A96A-954772267FF5"
)

var dropDatabases = []string{"test"}

func init() {
	var err error

	cfg, err = config.Get()
	if err != nil {
		log.ErrorC("Unable to access configurations", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})

	if err = mongo.NewDatastore(cfg.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		os.Exit(1)
	}

	if err = mongo.DropDatabases(dropDatabases); err != nil {
		log.ErrorC("failed to drop mongo databases", err, log.Data{"databases": dropDatabases})
	}

	// Remove test data that is left in mongo from previous test run
	if success := deleteMongoTestData(datasetName); !success {
		os.Exit(1)
	}

	// Create generic hierarchy for test (CPIH)
	if err = generateGenericHierarchy(); err != nil {
		log.ErrorC("neo4j datastore error", err, nil)
		os.Exit(1)
	}
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

	instanceID := oldInstanceResource.InstanceID

	var docs []*mongo.Doc

	dataset := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "datasets",
		Key:        "_id",
		Value:      datasetID,
	}

	importJob := &mongo.Doc{
		Database:   cfg.MongoImportsDB,
		Collection: "imports",
		Key:        "links.instances.id",
		Value:      instanceID,
	}

	instance := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "instances",
		Key:        "id",
		Value:      instanceID,
	}

	dimension := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: "dimension.options",
		Key:        "instance_id",
		Value:      instanceID,
	}

	if oldInstanceResource.Links.Edition != nil {
		edition := &mongo.Doc{
			Database:   cfg.MongoDB,
			Collection: "editions",
			Key:        "links.self.href",
			Value:      oldInstanceResource.Links.Edition.HRef,
		}
		docs = append(docs, edition)
	}

	filterBlueprint := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filters",
		Key:        "instance_id",
		Value:      instanceID,
	}

	filterOutput := &mongo.Doc{
		Database:   cfg.MongoFiltersDB,
		Collection: "filterOutputs",
		Key:        "instance_id",
		Value:      instanceID,
	}
	docs = append(docs, dataset, importJob, instance, dimension, filterBlueprint, filterOutput)

	// remove filter output
	if err = mongo.Teardown(docs...); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to remove previous test resource", err, log.Data{"instance_id": instanceID})
			successfullyRemovedMongoTestData = false
		}
	}

	log.Info("removed mongo test data", log.Data{"success": successfullyRemovedMongoTestData})

	return successfullyRemovedMongoTestData
}

func generateGenericHierarchy() error {
	datastore, err := neo4j.NewDatastore(cfg.Neo4jAddr, "", neo4j.GenericHierarchyCPIHTestData)
	if err != nil {
		log.ErrorC("unable to connect to neo4j", err, nil)
		return err
	}

	if err = datastore.CreateGenericHierarchy(genericHierarchy); err != nil {
		log.ErrorC("unable to create generic hierarchy", err, log.Data{"generic_hierarchy": genericHierarchy})
		return err
	}

	return nil
}