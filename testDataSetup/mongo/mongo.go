package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/go-ns/log"
)

var session *mgo.Session

// ManyDocs represents a list of objects that are able to query mongo db
type ManyDocs struct {
	Docs []Doc
}

// Doc contains information to be able to query mongo db
type Doc struct {
	Database   string
	Collection string
	ID         string
	Update     bson.M
}

// NewDatastore creates a new mgo.Session with a strong consistency and a write mode of "majority"
func newDatastore(uri string) error {
	if session == nil {
		var err error
		if session, err = mgo.Dial(uri); err != nil {
			return err
		}

		session.EnsureSafe(&mgo.Safe{WMode: "majority"})
		session.SetMode(mgo.Strong, true)
	}
	return nil
}

// Teardown is a way of cleaning up an individual document from mongo instance
func Teardown(database, collection, id string) error {
	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		return err
	}

	if err = newDatastore(config.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		return err
	}

	s := session.Copy()
	defer s.Clone()

	if err = s.DB(database).C(collection).Remove(bson.M{"_id": id}); err != nil {
		if err == mgo.ErrNotFound {
			log.Info("data does not exist, continue", nil)
			return nil
		}
		return err
	}

	return nil
}

// TeardownMany is a way of cleaning up many documents from mongo instance
func TeardownMany(d ManyDocs) error {
	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		return err
	}

	if err = newDatastore(config.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		return err
	}

	s := session.Copy()
	defer s.Clone()

	for _, doc := range d.Docs {
		if err = s.DB(doc.Database).C(doc.Collection).Remove(bson.M{"_id": doc.ID}); err != nil {
			return err
		}
	}

	return nil
}

// Setup is a way of loading in an individual document into a mongo instance
func Setup(database, collection, id string, update bson.M) error {
	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		return err
	}

	if err = newDatastore(config.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		return err
	}

	s := session.Copy()
	defer s.Clone()

	if err = s.DB(database).C(collection).Update(bson.M{"_id": id}, update); err != nil {
		return err
	}

	log.Info("SetUp completed", nil)
	return nil
}

// SetupMany is a way of loading in many documents into a mongo instance
func SetupMany(d ManyDocs) error {
	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		return err
	}

	if err = newDatastore(config.MongoAddr); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		return err
	}

	s := session.Copy()
	defer s.Clone()

	for _, doc := range d.Docs {
		if err = s.DB(doc.Database).C(doc.Collection).Update(bson.M{"_id": doc.ID}, doc.Update); err != nil {
			return err
		}
	}

	log.Info("SetUp completed", nil)
	return nil
}
