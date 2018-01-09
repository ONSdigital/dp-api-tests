package filterAPI

import (
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/go-ns/log"
	uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func setupInstance(instanceID string, update bson.M) error {

	if err := teardownInstance(instanceID); err != nil {
		return err
	}

	return mongo.Setup("datasets", "instances", "instance_id", instanceID, update)
}

func setupDimensionOptions(ID string, update bson.M) error {
	return mongo.Setup("datasets", "dimension.options", "_id", ID, update)
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

func teardownDimensionOptions(instanceID string) error {
	if err := mongo.Teardown("datasets", "dimension.options", "instance_id", instanceID); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}

func setupMultipleDimensionsAndOptions(instanceID string) error {
	var err error

	// setup age dimension options
	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidAgeDimensionData(instanceID, "27")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "age", "option": "27"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidAgeDimensionData(instanceID, "28")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "age", "option": "28"})
	}

	// setup sex dimension options
	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidSexDimensionData(instanceID, "male")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "sex", "option": "male"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidSexDimensionData(instanceID, "female")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "sex", "option": "female"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidSexDimensionData(instanceID, "unknown")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "sex", "option": "unknown"})
	}

	// setup Goods and services dimension options
	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidGoodsAndServicesDimensionData(instanceID, "Education")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "Goods and services", "option": "welfare"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidGoodsAndServicesDimensionData(instanceID, "health")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "Goods and services", "option": "welfare"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidGoodsAndServicesDimensionData(instanceID, "communication")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "Goods and services", "option": "welfare"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidGoodsAndServicesDimensionData(instanceID, "welfare")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "Goods and services", "option": "welfare"})
	}

	// setup time dimension options
	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidTimeDimensionData(instanceID, "March 1997")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "time", "option": "February 2007"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidTimeDimensionData(instanceID, "April 1997")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "time", "option": "February 2007"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidTimeDimensionData(instanceID, "June 1997")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "time", "option": "February 2007"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidTimeDimensionData(instanceID, "September 1997")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "time", "option": "February 2007"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidTimeDimensionData(instanceID, "December 1997")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "time", "option": "February 2007"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidTimeDimensionData(instanceID, "February 2007")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "time", "option": "February 2007"})
	}

	// setup residence type dimension options
	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidResidenceTypeDimensionData(instanceID, "Lives in a communal establishment")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "Residence Type", "option": "Lives in a communal establishment"})
	}

	if err = setupDimensionOptions(uuid.NewV4().String(), GetValidResidenceTypeDimensionData(instanceID, "Lives in a household")); err != nil {
		log.ErrorC("Unable to setup dimension option", err, log.Data{"instance_id": instanceID, "dimension": "Residence Type", "option": "Lives in a communal establishment"})
	}

	if err != nil {
		return err
	}

	return nil
}
