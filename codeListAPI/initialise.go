package codeListAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config
var ds *neo4j.Datastore

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	ds, err = neo4j.NewDatastore(cfg.Neo4jAddr, "", "")
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	ds.CodeListLabel = cfg.CodeListLabel
}
