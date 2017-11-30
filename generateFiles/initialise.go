package generateFiles

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	codeListDatabase = "codelists"
	datasetDatabase  = "datasets"
	filterDatabase   = "filters"
	importDatabase   = "imports"

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

	databases := []string{codeListDatabase, datasetDatabase, filterDatabase, importDatabase}
	if err = mongo.DropDatabases(databases); err != nil {
		log.ErrorC("Unable to remove all data from mongo db", err, nil)
		os.Exit(1)
	}

	if err = neo4j.DropDatabases(cfg.Neo4jAddr); err != nil {
		log.ErrorC("Failed to remove all data from neo4j database", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
