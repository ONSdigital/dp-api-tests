package importAPI

import "gopkg.in/mgo.v2/bson"

var validPOSTCreateJobJSON string = `{
  "recipe": "b944be78-f56d-409b-9ebd-ab2b77ffe187",
  "state": "created",
  "files": [
	{
	  "alias_name": "v4",
	  "url": "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv"
	}
  ]
}`

var validJSON string = `{
  "recipe": "b944be78-f56d-409b-9ebd-ab2b77ffe187",
  "state": "created",
  "files": [
	{
	  "alias_name": "v4",
	  "url": "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv"
	}
  ]
}`

// Invalid Json body without recipe
var invalidJSON string = `
{
  "state": "created",
  "files": [
	{
	  "alias_name": "v4",
	  "url": "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv"
	}
  ]
}`

var validPUTJobJSON string = `{
  "recipe": "b944be78-f56d-409b-9ebd-ab2b77ffe187",
  "state": "submitted",
  "files": [
	{
	  "alias_name": "v4",
	  "url": "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv"
	}
  ]
}`

// Invalid Syntax Json body
var invalidSyntaxJSON string = `
{
  "state": "created",
  "files": [
	{
	  "alias_name": "v4",
	  "url": "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv"
	
  ]
}`

var validPUTAddFilesJSON string = `{

	  "alias_name": "v5",
	  "url": "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv"
	
}`

// Files represents an object containing files information
type Files struct {
	AliasName string
	URL       string
}

// GenericObject represents a generic object structure
type GenericObject struct {
	ID   string
	HRef string
}

var files = Files{
	AliasName: "v4",
	URL:       "https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/CPIGrowth.csv",
}

var instances = GenericObject{
	ID:   instanceID,
	HRef: "http://localhost:22000/instances/" + instanceID,
}

var validCreatedImportJobData = bson.M{
	"$set": bson.M{
		"id":              jobID,
		"recipe":          "2080CACA-1A82-411E-AA46-F00804968E78",
		"state":           "Created",
		"files":           []Files{files},
		"links.instances": []GenericObject{instances},
		"links.self.id":   jobID,
		"links.self.href": "http://localhost:22000/jobs/" + jobID,
		"last_updated":    "2017-12-11", // TODO this should be an isodate
	},
}

var validSubmittedImportJobData = bson.M{
	"$set": bson.M{

		"id":              "01C24F0D-24BE-479F-962B-C76BCCD0AD00",
		"recipe":          "6C9D2696-131F-40C3-B598-12200C90415C",
		"state":           "Submitted",
		"files":           []Files{files},
		"links.instances": []GenericObject{instances},
		"links.self.id":   "01C24F0D-24BE-479F-962B-C76BCCD0AD00",
		"links.self.href": "http://localhost:22000/jobs/01C24F0D-24BE-479F-962B-C76BCCD0AD00",
		"last_updated":    "2017-12-11", // TODO this should be an isodate
	},
}
