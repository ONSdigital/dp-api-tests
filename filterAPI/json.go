package filterAPI

import "gopkg.in/mgo.v2/bson"

var ValidPublishedInstanceData = bson.M{
	"$set": bson.M{
		"_id":                   instanceID,
		"collection_id":         "108064B3-A808-449B-9041-EA3A2F72CFAA",
		"downloads.csv.url":     "http://localhost:8080/aws/census-2017-1-csv",
		"downloads.csv.size":    "10mb",
		"downloads.xls.url":     "http://localhost:8080/aws/census-2017-1-xls",
		"downloads.xls.size":    "24mb",
		"edition":               "2017",
		"headers":               []string{"time", "geography"},
		"id":                    instanceID,
		"last_updated":          "2017-09-08", // TODO Should be isodate
		"license":               "ONS License",
		"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
		"links.dataset.id":      "123",
		"links.dataset.href":    "http://localhost:8080/datasets/123",
		"links.dimensions.href": "http://localhost:8080/datasets/123/editions/2017/versions/1/dimensions",
		"links.edition.id":      "1",
		"links.edition.href":    "http://localhost:8080/datasets/123/editions/2017",
		"links.self.href":       "http://localhost:8080/instances/" + instanceID,
		"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
		"links.version.id":      "1",
		"release_date":          "2017-12-12", // TODO Should be isodate
		"state":                 "published",
		"total_inserted_observations": 1000,
		"total_observations":          1000,
		"version":                     1,
		"test_data":                   "true",
	},
}

type Dimension struct {
	DimensionURL string   `bson:"dimension_url"`
	Name         string   `bson:"name"`
	Options      []string `bson:"options"`
}

var dimension = Dimension{
	DimensionURL: "",
	Name:         "age",
	Options:      []string{"27", "28"},
}

var ValidCreatedFilterJob = bson.M{
	"$set": bson.M{
		"_id":                filterID,
		"dimension_list_url": "http://localhost:8080/instances/321/dimensions",
		"dimensions":         []Dimension{dimension},
		"filter_job_id":      filterJobID,
		"instance_id":        instanceID,
		"links.version.id":   "1",
		"links.version.href": "http://localhost:8080/datasets/123/editions/2017/versions/1",
		"state":              "created",
		"test_data":          "true",
	},
}

var ageDimension = Dimension{
	Name:    "age",
	Options: []string{"27"},
}

var sexDimension = Dimension{
	Name:    "sex",
	Options: []string{"male", "female"},
}
var goodsAndServicesDimension = Dimension{
	Name:    "Goods and services",
	Options: []string{"Education", "health", "communication"},
}

var timeDimension = Dimension{
	Name:    "time",
	Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
}

var ValidFilterJobWithMultipleDimensions = bson.M{
	"$set": bson.M{
		"_id":                filterID,
		"dimension_list_url": "http://localhost:8080/instances/321/dimensions",
		"dimensions":         []Dimension{ageDimension, sexDimension, goodsAndServicesDimension, timeDimension},
		"instance_id":        instanceID,
		"filter_job_id":      filterJobID,
		"links.version.id":   "1",
		"links.version.href": "http://localhost:8080/datasets/123/editions/2017/versions/1",
		"state":              "created",
		"test_data":          "true",
	},
}

var ValidSubmittedFilterJob = bson.M{
	"$set": bson.M{
		"_id":                filterID,
		"dimension_list_url": "http://localhost:8080/instances/321/dimensions",
		"dimensions":         []Dimension{ageDimension, sexDimension, goodsAndServicesDimension, timeDimension},
		"instance_id":        instanceID,
		"filter_job_id":      filterJobID,
		"links.version.id":   "1",
		"links.version.href": "http://localhost:8080/datasets/123/editions/2017/versions/1",
		"state":              "submitted",
		"test_data":          "true",
	},
}

var ValidPOSTCreateFilterJSON string = `{
	"instance_id": "` + instanceID + `" ,
	"state": "created",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		]
	  }
	]
  }`

// Invalid Json body without dataset filter id
var InvalidJSON = `
{
	"state": "created",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		]
	  }
	]
	}`

// Valid Json body with state and new dimension options
var ValidPUTUpdateFilterJobJSON = `
{
	"instance_id": "` + instanceID + `" ,
	"state": "submitted",
	"dimensions": [
	  {
		"name": "sex",
		"options": [
		  "intersex", "other"
		]
	  }
	]
	}`

// Invalid Syntax Json body
var InvalidSyntaxJSON = `
{
	"instance_id": "` + instanceID + `" ,
	"state": "created",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"

	  }
	]
	}`

var ValidPOSTMultipleDimensionsCreateFilterJSON = `{
	"instance_id": "` + instanceID + `" ,
	"state": "created",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27"
		]
	  },
	  {
		"name": "sex",
		"options": [
		  "male", "female"
		]
	  },
	  {
		"name": "Goods and services",
		"options": [
		  "Education", "health", "communication"
		]
	  },
	  {
		"name": "time",
		"options": [
		  "March 1997", "April 1997", "June 1997", "September 1997", "December 1997"
		]
	  }
	]
	}`

var ValidPOSTAddDimensionToFilterJobJSON = `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"
  ]
}`

var InvalidPOSTAddDimensionToFilterJobJSON = `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"

}`

var ValidPOSTCreateFilterSubmittedJobJSON = `{
	"instance_id": "` + instanceID + `" ,
	"state": "submitted",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		]
	  }
	]
  }`
