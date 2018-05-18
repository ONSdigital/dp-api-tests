package filterAPI

import (
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func GetValidPublishedInstanceDataBSON(instanceID, datasetID, edition string, version int) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                   instanceID,
			"collection_id":         "108064B3-A808-449B-9041-EA3A2F72CFAA",
			"dataset.id":            datasetID,
			"dataset.edition":       edition,
			"dataset.version":       version,
			"downloads.csv.href":    "http://localhost:8080/aws/census-2017-1-csv",
			"downloads.csv.size":    "10mb",
			"downloads.xls.href":    "http://localhost:8080/aws/census-2017-1-xls",
			"downloads.xls.size":    "24mb",
			"edition":               "2017",
			"headers":               []string{"time", "geography"},
			"id":                    instanceID,
			"last_updated":          "2017-09-08", // TODO Should be isodate
			"license":               "ONS License",
			"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":      datasetID,
			"links.dataset.href":    "http://localhost:8080/datasets/" + datasetID,
			"links.dimensions.href": "http://localhost:8080/datasets/" + datasetID + "/editions/" + edition + "/versions/" + strconv.Itoa(version) + "/dimensions",
			"links.edition.id":      edition,
			"links.edition.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/" + edition,
			"links.self.href":       "http://localhost:8080/instances/" + instanceID,
			"links.version.href":    "http://localhost:8080/datasets/" + datasetID + "/editions/" + edition + "/versions/" + strconv.Itoa(version),
			"links.version.id":      strconv.Itoa(version),
			"release_date":          "2017-12-12", // TODO Should be isodate
			"state":                 "published",
			"total_inserted_observations": 1000,
			"total_observations":          1000,
			"version":                     version,
			"test_data":                   "true",
		},
	}
}

func GetUnpublishedInstanceDataBSON(instanceID string, datasetID, edition string, version int) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                   instanceID,
			"collection_id":         "108064B3-A808-449B-9041-EA3A2F72CFAA",
			"dataset.id":            datasetID,
			"dataset.edition":       edition,
			"dataset.version":       version,
			"downloads.csv.href":    "http://localhost:8080/aws/census-2017-1-csv",
			"downloads.csv.size":    "10mb",
			"downloads.xls.href":    "http://localhost:8080/aws/census-2017-1-xls",
			"downloads.xls.size":    "24mb",
			"edition":               "2017",
			"headers":               []string{"time", "geography"},
			"id":                    instanceID,
			"last_updated":          "2017-09-08", // TODO Should be isodate
			"license":               "ONS License",
			"links.job.id":          "042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.job.href":        "http://localhost:8080/jobs/042e216a-7822-4fa0-a3d6-e3f5248ffc35",
			"links.dataset.id":      datasetID,
			"links.dataset.href":    "http://localhost:8080/datasets/123",
			"links.dimensions.href": "http://localhost:8080/datasets/123/editions/2017/versions/1/dimensions",
			"links.edition.id":      edition,
			"links.edition.href":    "http://localhost:8080/datasets/123/editions/2017",
			"links.self.href":       "http://localhost:8080/instances/" + instanceID,
			"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"links.version.id":      "1",
			"release_date":          "2017-12-12", // TODO Should be isodate
			"state":                 "associated",
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

func GetValidCreatedFilterBlueprintBSON(host, filterID, instanceID, filterBlueprintID, datasetID, edition string, version int) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":             filterID,
			"dataset.id":      datasetID,
			"dataset.edition": edition,
			"dataset.version": version,
			"dimensions": []Dimension{
				dimension(host, filterBlueprintID),
			},
			"filter_id":             filterBlueprintID,
			"instance_id":           instanceID,
			"links.dimensions.href": host + "/filters/" + filterBlueprintID + "/dimensions",
			"links.self.href":       host + "/filters/" + filterBlueprintID,
			"links.version.id":      "1",
			"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"published":             true,
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
			Name:    "aggregate",
			Options: []string{"cpi1dim1T60000", "cpi1dim1S10201", "cpi1dim1S10105"},
		}
	}
	return Dimension{
		URL:     host + "/filters/" + filterID + "/dimensions/aggregate",
		Name:    "aggregate",
		Options: []string{"cpi1dim1T60000", "cpi1dim1S10201", "cpi1dim1S10105"},
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

func GetValidFilterWithMultipleDimensionsBSON(host, filterID, instanceID, datasetID, edition, filterBlueprintID string, version int, published bool) bson.M {

	return bson.M{
		"$set": bson.M{
			"_id":                   filterID,
			"dataset.id":            datasetID,
			"dataset.edition":       edition,
			"dataset.version":       version,
			"dimensions":            []Dimension{ageDimension(host, filterBlueprintID), sexDimension(host, filterBlueprintID), goodsAndServicesDimension(host, filterBlueprintID), timeDimension(host, filterBlueprintID)},
			"instance_id":           instanceID,
			"filter_id":             filterBlueprintID,
			"links.dimensions.href": host + "/filters/" + filterBlueprintID + "/dimensions",
			"links.self.href":       host + "/filters/" + filterBlueprintID,
			"links.version.id":      "1",
			"links.version.href":    "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"published":             published,
			"test_data":             "true",
		},
	}
}

func GetValidFilterOutputWithMultipleDimensionsBSON(host, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition string, version int, published bool) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                         filterID,
			"dataset.id":                  datasetID,
			"dataset.edition":             edition,
			"dataset.version":             version,
			"dimensions":                  []Dimension{ageDimension(host, ""), sexDimension(host, ""), goodsAndServicesDimension(host, ""), timeDimension(host, "")},
			"downloads.csv.href":          "download-service-url.csv",
			"downloads.csv.private":       "private-s3-csv-location",
			"downloads.csv.public":        "https://s3-eu-west-1.amazonaws.com/dp-frontend-florence-file-uploads/2470609-cpicoicoptestcsv",
			"downloads.csv.size":          "12mb",
			"downloads.xls.href":          "download-service-url.xlsx",
			"downloads.xls.private":       "private-s3-xls-location",
			"downloads.xls.public":        "public-s3-xls-location",
			"downloads.xls.size":          "24mb",
			"instance_id":                 instanceID,
			"filter_id":                   filterOutputID,
			"links.filter_blueprint.href": host + "/filters/" + filterBlueprintID,
			"links.filter_blueprint.id":   filterBlueprintID,
			"links.self.href":             host + "/filter-outputs/" + filterOutputID,
			"links.version.id":            "1",
			"links.version.href":          "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"state":                       "completed",
			"published":                   published,
			"test_data":                   "true",
		},
	}
}

func GetValidFilterOutputWithPrivateDownloads(host, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition string, version int, published bool) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                         filterID,
			"dataset.id":                  datasetID,
			"dataset.edition":             edition,
			"dataset.version":             version,
			"dimensions":                  []Dimension{ageDimension(host, ""), sexDimension(host, ""), goodsAndServicesDimension(host, ""), timeDimension(host, "")},
			"downloads.csv.href":          "download-service-url.csv",
			"downloads.csv.private":       "s3://csv-exported/v4TestFile.csv",
			"downloads.csv.size":          "12mb",
			"downloads.xls.href":          "download-service-url.xlsx",
			"downloads.xls.private":       "s3://csv-exported/v4TestFile.xls",
			"downloads.xls.size":          "24mb",
			"instance_id":                 instanceID,
			"filter_id":                   filterOutputID,
			"links.filter_blueprint.href": host + "/filters/" + filterBlueprintID,
			"links.filter_blueprint.id":   filterBlueprintID,
			"links.self.href":             host + "/filter-outputs/" + filterOutputID,
			"links.version.id":            "1",
			"links.version.href":          "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"state":                       "completed",
			"published":                   published,
			"test_data":                   "true",
		},
	}
}

func GetValidFilterOutputBSON(host, filterID, instanceID, filterOutputID, filterBlueprintID, datasetID, edition, csvPublicLink, xlsPublicLink string, version int, dimension Dimension) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                         filterID,
			"dataset.id":                  datasetID,
			"dataset.edition":             edition,
			"dataset.version":             version,
			"dimensions":                  []Dimension{dimension},
			"downloads.csv.href":          "download-service-url.csv",
			"downloads.csv.private":       "private-s3-csv-location",
			"downloads.csv.public":        csvPublicLink,
			"downloads.csv.size":          "12mb",
			"downloads.xls.href":          "download-service-url.xlsx",
			"downloads.xls.private":       "private-s3-xls-location",
			"downloads.xls.public":        xlsPublicLink,
			"downloads.xls.size":          "24mb",
			"filter_id":                   filterOutputID,
			"instance_id":                 instanceID,
			"links.filter_blueprint.href": host + "/filters/" + filterBlueprintID,
			"links.filter_blueprint.id":   filterBlueprintID,
			"links.self.href":             host + "/filter-outputs/" + filterOutputID,
			"links.version.id":            "1",
			"links.version.href":          "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"published":                   true,
			"state":                       "completed",
			"test_data":                   "true",
		},
	}
}

func GetValidFilterOutputNoDimensionsBSON(host, filterID, instanceID, filterOutputID, filterBlueprintID string) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                         filterID,
			"downloads.csv.href":          "download-service-url.csv",
			"downloads.csv.size":          "12mb",
			"downloads.json.size":         "6mb",
			"downloads.xls.href":          "download-service-url.xlsx",
			"downloads.xls.size":          "24mb",
			"instance_id":                 instanceID,
			"filter_id":                   filterOutputID,
			"links.filter_blueprint.href": host + "/filters/" + filterBlueprintID,
			"links.filter_blueprint.id":   filterBlueprintID,
			"links.self.href":             host + "/filter-outputs/" + filterOutputID,
			"links.version.id":            "1",
			"links.version.href":          "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"state":                       "completed",
			"published":                   true,
			"test_data":                   "true",
		},
	}
}

func GetValidFilterOutputWithoutDownloadsBSON(host, filterID, instanceID, filterOutputID, datasetID, edition string, version int) bson.M {
	return bson.M{
		"$set": bson.M{
			"_id":                filterID,
			"dataset.id":         datasetID,
			"dataset.edition":    edition,
			"dataset.version":    version,
			"dimensions":         []Dimension{ageDimension(host, ""), sexDimension(host, ""), goodsAndServicesDimension(host, ""), timeDimension(host, "")},
			"instance_id":        instanceID,
			"filter_id":          filterOutputID,
			"links.self.href":    host + "/filters/" + filterOutputID,
			"links.version.id":   "1",
			"links.version.href": "http://localhost:8080/datasets/123/editions/2017/versions/1",
			"state":              "created",
			"published":          true,
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

func GetValidAggregateDimensionData(instanceID, option string) bson.M {
	return bson.M{
		"$set": bson.M{
			"instance_id":          instanceID,
			"name":                 "aggregate",
			"option":               option,
			"label":                "aggregate " + option,
			"links.code_list.id":   "64d384f1-ea3b-445c-8fb8-aa453f96e58f",
			"links.code_list.href": cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58f",
			"links.code.id":        option,
			"links.code.href":      cfg.DatasetAPIURL + "/code-lists/64d384f1-ea3b-445c-8fb8-aa453f96e58f/codes/" + option,
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

func GetValidPOSTCreateFilterJSON(datasetID, edition string, version int) string {
	return `{
	"dataset": {
		"id": "` + datasetID + `",
		"edition": "` + edition + `",
		"version": ` + strconv.Itoa(version) + `
	},
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
func GetInvalidJSON() string {
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
func GetInvalidDimensionJSON(datasetID, edition string, version int) string {
	return `
	{
		"dataset": {
			"id": "` + datasetID + `",
			"edition": "` + edition + `",
			"version": ` + strconv.Itoa(version) + `
		},
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
func GetInvalidDimensionOptionJSON(datasetID, edition string, version int) string {
	return `
	{
		"dataset": {
			"id": "` + datasetID + `",
			"edition": "` + edition + `",
			"version": ` + strconv.Itoa(version) + `
		},
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

func GetValidPUTFilterBlueprintJSON(version int, time time.Time) string {
	return `{
	  "dataset": {
			"version": ` + strconv.Itoa(version) + `
		},
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

func GetInValidPUTFilterBlueprintJSON(id, edition string, version int, time time.Time) string {
	json := `{"dataset": {`

	if id != "" {
		json = json + `"id": "` + id + `",`
	}

	if edition != "" {
		json = json + `"edition": "` + edition + `",`
	}

	json = json + `"version": ` + strconv.Itoa(version) + `}}`

	return json
}

func GetValidPUTFilterOutputWithCSVDownloadJSON() string {
	return `{
	  "downloads": {
			"csv": {
			  "href": "download-service-url.csv",
				"private": "private-s3-csv-location",
				"size": "12mb"
		  }
		}
  }`
}

func GetValidPUTFilterOutputWithCSVPublicLinkJSON() string {
	return `{
		"downloads": {
			"csv" : {
			  "public": "public-s3-csv-location"
		  }
	  }
	}`
}

func GetValidPUTFilterOutputWithXLSDownloadJSON() string {
	return `{
	  "downloads": {
			"xls": {
			  "href": "download-service-url.xlsx",
				"private": "private-s3-xls-location",
				"size": "24mb"
		  }
		}
  }`
}

func GetValidPUTFilterOutputWithXLSPublicLinkJSON() string {
	return `{
		"downloads": {
			"xls" : {
			  "public": "public-s3-xls-location"
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
