package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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
	Key        string
	Value      string
	Update     bson.M
}

// NewDatastore creates a new mgo.Session with a strong consistency and a write mode of "majority"
func NewDatastore(uri string) error {
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
func Teardown(database, collection, key, value string) error {
	s := session.Copy()
	defer s.Clone()

	if err := s.DB(database).C(collection).Remove(bson.M{key: value}); err != nil {
		if err == mgo.ErrNotFound {
			log.Info("data does not exist, continue", nil)
			return nil
		}
		return err
	}

	return nil
}

func RemoveAll(database, collection string) error {
	s := session.Copy()
	defer s.Clone()
	_, err := s.DB(database).C(collection).RemoveAll(nil)
	if err != nil {
		log.Info("error removing all data", nil)
		return err
	}

	return nil
}

// TeardownMany is a way of cleaning up many documents from mongo instance
func TeardownMany(d *ManyDocs) error {
	s := session.Copy()
	defer s.Clone()

	for _, doc := range d.Docs {
		if err := s.DB(doc.Database).C(doc.Collection).Remove(bson.M{doc.Key: doc.Value}); err != nil {
			if err == mgo.ErrNotFound {
				log.Info("data does not exist, continue", nil)
				continue
			}
			return err
		}
	}

	return nil
}

// Setup is a way of loading in an individual document into a mongo instance
func Setup(database, collection, key, value string, update bson.M) error {
	s := session.Copy()
	defer s.Clone()

	if _, err := s.DB(database).C(collection).Upsert(bson.M{key: value}, update); err != nil {
		log.ErrorC("mongodb datastore error", err, nil)
		return err
	}

	log.Info("SetUp completed", nil)
	return nil
}

// SetupMany is a way of loading in many documents into a mongo instance
func SetupMany(d *ManyDocs) error {
	s := session.Copy()
	defer s.Clone()

	for _, doc := range d.Docs {
		//log.Debug("got in for loop", log.Data{"key": key, "value": doc})
		if _, err := s.DB(doc.Database).C(doc.Collection).Upsert(bson.M{doc.Key: doc.Value}, doc.Update); err != nil {
			log.ErrorC("Unable to create document", err, nil)
			return err
		}
	}

	log.Info("SetUp completed", nil)
	return nil
}
