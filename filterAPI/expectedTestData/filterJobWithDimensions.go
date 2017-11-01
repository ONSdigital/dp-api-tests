package expectedTestData

import "github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"

func age(host, filterBlueprintID string) mongo.Dimension {
	if filterBlueprintID == "" {
		return mongo.Dimension{
			Name:    "age",
			Options: []string{"27"},
		}
	}

	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/age",
		Name:    "age",
		Options: []string{"27"},
	}
}

func sex(host, filterBlueprintID string) mongo.Dimension {
	if filterBlueprintID == "" {
		return mongo.Dimension{
			Name:    "sex",
			Options: []string{"male", "female"},
		}
	}

	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/sex",
		Name:    "sex",
		Options: []string{"male", "female"},
	}
}

func goodsAndServices(host, filterBlueprintID string) mongo.Dimension {
	if filterBlueprintID == "" {
		return mongo.Dimension{
			Name:    "Goods and services",
			Options: []string{"Education", "health", "communication"},
		}
	}

	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/Goods and services",
		Name:    "Goods and services",
		Options: []string{"Education", "health", "communication"},
	}
}

func time(host, filterBlueprintID string) mongo.Dimension {
	if filterBlueprintID == "" {
		return mongo.Dimension{
			Name:    "time",
			Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
		}
	}

	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/time",
		Name:    "time",
		Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
	}
}

func residenceType(host, filterBlueprintID string) mongo.Dimension {
	if filterBlueprintID == "" {
		return mongo.Dimension{
			Name:    "Residence Type",
			Options: []string{"Lives in a communal establishment", "Lives in a household"},
		}
	}

	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/Residence Type",
		Name:    "Residence Type",
		Options: []string{"Lives in a communal establishment", "Lives in a household"},
	}
}

// ExpectedFilterJob represents the expected data stored against a filter blueprint resource
func ExpectedFilterBlueprint(host, instanceID, filterBlueprintID string) mongo.Filter {
	return mongo.Filter{
		Dimensions: []mongo.Dimension{
			age(host, filterBlueprintID),
			sex(host, filterBlueprintID),
			goodsAndServices(host, filterBlueprintID),
			time(host, filterBlueprintID),
			residenceType(host, filterBlueprintID),
		},
		Links: mongo.LinkMap{
			Dimensions: mongo.LinkObject{
				HRef: host + "/filters/" + filterBlueprintID + "/dimensions",
			},
			Self: mongo.LinkObject{
				HRef: host + "/filters/" + filterBlueprintID,
			},
			Version: mongo.LinkObject{
				ID:   "1",
				HRef: "http://localhost:8080/datasets/123/editions/2017/versions/1",
			},
		},
	}
}

// ExpectedFilterOutput represents the expected data stored against a filter output resource
func ExpectedFilterOutput(host, instanceID, filterOutputID string) mongo.Filter {
	return mongo.Filter{
		FilterID:   filterOutputID,
		InstanceID: instanceID,
		Dimensions: []mongo.Dimension{
			age(host, ""),
			sex(host, ""),
			goodsAndServices(host, ""),
			time(host, ""),
		},
		Downloads: &mongo.Downloads{
			CSV: &mongo.DownloadItem{
				Size: "",
				URL:  "",
			},
			JSON: &mongo.DownloadItem{
				Size: "",
				URL:  "",
			},
			XLS: &mongo.DownloadItem{
				Size: "",
				URL:  "",
			},
		},
		Links: mongo.LinkMap{
			Self: mongo.LinkObject{
				HRef: host + "/filter-outputs/" + filterOutputID,
			},
			Version: mongo.LinkObject{
				ID:   "1",
				HRef: "http://localhost:8080/datasets/123/editions/2017/versions/1",
			},
		},
		State: "created",
	}
}

func updatedAge(host, filterBlueprintID string) mongo.Dimension {
	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/age",
		Name:    "age",
		Options: []string{"27", "28"},
	}
}

func updatedSex(host, filterBlueprintID string) mongo.Dimension {
	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/sex",
		Name:    "sex",
		Options: []string{"male", "female", "unknown"},
	}
}

func updatedGoodsAndServices(host, filterBlueprintID string) mongo.Dimension {
	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/Goods and services",
		Name:    "Goods and services",
		Options: []string{"Education", "health", "communication", "welfare"},
	}
}

func updatedTime(host, filterBlueprintID string) mongo.Dimension {
	return mongo.Dimension{
		URL:     host + "/filters/" + filterBlueprintID + "/dimensions/time",
		Name:    "time",
		Options: []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997", "February 2007"},
	}
}

// ExpectedFilterBlueprintUpdated represents the expected data stored against a filter job with dimensions
func ExpectedFilterBlueprintUpdated(host, instanceID, filterBlueprintID string) mongo.Filter {
	return mongo.Filter{
		Dimensions: []mongo.Dimension{
			updatedAge(host, filterBlueprintID),
			updatedSex(host, filterBlueprintID),
			updatedGoodsAndServices(host, filterBlueprintID),
			updatedTime(host, filterBlueprintID),
		},
		Links: mongo.LinkMap{
			Dimensions: mongo.LinkObject{
				HRef: host + "/filters/" + filterBlueprintID + "/dimensions",
			},
			Self: mongo.LinkObject{
				HRef: host + "/filters/" + filterBlueprintID,
			},
			Version: mongo.LinkObject{
				ID:   "1",
				HRef: "http://localhost:8080/datasets/123/editions/2017/versions/1",
			},
		},
	}
}
