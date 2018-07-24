package helpers

import (
	"time"

	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
)

// IDsAndTimestamps represents a list of ids and a mongo timestamp
type IDsAndTimestamps struct {
	DatasetPublished         string
	DatasetAssociated        string
	Dimension                string
	EditionPublished         string
	EditionUnpublished       string
	InstancePublished        string
	InstanceAssociated       string
	InstanceEditionConfirmed string
	InstanceCompleted        string
	InstanceSubmitted        string
	InstanceCreated          string
	InstanceInvalid          string
	Node                     string
	UniqueTimestamp          bson.MongoTimestamp
}

// GetIDsAndTimestamps returns an object containing a list of unique ids and timestamps
func GetIDsAndTimestamps() (it IDsAndTimestamps, err error) {
	it = IDsAndTimestamps{
		DatasetPublished:         uuid.NewV4().String(),
		DatasetAssociated:        uuid.NewV4().String(),
		Dimension:                uuid.NewV4().String(),
		EditionPublished:         uuid.NewV4().String(),
		EditionUnpublished:       uuid.NewV4().String(),
		InstancePublished:        uuid.NewV4().String(),
		InstanceAssociated:       uuid.NewV4().String(),
		InstanceEditionConfirmed: uuid.NewV4().String(),
		InstanceCompleted:        uuid.NewV4().String(),
		InstanceSubmitted:        uuid.NewV4().String(),
		InstanceCreated:          uuid.NewV4().String(),
		InstanceInvalid:          uuid.NewV4().String(),
		Node:                     uuid.NewV4().String(),
	}

	it.UniqueTimestamp, err = bson.NewMongoTimestamp(time.Now().UTC(), 1)
	if err != nil {
		return
	}

	return
}
