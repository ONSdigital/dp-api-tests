package codeListAPI

import (
	"os"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
)

var cfg *config.Config

const (
	database                  = "codelists"
	collection                = "codelists"
	firstCodeListID           = "1C322128-3FD5-44F0-BBAD-619779D8960E"
	firstCodeListFirstCodeID  = "45251AEA-B4DD-409C-8C0E-CD5867399843"
	firstCodeListSecondCodeID = "0A78FEAC-E5D5-48B3-BA91-270468B432D1"
	firstCodeListThirdCodeID  = "4A335104-8C52-49C6-BA68-828AE16F9EBB"
	secondCodeListID          = "C5FA175A-7EA0-4B39-B252-7B52BE75C9DE"
	thirdCodelistID           = "5A561370-9AB5-48A4-A619-BEC996DD0BDA"
	invalidCodeListID         = "1C3221283FD544F0BBAD619779D8960E"
	firstCode                 = "LS_00998877"
	secondCode                = "LS_00998811"
	thirdCode                 = "LS_00998822"
	invalidCode               = "AC!@£$)98"
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

	test := &mongo.Doc{
		Database:   database,
		Collection: collection,
		Key:        "test_data",
		Value:      "true",
	}

	if err = mongo.Teardown(test); err != nil {
		log.ErrorC("Unable to remove all test data from mongo db", err, nil)
		os.Exit(1)
	}

	log.Debug("config is:", log.Data{"config": cfg})
}
