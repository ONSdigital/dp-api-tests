package filterAPI

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

func GetValidPublishedInstanceData(instanceID string) bson.M {
	return bson.M{
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
}

type Dimension struct {
	URL     string   `bson:"dimension_url"`
	Name    string   `bson:"name"`
	Options []string `bson:"options"`
}

func dimension(host, filterBlueprintID string) Dimension {
	return Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/age",
		Name:    "age",
		Options: []string{"27", "28"},
	}
}

func GetValidCreatedFilterBlueprint(host, filterID, instanceID, filterBlueprintID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id": filterID,
			"dimensions": []Dimension{
				dimension(host, filterBlueprintID),
			},
			"filter_id":             filterBlueprintID,
			"instance_id":           instanceID,
			"links.dimensions.href": host + "/filters/" + filterBlueprintID + "/dimensions",
			"links.self.href":       host + "/filters/" + filterBlueprintID,
			"links.version.id":      "1",
			"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"test_data":             "true",
		},
	}
}

func ageDimension(host, filterBlueprintID string) Dimension {
	return Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/age",
		Name:    "age",
		Options: []string{"27"},
	}
}

func sexDimension(host, filterBlueprintID string) Dimension {
	return Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/sex",
		Name:    "sex",
		Options: []string{"male", "female"},
	}
}

func goodsAndServicesDimension(host, filterBlueprintID string) Dimension {
	return Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/Goods and services",
		Name:    "Goods and services",
		Options: []string{"Education", "health", "communication"},
	}
}

func timeDimension(host, filterBlueprintID string) Dimension {
	return Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/time",
		Name:    "time",
		Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
	}
}

func GetValidFilterWithMultipleDimensions(host, filterID, instanceID, filterBlueprintID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                   filterID,
			"dimensions":            []Dimension{ageDimension(host, filterBlueprintID), sexDimension(host, filterBlueprintID), goodsAndServicesDimension(host, filterBlueprintID), timeDimension(host, filterBlueprintID)},
			"instance_id":           instanceID,
			"filter_id":             filterBlueprintID,
			"links.dimensions.href": host + "/filters/" + filterBlueprintID + "/dimensions",
			"links.self.href":       host + "/filters/" + filterBlueprintID,
			"links.version.id":      "1",
			"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"test_data":             "true",
		},
	}
}

func GetValidSubmittedFilterJob(host, filterID, instanceID, filterBlueprintID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                   filterID,
			"dimensions":            []Dimension{ageDimension(host, filterBlueprintID), sexDimension(host, filterBlueprintID), goodsAndServicesDimension(host, filterBlueprintID), timeDimension(host, filterBlueprintID)},
			"instance_id":           instanceID,
			"filter_id":             filterBlueprintID,
			"links.dimensions.href": host + "/filters/" + filterBlueprintID + "/dimensions",
			"links.self.href":       host + "/filters/" + filterBlueprintID,
			"links.version.id":      "1",
			"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"test_data":             "true",
		},
	}
}

func GetValidPOSTCreateFilterJSON(instanceID string) string {
	return `{
	"instance_id": "` + instanceID + `" ,
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		]
	  }
	]
  }`
}

// Invalid Json body without dataset filter id
func GetInvalidJSON(instanceID string) string {
	return `
{
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		]
	  }
	]
	}`
}

// GetValidPUTUpdateFilterBlueprintJSON Json body with state and new dimension options
func GetValidPUTUpdateFilterBlueprintJSON(instanceID string) string {
	return `
{
	"instance_id": "` + instanceID + `" ,
	"dimensions": [
	  {
		"name": "sex",
		"options": [
		  "intersex", "other"
		]
	  }
	]
	}`
}

// Invalid Syntax Json body
func GetInvalidSyntaxJSON(instanceID string) string {
	return `
{
	"instance_id": "` + instanceID + `" ,
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"

	  }
	]
	}`
}

func GetValidPOSTMultipleDimensionsCreateFilterJSON(instanceID string) string {
	return `{
	"instance_id": "` + instanceID + `" ,
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
}

func GetValidPOSTAddDimensionToFilterBlueprintJSON() string {
	return `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"
  ]
}`
}

func GetInvalidPOSTAddDimensionToFilterBlueprintJSON() string {
	return `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"

}`
}

func GetValidPOSTCreateFilterSubmittedJobJSON(instanceID string) string {
	return `{
	"instance_id": "` + instanceID + `" ,
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		]
	  }
	]
  }`
}

func GetValidPUTFilterBlueprintJSON(instanceID string, time time.Time) string {
	return `{
	  "instance_id": "` + instanceID + `",
	  "events": {
		  "info": [
		    {
		      "message": "blueprint has created filter output resource",
					"time": "` + time.String() + `",
					"type": "info"
	      }
	    ]
		}
  }`
}
