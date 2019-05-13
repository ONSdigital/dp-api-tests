package codeListAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	collection = "test"

	firstCodeListID           = "cpih1dim1aggid"
	firstCodeListEdition      = "one-off"
	firstCodeListFirstCodeID  = "cpih1dim1S90401"
	firstCodeListFirstLabel   = "09.4.1 Recreational and sporting services"
	firstCodeListSecondCodeID = "cpih1dim1S90501"
	firstCodeListSecondLabel  = "09.5.1 Books"
	firstCodeListThirdCodeID  = "cpih1dim1S90402"
	firstCodeListThirdLabel   = "09.4.2 Cultural services"

	secondCodeListID          = "uk-only"

	invalidCodeListID = "1C3221283FD544F0BBAD619779D8960E"
	firstCode         = "cpih1dim1S90401"
	secondCode        = "cpih1dim1S90501"
	thirdCode         = "cpih1dim1S90402"
	invalidCode       = "AC!@£$)98"
)

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.ErrorC("Unable to access configurations", err, nil)
		os.Exit(1)
	}

	store, err := neo4j.NewDatastore(cfg.Neo4jAddr, "", neo4j.GenericHierarchyCPIHTestData)
	if err != nil {
		log.ErrorC("neo4j datastore error", err, nil)
		os.Exit(1)
	}

	err = store.CreateCPIHCodeList()
	if err != nil {
		log.ErrorC("neo4j datastore error", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
