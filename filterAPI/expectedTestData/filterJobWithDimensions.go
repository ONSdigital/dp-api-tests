package expectedTestData

import "github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"

var age = mongo.Dimension{
	DimensionURL: "",
	Name:         "age",
	Options:      []string{"27"},
}

var sex = mongo.Dimension{
	DimensionURL: "",
	Name:         "sex",
	Options:      []string{"male", "female"},
}

var goodsAndServices = mongo.Dimension{
	DimensionURL: "",
	Name:         "Goods and services",
	Options:      []string{"Education", "health", "communication"},
}

var time = mongo.Dimension{
	DimensionURL: "",
	Name:         "time",
	Options:      []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997"},
}

var residenceType = mongo.Dimension{
	DimensionURL: "http://localhost:22100/filter/321/dimensions/Residence Type",
	Name:         "Residence Type",
	Options:      []string{"Lives in a communal establishment", "Lives in a household"},
}

// ExpectedFilterJob represents the expected data stored against a filter job with dimensions
var ExpectedFilterJob = mongo.FilterJob{
	DimensionListURL: "http://localhost:8080/instances/321/dimensions",
	FilterID:         "321",
	InstanceID:       "789",
	Dimensions: []mongo.Dimension{
		age,
		sex,
		goodsAndServices,
		time,
		residenceType,
	},
	State: "created",
	Links: mongo.LinkMap{
		Version: mongo.LinkObject{
			ID:   "1",
			HRef: "http://localhost:8080/datasets/123/editions/2017/versions/1",
		},
	},
}

var updatedAge = mongo.Dimension{
	DimensionURL: "",
	Name:         "age",
	Options:      []string{"27", "28"},
}

var updatedSex = mongo.Dimension{
	DimensionURL: "",
	Name:         "sex",
	Options:      []string{"male", "female", "unknown"},
}

var updatedGoodsAndServices = mongo.Dimension{
	DimensionURL: "",
	Name:         "Goods and services",
	Options:      []string{"Education", "health", "communication", "welfare"},
}

var updatedTime = mongo.Dimension{
	DimensionURL: "",
	Name:         "time",
	Options:      []string{"March 1997", "April 1997", "June 1997", "September 1997", "December 1997", "February 2007"},
}

// ExpectedFilterJobUpdated represents the expected data stored against a filter job with dimensions
var ExpectedFilterJobUpdated = mongo.FilterJob{
	DimensionListURL: "http://localhost:8080/instances/321/dimensions",
	FilterID:         "321",
	InstanceID:       "789",
	Dimensions: []mongo.Dimension{
		updatedAge,
		updatedSex,
		updatedGoodsAndServices,
		updatedTime,
	},
	State: "created",
	Links: mongo.LinkMap{
		Version: mongo.LinkObject{
			ID:   "1",
			HRef: "http://localhost:8080/datasets/123/editions/2017/versions/1",
		},
	},
}
