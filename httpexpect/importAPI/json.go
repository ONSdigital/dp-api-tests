package importAPI

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
