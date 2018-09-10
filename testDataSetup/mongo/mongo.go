package mongo

import (
	"github.com/ONSdigital/dp-api-tests/identityAPIModels"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	datasetAPIModel "github.com/ONSdigital/dp-dataset-api/models"
	importAPIModel "github.com/ONSdigital/dp-import-api/models"
	"github.com/ONSdigital/go-ns/log"
)

var session *mgo.Session

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

// Teardown is a way of cleaning up any number of documents from mongo instance
func Teardown(d ...*Doc) error {
	s := session.Copy()
	defer s.Close()

	for _, doc := range d {
		if err := s.DB(doc.Database).C(doc.Collection).Remove(bson.M{doc.Key: doc.Value}); err != nil {
			if err == mgo.ErrNotFound {
				log.Info("data does not exist, continue", log.Data{
					"database":   doc.Database,
					"collection": doc.Collection,
					"key":        doc.Key,
					"value":      doc.Value,
				})
				continue
			}
			return err
		}
	}

	return nil
}

// Setup is a way of loading any number of documents into a mongo instance
func Setup(d ...*Doc) error {
	if err := Teardown(d...); err != nil {
		log.ErrorC("Unable to teardown previous document", err, nil)
		return err
	}

	s := session.Copy()
	defer s.Close()

	for _, doc := range d {
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

// GetJob retrieves a job document from mongo
func GetJob(database, collection, key, value string) (importAPIModel.Job, error) {
	s := session.Copy()
	defer s.Close()

	var job importAPIModel.Job
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

//Dataset contains all the metadata which does not change across editions and versions of a dataset
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

// EditionUpdate represents an evolving edition containing both the next and current edition
type EditionUpdate struct {
	ID      string   `bson:"_id,omitempty"         json:"id,omitempty"`
	Current *Edition `bson:"current,omitempty"     json:"current,omitempty"`
	Next    *Edition `bson:"next,omitempty"        json:"next,omitempty"`
}

// Edition represents information related to a single edition for a dataset
type Edition struct {
	Edition   string        `bson:"edition,omitempty"      json:"edition,omitempty"`
	ID        string        `bson:"id,omitempty"          json:"id,omitempty"`
	Links     *EditionLinks `bson:"links,omitempty"        json:"links,omitempty"`
	State     string        `bson:"state,omitempty"        json:"state,omitempty"`
	time.Time `bson:"last_updated,omitempty" json:"-"`
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
	Alerts        *[]Alert            `bson:"alerts,omitempty"         json:"alerts,omitempty"`
	CollectionID  string              `bson:"collection_id,omitempty"  json:"collection_id,omitempty"`
	Dimensions    []CodeList          `bson:"dimensions,omitempty"     json:"dimensions,omitempty"`
	Downloads     *DownloadList       `bson:"downloads,omitempty"      json:"downloads,omitempty"`
	Edition       string              `bson:"edition,omitempty"        json:"edition,omitempty"`
	ID            string              `bson:"id,omitempty"             json:"id,omitempty"`
	LatestChanges []LatestChange      `bson:"latest_changes,omitempty" json:"latest_changes,omitempty"`
	Links         *VersionLinks       `bson:"links,omitempty"          json:"links,omitempty"`
	ReleaseDate   string              `bson:"release_date,omitempty"   json:"release_date,omitempty"`
	State         string              `bson:"state,omitempty"          json:"state,omitempty"`
	Temporal      []TemporalFrequency `bson:"temporal,omitempty"       json:"temporal,omitempty"`
	LastUpdated   time.Time           `bson:"last_updated,omitempty"   json:"-"`
	Version       int                 `bson:"version,omitempty"        json:"version,omitempty"`
	UsageNotes    *[]UsageNote        `bson:"usage_notes,omitempty"     json:"usage_notes,omitempty"`
}

type UsageNote struct {
	Title string `bson:"title,omitempty"    json:"title,omitempty"`
	Note  string `bson:"note,omitempty"     json:"note,omitempty"`
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
	Label       string `json:"label"`
	Name        string `json:"name"`
}

// DownloadList represents a list of objects of containing information on the downloadable files
type DownloadList struct {
	CSV  *DownloadObject `bson:"csv,omitempty" json:"csv,omitempty"`
	CSVW *DownloadObject `bson:"csvw,omitempty" json:"csvw,omitempty"`
	XLS  *DownloadObject `bson:"xls,omitempty" json:"xls,omitempty"`
}

// DownloadObject represents information on the downloadable file
type DownloadObject struct {
	URL     string `bson:"href,omitempty"     json:"href,omitempty"`
	Size    string `bson:"size,omitempty"    json:"size,omitempty"`
	Public  string `bson:"public,omitempty"  json:"public,omitempty"`
	Private string `bson:"private,omitempty" json:"private,omitempty"`
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
	Alerts            *[]Alert             `bson:"alerts,omitempty"                      json:"alerts,omitempty"`
	InstanceID        string               `bson:"id,omitempty"                          json:"id,omitempty"`
	CollectionID      string               `bson:"collection_id,omitempty"               json:"collection_id,omitempty"`
	Dimensions        []CodeList           `bson:"dimensions,omitempty"                  json:"dimensions,omitempty"`
	Downloads         *DownloadList        `bson:"downloads,omitempty"                   json:"downloads,omitempty"`
	Edition           string               `bson:"edition,omitempty"                     json:"edition,omitempty"`
	Events            *[]InstanceEvent     `bson:"events,omitempty"                      json:"events,omitempty"`
	Headers           *[]string            `bson:"headers,omitempty"                     json:"headers,omitempty"`
	ImportTasks       *InstanceImportTasks `bson:"import_tasks,omitempty"                json:"import_tasks,omitempty"`
	LatestChanges     []LatestChange       `bson:"latest_changes,omitempty"              json:"latest_changes,omitempty"`
	Links             InstanceLinks        `bson:"links,omitempty"                       json:"links,omitempty"`
	ReleaseDate       string               `bson:"release_date,omitempty"                json:"release_date,omitempty"`
	State             string               `bson:"state,omitempty"                       json:"state,omitempty"`
	Temporal          []TemporalFrequency  `bson:"temporal,omitempty"                    json:"temporal,omitempty"`
	TotalObservations int64                `bson:"total_observations,omitempty"          json:"total_observations,omitempty"`
	Version           int                  `bson:"version,omitempty"                     json:"version,omitempty"`
	LastUpdated       time.Time            `bson:"last_updated,omitempty"                json:"last_updated,omitempty"`
	UniqueTimestamp   bson.MongoTimestamp  `bson:"unique_timestamp"                      json:"-"`
}

// InstanceImportTasks represent an object containing specific lists of tasks for import process
type InstanceImportTasks struct {
	ImportObservations  *ImportObservationsTask `bson:"import_observations,omitempty"  json:"import_observations,omitempty"`
	BuildHierarchyTasks []*BuildHierarchyTask   `bson:"build_hierarchies,omitempty"    json:"build_hierarchies,omitempty"`
	SearchTasks         []*BuildSearchIndexTask `bson:"build_search_indexes,omitempty" json:"build_search_indexes,omitempty"`
}

// ImportObservationsTask represents the task of importing instance observation data into the database.
type ImportObservationsTask struct {
	State                string `bson:"state,omitempty"                       json:"state,omitempty"`
	InsertedObservations int64  `bson:"total_inserted_observations,omitempty" json:"total_inserted_observations"`
}

// BuildHierarchyTask represents a task of importing a single hierarchy.
type BuildHierarchyTask struct {
	State         string `bson:"state,omitempty"          json:"state,omitempty"`
	DimensionName string `bson:"dimension_name,omitempty" json:"dimension_name,omitempty"`
	CodeListID    string `bson:"code_list_id,omitempty"   json:"code_list_id,omitempty"`
}

// BuildSearchIndexTask represents a task of importing a single hierarchy into search.
type BuildSearchIndexTask struct {
	State         string `bson:"state,omitempty"          json:"state,omitempty"`
	DimensionName string `bson:"dimension_name,omitempty" json:"dimension_name,omitempty"`
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

// InstanceEvent represents the event structure against an instance
// until it is deprecated to the event structure below
type InstanceEvent struct {
	Message       string `bson:"message,omitempty" json:"message,omitempty"`
	MessageOffset string `bson:"message_offset,omitempty" json:"message_offset,omitempty"`
	Type          string `bson:"type,omitempty" json:"type,omitempty"`
}

// Event structure
type Event struct {
	Type string    `bson:"type,omitempty" json:"type"`
	Time time.Time `bson:"time,omitempty" json:"time"`
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
func GetEdition(database, collection, key, value string) (EditionUpdate, error) {
	s := session.Copy()
	defer s.Close()

	var edition EditionUpdate
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

// GetInstance retrieves an instance document from mongo
func GetInstance(database, collection, key, value string) (Instance, error) {
	s := session.Copy()
	defer s.Close()

	var instance Instance
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&instance); err != nil {
		return instance, err
	}

	return instance, nil
}

// GetDimensionOption retrieves a dimension option document from mongo
func GetDimensionOption(database, collection, key, value string) (dimensionOption datasetAPIModel.DimensionOption, err error) {
	s := session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Find(bson.M{key: value}).One(&dimensionOption); err != nil {
		return
	}

	return
}

// CountDimensionOptions retrieves a count of the number of dimension options exist for an instance in mongo
func CountDimensionOptions(database, collection, key, value string) (int, error) {
	s := session.Copy()
	defer s.Close()

	return s.DB(database).C(collection).Find(bson.M{key: value}).Count()
}

func GetIdentity(database, collection, key, value string) (*identityAPIModels.Mongo, error) {
	s := session.Copy()
	defer s.Close()

	var i identityAPIModels.Mongo
	if err := s.DB(database).C(collection).Find(bson.M{key: value}).One(&i); err != nil {
		return nil, err
	}

	return &i, nil
}

func GetIdentities(database, collection string) ([]identityAPIModels.Mongo, error) {
	s := session.Copy()
	defer s.Close()

	var results []identityAPIModels.Mongo
	if err := s.DB(database).C(collection).Find(nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// Possible values for flagging whether a filter resource (output or blueprint)
// is a filter against a published or unpublished version
var (
	Published   = true
	Unpublished = false
)

// Filter represents a structure for a filter blueprint or output
type Filter struct {
	InstanceID      string              `bson:"instance_id"          json:"instance_id"`
	UniqueTimestamp bson.MongoTimestamp `bson:"unique_timestamp"     json:"-"`
	Dimensions      []Dimension         `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
	Downloads       *Downloads          `bson:"downloads,omitempty"  json:"downloads,omitempty"`
	Events          []*Event            `bson:"events,omitempty"     json:"events,omitempty"`
	FilterID        string              `bson:"filter_id"            json:"filter_id,omitempty"`
	State           string              `bson:"state,omitempty"      json:"state,omitempty"`
	Links           LinkMap             `bson:"links"                json:"links,omitempty"`
	Published       *bool               `bson:"published,omitempty"  json:"published,omitempty"`
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
	CSV *DownloadItem `bson:"csv,omitempty"  json:"csv,omitempty"`
	XLS *DownloadItem `bson:"xls,omitempty"  json:"xls,omitempty"`
}

// DownloadItem represents an object containing information for the download item
type DownloadItem struct {
	HRef    string `bson:"href,omitempty"    json:"href,omitempty"`
	Private string `bson:"private,omitempty" json:"private,omitempty"`
	Public  string `bson:"public,omitempty"  json:"public,omitempty"`
	Size    string `bson:"size,omitempty"    json:"size,omitempty"`
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
