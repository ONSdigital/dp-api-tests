package generateFiles

func createValidJobJSON(recipe, location string) string {
	body := `{
		"recipe": "` + recipe + `",
		"state": "created",
		"files": [{
			"alias_name": "CPIH",
			"url": "` + location + `"
		}]
	}`

	return body
}

func GetValidPOSTCreateFilterJSON(instanceID string) string {
	return `{
	"instance_id": "` + instanceID + `" ,
	"dimensions": [
  	{
	  	"name": "geography",
			"options": ["K02000001"]
		},
		{
			"name": "aggregate",
			"options": ["cpi1dim1A0","cpi1dim1G10100","cpi1dim1G10200","cpi1dim1G20100","cpi1dim1G20200","cpi1dim1G30100","cpi1dim1G30200","cpi1dim1G40100","cpi1dim1G40300","cpi1dim1G40400","cpi1dim1G40500","cpi1dim1G50100","cpi1dim1G50200","cpi1dim1G50300","cpi1dim1G50400","cpi1dim1G50500","cpi1dim1G50600","cpi1dim1G60100","cpi1dim1G70100","cpi1dim1G70200","cpi1dim1G70300","cpi1dim1G80100","cpi1dim1G80200",
				"cpi1dim1G90100","cpi1dim1G90300","cpi1dim1G90400","cpi1dim1G90500","cpi1dim1G90600","cpi1dim1G100000","cpi1dim1G110100","cpi1dim1G110200","cpi1dim1G120100","cpi1dim1G120300","cpi1dim1G120500","cpi1dim1G120600","cpi1dim1G120700","cpi1dim1S10101", "cpi1dim1S10102"]
		},
		{
			"name": "time",
			"options": ["Jan-96"]
		}
	]
}`
}

var validPOSTCreateDatasetJSON = `
{
  "contacts": [
    {
      "email": "cpi@onstest.gov.uk",
      "name": "Automation Tester",
      "telephone": "+44 (0)1633 123456"
    }
  ],
  "description": "Comprehensive database of time series covering measures of inflation data including CPIH, CPI and RPI.",
  "keywords": [
    "cpi"
  ],
	"license": "ONS license",
	"links": {
		"access_rights": {
			"href": "http://ons.gov.uk/accessrights"
		}
	},
  "methodologies": [
    {
      "description": "Consumer price inflation is the rate at which the prices of the goods and services bought by households rise or fall, and is estimated by using consumer price indices.",
      "href": "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
      "title": "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)"
    }
  ],
  "national_statistic": true,
  "next_release": "17 October 2017",
  "publications": [
	  {
		  "description": "Price indices, percentage changes and weights for the different measures of consumer price inflation.",
      "href": "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017",
			"title": "UK consumer price inflation: August 2017"
		}
	],
	"publisher": {
	  "name": "Automation Tester",
		"type": "publisher",
		"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/bulletins/consumerpriceinflation/aug2017"
	},
	"qmi": {
	  "description": "Consumer price inflation is the rate at which the prices of goods and services bought by households rise and fall",
		"href": "https://www.ons.gov.uk/economy/inflationandpriceindices/qmis/consumerpriceinflationqmi",
	  "title": "Consumer Price Inflation (includes all 3 indices – CPIH, CPI and RPI)"
	},
	"related_datasets": [
	  {
		  "href": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceindices",
			"title": "Consumer Price Inflation time series dataset"
		}
	],
	"release_frequency": "Monthly",
	"state": "created",
	"theme": "Goods and services",
	"title": "CPI",
	"unit_of_measure": "Pounds Sterling",
	"uri": "https://www.ons.gov.uk/economy/inflationandpriceindices/datasets/consumerpriceinflation"
}`

var validPUTInstanceMetadataJSON = `
{
	"alerts": [
	  {
		  "date": "2017-04-05",
		  "description": "All data entries (observations) for Plymouth have been updated",
			"type": "Correction"
	  }
	],
	"edition": "2017",
	"latest_changes": [
	  {
		  "description": "change to the period frequency from quarterly to monthly",
			"name": "Changes to the period frequency",
			"type": "Summary of Changes"
	  }
	],
	"links": {
		"spatial": {
			"href": "http://ons.gov.uk/geography-list"
		}
	},
	"release_date": "2017-11-11",
  "state": "edition-confirmed",
	"temporal": [
		{
			"start_date": "2014-10-10",
			"end_date": "2016-10-10",
			"frequency": "monthly"
		}
	]
}`

var validPUTUpdateVersionToAssociatedJSON = `
{
  "collection_id": "308064B3-A808-449B-9041-EA3A2F72CFAC",
  "state": "associated"
}`
