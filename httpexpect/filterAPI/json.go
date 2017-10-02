package filterAPI

var validPOSTCreateFilterJSON string = `{
	"dataset_filter_id": "census",
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
var invalidJSON string = `
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

var validPUTUpdateFilterJobJSON string = `
{
	"dataset_filter_id": "census",
	"state": "submitted",
	"dimensions": [
	  {
		"name": "sex",
		"options": [
		  "male", "female"
		]
	  }
	]
	}`

// Invalid Syntax Json body
var invalidSyntaxJSON string = `
{
	"dataset_filter_id": "census",
	"state": "created",
	"dimensions": [
	  {
		"name": "age",
		"options": [
		  "27", "28"
		
	  }
	]
	}`

var validPOSTMultipleDimensionsCreateFilterJSON string = `{
	"dataset_filter_id": "census 123",
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

var validPOSTAddDimensionToFilterJobJSON string = `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"
  ]
}`

var invalidPOSTAddDimensionToFilterJobJSON string = `
{
  "options": [
    "Lives in a communal establishment", "Lives in a household"
  
}`

var validPOSTCreateFilterSubmittedJobJSON string = `{
	"dataset_filter_id": "census",
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
