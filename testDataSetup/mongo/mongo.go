package mongo

import (
	"time"

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

// DropDatabases cleans out all data by removing the databases specified
func DropDatabases(databases []string) error {
	log.Info("the following databases about to be dropped in mongo", log.Data{"databases": databases})

	s := session.Copy()
	defer s.Close()

	for _, db := range databases {
		log.Info("dropping database", log.Data{"database": db})
		if err := s.DB(db).DropDatabase(); err != nil {
			return err
		}
	}

	return nil
}

// Teardown is a way of cleaning up an individual document from mongo instance
func Teardown(database, collection, key, value string) error {
	s := session.Copy()
	defer s.Close()

	if _, err := s.DB(database).C(collection).RemoveAll(bson.M{key: value}); err != nil {
		if err == mgo.ErrNotFound {
			log.Info("data does not exist, continue", nil)
			return nil
		}
		return err
	}

	return nil
}

// TeardownAll removes all documents from collection
func TeardownAll(database, collection string) error {
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

// ------------------------------------------------------------------------

// Job for importing datasets
type Job struct {
	ID            string          `bson:"id,omitempty"             json:"id,omitempty"`
	RecipeID      string          `bson:"recipe,omitempty"         json:"recipe,omitempty"`
	State         string          `bson:"state,omitempty"          json:"state,omitempty"`
	UploadedFiles *[]UploadedFile `bson:"files,omitempty"          json:"files,omitempty"`
	Links         LinksMap        `bson:"links,omitempty"          json:"links,omitempty"`
	LastUpdated   time.Time       `bson:"last_updated,omitempty"   json:"last_updated,omitempty"`
}

// UploadedFile used for a file which has been uploaded to a bucket
type UploadedFile struct {
	AliasName string `bson:"alias_name" json:"alias_name" avro:"alias-name"`
	URL       string `bson:"url"        json:"url"        avro:"url"`
}

// LinksMap represents an object containing a set of links
type LinksMap struct {
	Instances []IDLink `bson:"instances,omitempty" json:"instances,omitempty"`
	Self      IDLink   `bson:"self,omitempty" json:"self,omitempty"`
}

// GetJob retrieves a job document from mongo
func GetJob(database, collection, key, value string) (Job, error) {
	s := session.Copy()
	defer s.Close()

	var job Job
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&job); err != nil {
		return job, err
	}

	return job, nil
}

// DatasetUpdate represents an evolving dataset with the current dataset and the updated dataset
type DatasetUpdate struct {
	ID      string   `bson:"_id,omitempty"         json:"id,omitempty"`
	Current *Dataset `bson:"current,omitempty"     json:"current,omitempty"`
	Next    *Dataset `bson:"next,omitempty"        json:"next,omitempty"`
}

type Dataset struct {
	CollectionID      string           `bson:"collection_id,omitempty"          json:"collection_id,omitempty"`
	Contacts          []ContactDetails `bson:"contacts,omitempty"               json:"contacts,omitempty"`
	Description       string           `bson:"description,omitempty"            json:"description,omitempty"`
	Keywords          []string         `bson:"keywords,omitempty"               json:"keywords,omitempty"`
	ID                string           `bson:"_id,omitempty"                    json:"id,omitempty"`
	License           string           `bson:"license,omitempty"                json:"license,omitempty"`
	Links             *DatasetLinks    `bson:"links,omitempty"                  json:"links,omitempty"`
	Methodologies     []GeneralDetails `bson:"methodologies,omitempty"          json:"methodologies,omitempty"`
	NationalStatistic *bool            `bson:"national_statistic,omitempty"     json:"national_statistic,omitempty"`
	NextRelease       string           `bson:"next_release,omitempty"           json:"next_release,omitempty"`
	Publications      []GeneralDetails `bson:"publications,omitempty"           json:"publications,omitempty"`
	Publisher         *Publisher       `bson:"publisher,omitempty"              json:"publisher,omitempty"`
	QMI               *GeneralDetails  `bson:"qmi,omitempty"                    json:"qmi,omitempty"`
	RelatedDatasets   []GeneralDetails `bson:"related_datasets,omitempty"       json:"related_datasets,omitempty"`
	ReleaseFrequency  string           `bson:"release_frequency,omitempty"      json:"release_frequency,omitempty"`
	State             string           `bson:"state,omitempty"                  json:"state,omitempty"`
	Theme             string           `bson:"theme,omitempty"                  json:"theme,omitempty"`
	Title             string           `bson:"title,omitempty"                  json:"title,omitempty"`
	UnitOfMeasure     string           `bson:"unit_of_measure,omitempty"        json:"unit_of_measure,omitempty"`
	URI               string           `bson:"uri,omitempty"                    json:"uri,omitempty"`
}

// ContactDetails represents an object containing information of the contact
type ContactDetails struct {
	Email     string `bson:"email,omitempty"      json:"email,omitempty"`
	Name      string `bson:"name,omitempty"       json:"name,omitempty"`
	Telephone string `bson:"telephone,omitempty"  json:"telephone,omitempty"`
}

// DatasetLinks represents a list of specific links related to the dataset resource
type DatasetLinks struct {
	AccessRights  *LinkObject `bson:"access_rights,omitempty"   json:"access_rights,omitempty"`
	Editions      *LinkObject `bson:"editions,omitempty"        json:"editions,omitempty"`
	LatestVersion *LinkObject `bson:"latest_version,omitempty"  json:"latest_version,omitempty"`
	Self          *LinkObject `bson:"self,omitempty"            json:"self,omitempty"`
}

// GeneralDetails represents generic fields stored against an object (reused)
type GeneralDetails struct {
	Description string `bson:"description,omitempty"    json:"description,omitempty"`
	HRef        string `bson:"href,omitempty"           json:"href,omitempty"`
	Title       string `bson:"title,omitempty"          json:"title,omitempty"`
}

// LinkObject represents a generic structure for all links
type LinkObject struct {
	ID   string `bson:"id,omitempty"    json:"id,omitempty"`
	HRef string `bson:"href,omitempty"  json:"href,omitempty"`
}

// Publisher represents an object containing information of the publisher
type Publisher struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Type string `bson:"type,omitempty" json:"type,omitempty"`
	HRef string `bson:"href,omitempty" json:"href,omitempty"`
}

// Edition represents information related to a single edition for a dataset
type Edition struct {
	Edition     string        `bson:"edition,omitempty"      json:"edition,omitempty"`
	ID          string        `bson:"id,omitempty"          json:"id,omitempty"`
	Links       *EditionLinks `bson:"links,omitempty"        json:"links,omitempty"`
	State       string        `bson:"state,omitempty"        json:"state,omitempty"`
	LastUpdated time.Time     `bson:"last_updated,omitempty" json:"-"`
}

// EditionLinks represents a list of specific links related to the edition resource of a dataset
type EditionLinks struct {
	Dataset       *LinkObject `bson:"dataset,omitempty"        json:"dataset,omitempty"`
	LatestVersion *LinkObject `bson:"latest_version,omitempty" json:"latest_version,omitempty"`
	Self          *LinkObject `bson:"self,omitempty"           json:"self,omitempty"`
	Versions      *LinkObject `bson:"versions,omitempty"       json:"versions,omitempty"`
}

// Version represents information related to a single version for an edition of a dataset
type Version struct {
	Alerts        *[]Alert             `bson:"alerts,omitempty"         json:"alerts,omitempty"`
	CollectionID  string               `bson:"collection_id,omitempty"  json:"collection_id,omitempty"`
	Dimensions    []CodeList           `bson:"dimensions,omitempty"     json:"dimensions,omitempty"`
	Downloads     *DownloadList        `bson:"downloads,omitempty"      json:"downloads,omitempty"`
	Edition       string               `bson:"edition,omitempty"        json:"edition,omitempty"`
	ID            string               `bson:"id,omitempty"             json:"id,omitempty"`
	LatestChanges *[]LatestChange      `bson:"latest_changes,omitempty" json:"latest_changes,omitempty"`
	Links         *VersionLinks        `bson:"links,omitempty"          json:"links,omitempty"`
	ReleaseDate   string               `bson:"release_date,omitempty"   json:"release_date,omitempty"`
	State         string               `bson:"state,omitempty"          json:"state,omitempty"`
	Temporal      *[]TemporalFrequency `bson:"temporal,omitempty"       json:"temporal,omitempty"`
	LastUpdated   time.Time            `bson:"last_updated,omitempty"   json:"-"`
	Version       int                  `bson:"version,omitempty"        json:"version,omitempty"`
}

// Alert represents an object containing information on an alert
type Alert struct {
	Date        string `bson:"date,omitempty"        json:"date,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Type        string `bson:"type,omitempty"        json:"type,omitempty"`
}

// CodeList for a dimension within an instance
type CodeList struct {
	Description string `json:"description"`
	HRef        string `json:"href"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}

// DownloadList represents a list of objects of containing information on the downloadable files
type DownloadList struct {
	CSV *DownloadObject `bson:"csv,omitempty" json:"csv,omitempty"`
	XLS *DownloadObject `bson:"xls,omitempty" json:"xls,omitempty"`
}

// DownloadObject represents information on the downloadable file
type DownloadObject struct {
	URL  string `bson:"url,omitempty"  json:"url,omitempty"`
	Size string `bson:"size,omitempty" json:"size,omitempty"`
}

// TemporalFrequency represents a frequency for a particular period of time
type TemporalFrequency struct {
	EndDate   string `bson:"end_date,omitempty"    json:"end_date,omitempty"`
	Frequency string `bson:"frequency,omitempty"   json:"frequency,omitempty"`
	StartDate string `bson:"start_date,omitempty"  json:"start_date,omitempty"`
}

// LatestChange represents an object contining
// information on a single change between versions
type LatestChange struct {
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Name        string `bson:"name,omitempty"        json:"name,omitempty"`
	Type        string `bson:"type,omitempty"        json:"type,omitempty"`
}

// VersionLinks represents a list of specific links related to the version resource for an edition of a dataset
type VersionLinks struct {
	Dataset    *LinkObject `bson:"dataset,omitempty"     json:"dataset,omitempty"`
	Dimensions *LinkObject `bson:"dimensions,omitempty"  json:"dimensions,omitempty"`
	Edition    *LinkObject `bson:"edition,omitempty"     json:"edition,omitempty"`
	Self       *LinkObject `bson:"self,omitempty"        json:"self,omitempty"`
	Spatial    *LinkObject `bson:"spatial,omitempty"     json:"spatial,omitempty"`
	Version    *LinkObject `bson:"version,omitempty"     json:"-"`
}

// Instance which presents a single dataset being imported
type Instance struct {
	Alerts               *[]Alert             `bson:"alerts,omitempty"         json:"alerts,omitempty"`
	InstanceID           string               `bson:"id,omitempty"                          json:"id,omitempty"`
	CollectionID         string               `bson:"collection_id,omitempty"               json:"collection_id,omitempty"`
	Dimensions           []CodeList           `bson:"dimensions,omitempty"                  json:"dimensions,omitempty"`
	Downloads            *DownloadList        `bson:"downloads,omitempty"                   json:"downloads,omitempty"`
	Edition              string               `bson:"edition,omitempty"                     json:"edition,omitempty"`
	Events               *[]Event             `bson:"events,omitempty"                      json:"events,omitempty"`
	Headers              *[]string            `bson:"headers,omitempty"                     json:"headers,omitempty"`
	InsertedObservations *int                 `bson:"total_inserted_observations,omitempty" json:"total_inserted_observations,omitempty"`
	LatestChanges        *[]LatestChange      `bson:"latest_changes,omitempty" json:"latest_changes,omitempty"`
	Links                InstanceLinks        `bson:"links,omitempty"                       json:"links,omitempty"`
	ReleaseDate          string               `bson:"release_date,omitempty"                json:"release_date,omitempty"`
	State                string               `bson:"state,omitempty"                       json:"state,omitempty"`
	Temporal             *[]TemporalFrequency `bson:"temporal,omitempty"                    json:"temporal,omitempty"`
	TotalObservations    *int                 `bson:"total_observations,omitempty"          json:"total_observations,omitempty"`
	Version              int                  `bson:"version,omitempty"                     json:"version,omitempty"`
	LastUpdated          time.Time            `bson:"last_updated,omitempty"                json:"last_updated,omitempty"`
}

// InstanceLinks holds all links for an instance
type InstanceLinks struct {
	Job        *IDLink `bson:"job,omitempty"        json:"job"`
	Dataset    *IDLink `bson:"dataset,omitempty"    json:"dataset,omitempty"`
	Dimensions *IDLink `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
	Edition    *IDLink `bson:"edition,omitempty"    json:"edition,omitempty"`
	Version    *IDLink `bson:"version,omitempty"    json:"version,omitempty"`
	Self       *IDLink `bson:"self,omitempty"       json:"self,omitempty"`
	Spatial    *IDLink `bson:"spatial,omitempty"    json:"spatial,omitempty"`
}

// IDLink holds the id and a link to the resource
type IDLink struct {
	ID   string `bson:"id,omitempty"   json:"id,omitempty"`
	HRef string `bson:"href,omitempty" json:"href,omitempty"`
}

// Event which has happened to an instance
type Event struct {
	Type          string     `bson:"type,omitempty"           json:"type"`
	Time          *time.Time `bson:"time,omitempty"           json:"time"`
	Message       string     `bson:"message,omitempty"        json:"message"`
	MessageOffset string     `bson:"message_offset,omitempty" json:"message_offset"`
}

// GetDataset retrieves a dataset document from mongo
func GetDataset(database, collection, key, value string) (DatasetUpdate, error) {
	s := session.Copy()
	defer s.Close()

	var dataset DatasetUpdate
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&dataset); err != nil {
		return dataset, err
	}

	return dataset, nil
}

// GetEdition retrieves an edition document from mongo
func GetEdition(database, collection, key, value string) (Edition, error) {
	s := session.Copy()
	defer s.Close()

	var edition Edition
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&edition); err != nil {
		return edition, err
	}

	return edition, nil
}

// GetVersion retrieves a version document from mongo
func GetVersion(database, collection, key, value string) (Version, error) {
	s := session.Copy()
	defer s.Close()

	var version Version
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&version); err != nil {
		return version, err
	}

	return version, nil
}

// GetInstance retrieves a version document from mongo
func GetInstance(database, collection, key, value string) (Instance, error) {
	s := session.Copy()
	defer s.Close()

	var instance Instance
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&instance); err != nil {
		return instance, err
	}

	return instance, nil
}

// CountDimensionOptions retrieves a count of the number of dimension options exist for an instance in mongo
func CountDimensionOptions(database, collection, key, value string) (int, error) {
	s := session.Copy()
	defer s.Close()

	var count int
	count, err := s.DB(database).C(collection).Find(bson.M{key: value}).Count()
	if err != nil {
		return count, err
	}

	return count, nil
}

// Filter represents a structure for a filter blueprint or output
type Filter struct {
	InstanceID string      `bson:"instance_id"          json:"instance_id"`
	Dimensions []Dimension `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
	Downloads  *Downloads  `bson:"downloads,omitempty"  json:"downloads,omitempty"`
	Events     *Events     `bson:"events,omitempty"     json:"events,omitempty"`
	FilterID   string      `bson:"filter_id"            json:"filter_id,omitempty"`
	State      string      `bson:"state,omitempty"      json:"state,omitempty"`
	Links      LinkMap     `bson:"links"                json:"links,omitempty"`
}

// LinkMap contains a named LinkObject for each link to other resources
type LinkMap struct {
	Dimensions      LinkObject `bson:"dimensions"       json:"dimensions,omitempty"`
	FilterBlueprint LinkObject `bson:"filter_blueprint" json:"filter_blueprint,omitempty"`
	Self            LinkObject `bson:"self"             json:"self,omitempty"`
	Version         LinkObject `bson:"version"          json:"version,omitempty"`
}

// Dimension represents an object containing a list of dimension values and the dimension name
type Dimension struct {
	URL     string   `bson:"dimension_url"           json:"dimension_url"`
	Name    string   `bson:"name"                    json:"name"`
	Options []string `bson:"options"                 json:"options"`
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

// Events represents a list of array objects containing event information against the filter blueprint or output
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

// GetFilter retrieves a filter document from mongo
func GetFilter(database, collection, key, value string) (Filter, error) {
	s := session.Copy()
	defer s.Close()

	var filter Filter
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&filter); err != nil {
		return filter, err
	}

	return filter, nil
}
