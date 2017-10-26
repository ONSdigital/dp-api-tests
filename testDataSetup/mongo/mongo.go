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
	defer s.Close()

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
	defer s.Close()
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
	defer s.Close()

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
	defer s.Close()

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
	defer s.Close()

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

// FilterJob represents a structure for a filter job
type FilterJob struct {
	InstanceID       string      `bson:"instance_id"          json:"instance_id"`
	DimensionListURL string      `bson:"dimension_list_url"   json:"dimension_list_url,omitempty"`
	Dimensions       []Dimension `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
	Downloads        *Downloads  `bson:"downloads,omitempty" json:"downloads,omitempty"`
	Events           *Events     `bson:"events,omitempty"    json:"events,omitempty"`
	FilterID         string      `bson:"filter_job_id"        json:"filter_job_id,omitempty"`
	State            string      `bson:"state"                json:"state,omitempty"`
	Links            LinkMap     `bson:"links"                json:"links,omitempty"`
}

// LinkMap contains a named LinkObject for each link to other resources
type LinkMap struct {
	Version LinkObject `bson:"version" json:"version,omitempty"`
}

// LinkObject represents a generic structure for all links
type LinkObject struct {
	ID   string `bson:"id"   json:"id,omitempty"`
	HRef string `bson:"href" json:"href,omitempty"`
}

// Dimension represents an object containing a list of dimension values and the dimension name
type Dimension struct {
	DimensionURL string   `bson:"dimension_url"           json:"dimension_url"`
	Name         string   `bson:"name"                    json:"name"`
	Options      []string `bson:"options"                 json:"options"`
}

// Downloads represents a list of file types possible to download
type Downloads struct {
	CSV  *DownloadItem `bson:"csv,omitempty"  json:"csv,omitempty"`
	JSON *DownloadItem `bson:"json,omitempty" json:"json,omitempty"`
	XLS  *DownloadItem `bson:"xls,omitempty"  json:"xls,omitempty"`
}

// DownloadItem represents an object containing information for the download item
type DownloadItem struct {
	Size string `bson:"size" json:"size"`
	URL  string `bson:"url"  json:"url"`
}

// Events represents a list of array objects containing event information against the filter job
type Events struct {
	Error *[]EventItem `bson:"error,omitempty" json:"error,omitempty"`
	Info  *[]EventItem `bson:"info,omitempty"  json:"info,omitempty"`
}

// EventItem represents an event object containing event information
type EventItem struct {
	Message string `bson:"message,omitempty" json:"message,omitempty"`
	Time    string `bson:"time,omitempty"    json:"time,omitempty"`
	Type    string `bson:"type,omitempty"    json:"type,omitempty"`
}

// GetFilterJob retrieves a document from mongo
func GetFilterJob(database, collection, key, value string) (FilterJob, error) {
	s := session.Copy()
	defer s.Close()

	var filterJob FilterJob
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&filterJob); err != nil {
		return filterJob, err
	}

	return filterJob, nil
}
