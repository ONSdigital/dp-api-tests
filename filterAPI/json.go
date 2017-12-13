package filterAPI

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

func GetValidPublishedInstanceDataBSON(instanceID string) bson.M {
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
			"links.edition.id":      "2017",
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

func GetValidCreatedFilterBlueprintBSON(host, filterID, instanceID, filterBlueprintID string) bson.M {
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

func ageDimension(host, filterID string) Dimension {
	if filterID == "" {
		return Dimension{
			Name:    "age",
			Options: []string{"27"},
		}
	}
	return Dimension{
		URL:     host + "/filters/" + filterID + "/dimensions/age",
		Name:    "age",
		Options: []string{"27"},
	}
}

func sexDimension(host, filterID string) Dimension {
	if filterID == "" {
		return Dimension{
			Name:    "sex",
			Options: []string{"male", "female"},
		}
	}
	return Dimension{
		URL:     host + "/filters/" + filterID + "/dimensions/sex",
		Name:    "sex",
		Options: []string{"male", "female"},
	}
}

func goodsAndServicesDimension(host, filterID string) Dimension {
	if filterID == "" {
		return Dimension{
			Name:    "Goods and services",
			Options: []string{"Education", "health", "communication"},
		}
	}
	return Dimension{
		URL:     host + "/filters/" + filterID + "/dimensions/Goods and services",
		Name:    "Goods and services",
		Options: []string{"Education", "health", "communication"},
	}
}

func timeDimension(host, filterID string) Dimension {
	if filterID == "" {
		return Dimension{
			Name:    "time",
			Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
		}
	}
	return Dimension{
		URL:     host + "/filters/" + filterID + "/dimensions/time",
		Name:    "time",
		Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
	}
}

func GetValidFilterWithMultipleDimensionsBSON(host, filterID, instanceID, filterBlueprintID string) bson.M {
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

func GetValidFilterOutputWithMultipleDimensionsBSON(host, filterID, instanceID, filterOutputID, filterBlueprintID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                         filterID,
			"dimensions":                  []Dimension{ageDimension(host, ""), sexDimension(host, ""), goodsAndServicesDimension(host, ""), timeDimension(host, "")},
			"downloads.csv.url":           "s3-csv-location",
			"downloads.csv.size":          "12mb",
			"downloads.json.url":          "s3-json-location",
			"downloads.json.size":         "6mb",
			"downloads.xls.url":           "s3-xls-location",
			"downloads.xls.size":          "24mb",
			"instance_id":                 instanceID,
			"filter_id":                   filterOutputID,
			"links.filter_blueprint.href": host + "/filters/" + filterBlueprintID,
			"links.filter_blueprint.id":   filterBlueprintID,
			"links.self.href":             host + "/filter-outputs/" + filterOutputID,
			"links.version.id":            "1",
			"links.version.href":          "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"state":                       "completed",
			"test_data":                   "true",
		},
	}
}

func GetValidFilterOutputWithoutDownloadsBSON(host, filterID, instanceID, filterOutputID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                filterID,
			"dimensions":         []Dimension{ageDimension(host, ""), sexDimension(host, ""), goodsAndServicesDimension(host, ""), timeDimension(host, "")},
			"instance_id":        instanceID,
			"filter_id":          filterOutputID,
			"links.self.href":    host + "/filters/" + filterOutputID,
			"links.version.id":   "1",
			"links.version.href": "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"state":              "created",
			"test_data":          "true",
		},
	}
}

func GetValidAge27DimensionData(instanceID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "age",
			"option":               "27",
			"label":                "age 27",
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code.id":        "27",
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/27",
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func GetValidAgeDimensionData(instanceID, option string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "age",
			"option":               option,
			"label":                "age " + option,
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a",
			"links.code.id":        option,
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58a/codes/" + option,
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func GetValidSexDimensionData(instanceID, option string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "sex",
			"option":               option,
			"label":                "sex " + option,
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58b",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58b",
			"links.code.id":        option,
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58b/codes/" + option,
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func GetValidGoodsAndServicesDimensionData(instanceID, option string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "Goods and services",
			"option":               option,
			"label":                "Goods and services " + option,
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58c",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58c",
			"links.code.id":        option,
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58c/codes/" + option,
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func GetValidTimeDimensionData(instanceID, option string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "time",
			"option":               option,
			"label":                "time" + option,
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58d",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58d",
			"links.code.id":        option,
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58d/codes/" + option,
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
		},
	}
}

func GetValidResidenceTypeDimensionData(instanceID, option string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "Residence Type",
			"option":               option,
			"label":                "Residence Type " + option,
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58e",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58e",
			"links.code.id":        option,
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58e/codes/" + option,
			"last_updated":         "2017-09-09", // TODO Should be isodate
			"test_data":            "true",
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
		  "27", "42"
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
		  "27", "42"
		]
	  }
	]
	}`
}

// GetInvalidDimensionJSON contains an invalid dimension for instance
func GetInvalidDimensionJSON(instanceID string) string {
	return `
{
	"instance_id": "` + instanceID + `",
	"dimensions": [
	  {
		"name": "weight",
		"options": [
		  "27", "42"
		]
	  }
	]
	}`
}

// GetInvalidDimensionOptionJSON contains an invalid dimension for instance
func GetInvalidDimensionOptionJSON(instanceID string) string {
	return `
{
	"instance_id": "` + instanceID + `",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "33"
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

func GetValidPOSTDimensionToFilterBlueprintJSON() string {
	return `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"
  ]
}`
}

func GetInvalidPOSTDimensionToFilterBlueprintJSON() string {
	return `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"

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

func GetValidPUTFilterOutputWithCSVDownloadJSON() string {
	return `{
	  "downloads": {
			"csv": {
			  "url": "s3-csv-location",
				"size": "12mb"
		  }
		}
  }`
}

func GetValidPUTFilterOutputWithXLSDownloadJSON() string {
	return `{
	  "downloads": {
			"xls": {
			  "url": "s3-xls-location",
				"size": "24mb"
		  }
		}
  }`
}

func GetValidPUTFilterOutputWithDimensionsJSON() string {
	return `{
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
