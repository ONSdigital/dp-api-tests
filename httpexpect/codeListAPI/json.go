package codeListAPI

import "gopkg.in/mgo.v2/bson"

var validFirstCodeListData = bson.M{
	"$set": bson.M{

		"_id":              firstCodeListID,
		"name":             "First Code List",
		"links.self.id":    firstCodeListID,
		"links.self.href":  "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E",
		"links.codes.id":   "code",
		"links.codes.href": "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E/codes",
		"test_data":        "true",
	},
}

var validFirstCodeListFirstCodeData = bson.M{
	"$set": bson.M{

		"_id":                  firstCodeListFirstCodeID,
		"label":                "First Code List label one",
		"code":                 "LS_00998877",
		"links.self.href":      "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E",
		"links.code_list.id":   firstCodeListID,
		"links.code_list.href": "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E/codes",
		"test_data":            "true",
	},
}

var validFirstCodeListSecondCodeData = bson.M{
	"$set": bson.M{

		"_id":                  firstCodeListSecondCodeID,
		"label":                "First Code List label two",
		"code":                 "LS_00998811",
		"links.self.href":      "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E",
		"links.code_list.id":   firstCodeListID,
		"links.code_list.href": "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E/codes",
		"test_data":            "true",
	},
}

var validFirstCodeListThirdCodeData = bson.M{
	"$set": bson.M{

		"_id":                  firstCodeListThirdCodeID,
		"label":                "First Code List label three",
		"code":                 "LS_00998822",
		"links.self.href":      "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E",
		"links.code_list.id":   firstCodeListID,
		"links.code_list.href": "http://localhost:22400/code-lists/1C322128-3FD5-44F0-BBAD-619779D8960E/codes",
		"test_data":            "true",
	},
}

var validSecondCodeListData = bson.M{
	"$set": bson.M{

		"_id":              secondCodeListID,
		"name":             "Second Code List",
		"links.self.id":    secondCodeListID,
		"links.self.href":  "http://localhost:22400/code-lists/C5FA175A-7EA0-4B39-B252-7B52BE75C9DE",
		"links.codes.id":   "code",
		"links.codes.href": "http://localhost:22400/code-lists/C5FA175A-7EA0-4B39-B252-7B52BE75C9DE/codes",
		"test_data":        "true",
	},
}

var validThirdCodeListData = bson.M{
	"$set": bson.M{

		"_id":              thirdCodelistID,
		"name":             "Third Code List",
		"links.self.id":    thirdCodelistID,
		"links.self.href":  "http://localhost:22400/code-lists/5A561370-9AB5-48A4-A619-BEC996DD0BDA",
		"links.codes.id":   "code",
		"links.codes.href": "http://localhost:22400/code-lists/5A561370-9AB5-48A4-A619-BEC996DD0BDA/codes",
		"test_data":        "true",
	},
}
