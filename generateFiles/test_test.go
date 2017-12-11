package generateFiles

import (
	"net/http"
	"os"
	"strings"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfulEndToEndProcess(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)
	recipeAPI := httpexpect.New(t, cfg.RecipeAPIURL)
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)

	hasRemovedAllResources := true
	filename := "v4TestFile.csv"
	recipe := "2943f3c5-c3f1-4a9a-aa6e-14d21c33524c"

	// Get dataset ID from recipe API
	recipeResponse := recipeAPI.GET("/recipes/{recipe}", recipe).
		Expect().Status(http.StatusOK).JSON().Object()

	recipeResponse.Value("id").NotNull()

	Convey("Given a v4 file exists in aws", t, func() {
		// Send v4 file to aws
		_, err := sendV4FileToAWS(region, bucketName, filename)
		if err != nil {
			os.Exit(1)
		}

		// import API expects a s3 url as the location of the file
		location := "s3://" + bucketName + "/" + filename

		Convey("When a job is imported and the version of the dataset is published", func() {

			// STEP 1 - Create job with state created
			postJobResponse := importAPI.POST("/jobs").WithBytes([]byte(createValidJobJSON(recipe, location))).
				Expect().Status(http.StatusCreated).JSON().Object()

			postJobResponse.Value("id").NotNull()
			jobID := postJobResponse.Value("id").String().Raw()

			postJobResponse.Value("files").Array().Element(0).Object().Value("alias_name").Equal("CPIH")
			postJobResponse.Value("files").Array().Element(0).Object().Value("url").Equal("s3://ons-dp-cmd-test/v4TestFile.csv")

			postJobResponse.Value("last_updated").NotNull()
			postJobResponse.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").NotNull()
			postJobResponse.Value("links").Object().Value("self").Object().Value("href").String().Match(cfg.ImportAPIURL + "/jobs/" + jobID + "$")
			postJobResponse.Value("recipe").Equal(recipe)
			postJobResponse.Value("state").Equal("created")

			instanceID := postJobResponse.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").String().Raw()

			// Check for instance creation
			instanceResource, err := mongo.GetInstance("datasets", "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve instance resource", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(len(instanceResource.Dimensions), ShouldEqual, 3)
			So(instanceResource.Links.Job.ID, ShouldEqual, jobID)
			So(instanceResource.Links.Job.HRef, ShouldEqual, cfg.ImportAPIURL+"/jobs/"+jobID)
			So(instanceResource.Links.Dataset.ID, ShouldEqual, datasetName)
			So(instanceResource.Links.Dataset.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName)
			So(instanceResource.Links.Self.HRef, ShouldEqual, cfg.DatasetAPIURL+"/instances/"+instanceID)
			So(instanceResource.State, ShouldEqual, "created")

			// STEP 2 - Create dataset with dataset id from previous response
			postDatasetResponse := datasetAPI.POST("/datasets/{id}", datasetName).WithHeader(internalTokenHeader, internalTokenID).WithBytes([]byte(validPOSTCreateDatasetJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()

			postDatasetResponse.Value("next").Object().Value("links").Object().Value("self").Object().Value("href").String().Match(cfg.DatasetAPIURL + "/datasets/" + datasetName + "$")
			postDatasetResponse.Value("next").Object().Value("state").Equal("created")

			// STEP 3 - Update job state to submitted
			importAPI.PUT("/jobs/{id}", jobID).WithHeader(internalTokenHeader, internalTokenID).WithBytes([]byte(`{"state":"submitted"}`)).
				Expect().Status(http.StatusOK)

			// Check import job state is completed or submitted
			jobResource, err := mongo.GetJob("imports", "imports", "id", jobID)
			if err != nil {
				log.ErrorC("Unable to retrieve job resource", err, log.Data{"job_id": jobID})
				os.Exit(1)
			}

			So(jobResource.State, ShouldNotEqual, "created")

			var stateHasChanged bool
			if jobResource.State == "completed" || jobResource.State == "submitted" {
				stateHasChanged = true
			}

			So(stateHasChanged, ShouldEqual, true)

			// Check instance has updated with headers, state is completed, total_observations and total_inserted_observations
			totalObservations := 1513

			tryAgain := true
			for tryAgain {
				instanceResource, err = mongo.GetInstance("datasets", "instances", "id", instanceID)
				if err != nil {
					log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
					os.Exit(1)
				}
				if instanceResource.State == "completed" {
					tryAgain = false
				} else {
					So(instanceResource.State, ShouldEqual, "submitted")
				}
			}

			So(instanceResource.Headers, ShouldResemble, &[]string{"V4_0", "Time_codelist", "Time", "Geography_codelist", "Geography", "cpi1dim1aggid", "Aggregate"})
			So(instanceResource.InsertedObservations, ShouldResemble, &totalObservations)
			So(instanceResource.State, ShouldEqual, "completed")
			So(instanceResource.TotalObservations, ShouldResemble, &totalObservations)

			// Check dimension options
			count, err := mongo.CountDimensionOptions("datasets", "dimension.options", "instance_id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve dimension option resources", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(count, ShouldEqual, 140)

			// STEP 4 - Update instance with meta data and change state to `edition-confirmed`
			datasetAPI.PUT("/instances/{instance_id}", instanceID).WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(validPUTInstanceMetadataJSON)).Expect().Status(http.StatusOK)

			// Check instance has updated
			instanceResource, err = mongo.GetInstance("datasets", "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve instance resource", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(instanceResource.Alerts, ShouldNotBeNil)
			So(instanceResource.Edition, ShouldEqual, "2017")
			So(instanceResource.LatestChanges, ShouldNotBeNil)
			So(instanceResource.Links.Dimensions.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1/dimensions")
			So(instanceResource.Links.Edition.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017")
			So(instanceResource.Links.Edition.ID, ShouldEqual, "2017")
			So(instanceResource.Links.Spatial.HRef, ShouldEqual, "http://ons.gov.uk/geography-list")
			So(instanceResource.Links.Version.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
			So(instanceResource.Links.Version.ID, ShouldEqual, "1")
			So(instanceResource.ReleaseDate, ShouldEqual, "2017-11-11")
			So(instanceResource.State, ShouldEqual, "edition-confirmed")
			So(instanceResource.Temporal, ShouldNotBeNil)
			So(instanceResource.Version, ShouldEqual, 1)

			// Check Edition has been created
			editionResource, err := mongo.GetEdition("datasets", "editions", "links.self.href", instanceResource.Links.Edition.HRef)
			if err != nil {
				log.ErrorC("Unable to retrieve edition resource", err, log.Data{"links.self.href": instanceResource.Links.Edition.HRef})
				os.Exit(1)
			}

			So(editionResource.Edition, ShouldEqual, "2017")
			So(editionResource.Links.Dataset.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName)
			So(editionResource.Links.Dataset.ID, ShouldEqual, datasetName)
			So(editionResource.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
			So(editionResource.Links.LatestVersion.ID, ShouldEqual, "1")
			So(editionResource.Links.Self.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017")
			So(editionResource.Links.Versions.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions")
			So(editionResource.State, ShouldEqual, "created")

			// STEP 5 - Update version with collection_id and change state to associated
			datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetName, "2017", "1").WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(validPUTUpdateVersionToAssociatedJSON)).Expect().Status(http.StatusOK)

			versionResource, err := mongo.GetVersion("datasets", "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(versionResource.CollectionID, ShouldEqual, "308064B3-A808-449B-9041-EA3A2F72CFAC")
			So(versionResource.State, ShouldEqual, "associated")

			// Check dataset has updated
			datasetResource, err := mongo.GetDataset("datasets", "datasets", "_id", datasetName)
			if err != nil {
				log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
				os.Exit(1)
			}

			So(datasetResource.Current, ShouldBeNil)
			So(datasetResource.Next.CollectionID, ShouldEqual, "308064B3-A808-449B-9041-EA3A2F72CFAC")
			So(datasetResource.Next.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
			So(datasetResource.Next.Links.LatestVersion.ID, ShouldEqual, "1")
			So(datasetResource.Next.State, ShouldEqual, "associated")

			// STEP 6 -  Update version to a state of published
			datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetName, "2017", "1").WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(`{"state":"published"}`)).Expect().Status(http.StatusOK)

			versionResource, err = mongo.GetVersion("datasets", "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(versionResource.State, ShouldEqual, "published")

			// Check edition has updated
			editionResource, err = mongo.GetEdition("datasets", "editions", "links.self.href", instanceResource.Links.Edition.HRef)
			if err != nil {
				log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
				os.Exit(1)
			}

			So(editionResource.State, ShouldEqual, "published")

			// Check dataset has updated
			datasetResource, err = mongo.GetDataset("datasets", "datasets", "_id", datasetName)
			if err != nil {
				log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
				os.Exit(1)
			}

			So(datasetResource.Current, ShouldNotBeNil)
			So(datasetResource.Current.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
			So(datasetResource.Current.State, ShouldEqual, "published")

			Convey("Then an api customer should be able to get a csv and xlsx download link", func() {
				// TODO Get downloads link from version document

				// read file from s3
				// check the number of rows and anything else (e.g. meta data)

			})

			Convey("Then an api customer should be able to filter a dataset and be able to download a csv and xlsx download of the data", func() {
				filterBlueprintResponse := filterAPI.POST("/filters").WithQuery("submitted", "true").
					WithBytes([]byte(GetValidPOSTCreateFilterJSON(instanceID))).Expect().Status(http.StatusCreated).JSON().Object()

				filterBlueprintID := filterBlueprintResponse.Value("filter_id").String().Raw()

				filterBlueprintResponse.Value("filter_id").NotNull()
				filterBlueprintResponse.Value("instance_id").Equal(instanceID)
				filterBlueprintResponse.Value("dimensions").Array().Element(0).Object().Value("name").Equal("geography")
				filterBlueprintResponse.Value("dimensions").Array().Element(0).Object().Value("options").Array().Length().Equal(1)
				filterBlueprintResponse.Value("dimensions").Array().Element(1).Object().Value("name").Equal("aggregate")
				filterBlueprintResponse.Value("dimensions").Array().Element(1).Object().Value("options").Array().Length().Equal(38)
				filterBlueprintResponse.Value("dimensions").Array().Element(2).Object().Value("name").Equal("time")
				filterBlueprintResponse.Value("dimensions").Array().Element(2).Object().Value("options").Array().Length().Equal(1)
				filterBlueprintResponse.Value("links").Object().Value("dimensions").Object().Value("href").String().Match("/filters/" + filterBlueprintID + "/dimensions$")
				filterBlueprintResponse.Value("links").Object().Value("self").Object().Value("href").String().Match("/filters/(.+)$")
				filterBlueprintResponse.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetName + "/editions/2017/versions/1$")
				filterBlueprintResponse.Value("links").Object().Value("version").Object().Value("id").Equal("1")
				filterBlueprintResponse.Value("links").Object().Value("filter_output").Object().Value("href").String().Match("/filter-outputs/(.+)$")
				filterBlueprintResponse.Value("links").Object().Value("filter_output").Object().Value("id").NotNull()

				filterOutputID := filterBlueprintResponse.Value("links").Object().Value("filter_output").Object().Value("id").String().Raw()

				filterOutputResource, err := mongo.GetFilter("filters", "filterOutputs", "filter_id", filterOutputID)
				if err != nil {
					log.ErrorC("Unable to retrieve filter output document", err, log.Data{"filter_output_id": filterOutputID})
					os.Exit(1)
				}

				So(filterOutputResource.FilterID, ShouldEqual, filterOutputID)
				So(filterOutputResource.InstanceID, ShouldEqual, instanceID)
				So(filterOutputResource.State, ShouldEqual, "created")

				tryAgain := true
				for tryAgain {
					filterOutputResource, err = mongo.GetFilter("filters", "filterOutputs", "filter_id", filterOutputID)
					if err != nil {
						log.ErrorC("Unable to retrieve filter output document", err, log.Data{"filter_output_id": filterOutputID})
						os.Exit(1)
					}
					if filterOutputResource.State == "completed" {
						tryAgain = false
					}
				}

				So(filterOutputResource.FilterID, ShouldEqual, filterOutputID)
				So(filterOutputResource.InstanceID, ShouldEqual, instanceID)
				So(filterOutputResource.State, ShouldEqual, "completed")
				So(filterOutputResource.Downloads.CSV, ShouldNotBeNil)
				So(filterOutputResource.Downloads.XLS, ShouldNotBeNil)

				var locationCSV, locationXLS string
				var filename []string
				if filterOutputResource.Downloads.CSV != nil {
					locationCSV = filterOutputResource.Downloads.CSV.URL
					filename = strings.Split(locationCSV, "https://csv-exported.s3.eu-west-1.amazonaws.com/")
				}

				if filterOutputResource.Downloads.XLS != nil {
					locationXLS = filterOutputResource.Downloads.XLS.URL
				}

				log.Info("My filtered files on aws", log.Data{"csv_location": locationCSV, "xls_location": locationXLS, "filename": filename[0]})

				// TODO get csv file and xlsx file
				if err = getS3File("eu-west-1", "csv-exported", filename[0]); err != nil {
					//log.ErrorC("failed to find filtered csv file", err, log.Data{"filename": filename + ".csv"})
				}

				// remove filter blueprint
				if err = mongo.Teardown("filters", "filters", "filter_id", filterBlueprintID); err != nil {
					if err != mgo.ErrNotFound {
						log.ErrorC("failed to remove filter blueprint resource", err, log.Data{"filter_blueprint_id": filterBlueprintID})
						hasRemovedAllResources = false
					}
				}

				// remove filter output
				if err = mongo.Teardown("filters", "filterOutputs", "filter_id", filterOutputID); err != nil {
					if err != mgo.ErrNotFound {
						log.ErrorC("failed to remove filter output resource", err, log.Data{"filter_output_id": filterOutputID})
						hasRemovedAllResources = false
					}
				}

				if err = deleteS3File("eu-west-1", "csv-exported", locationCSV); err != nil {
					log.ErrorC("Failed to remove filtered csv file from s3", err, log.Data{"location": locationCSV})
					hasRemovedAllResources = false
				}

				if err = deleteS3File("eu-west-1", "csv-exported", locationXLS); err != nil {
					log.ErrorC("Failed to remove filtered xls file from s3", err, log.Data{"location": locationXLS})
					hasRemovedAllResources = false
				}
			})

			// delete dataset
			if err = mongo.Teardown("datasets", "datasets", "_id", datasetName); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove dataset resource", err, log.Data{"dataset_id": datasetName})
					hasRemovedAllResources = false
				}
			}

			// remove job
			if err = mongo.Teardown("imports", "imports", "id", jobID); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove job resource", err, log.Data{"job_id": jobID})
					hasRemovedAllResources = false
				}
			}

			// remove instance/versions
			if err = mongo.Teardown("datasets", "instances", "id", instanceID); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove instance resource", err, log.Data{"instance_id": instanceID})
					hasRemovedAllResources = false
				}
			}

			// remove dimension options
			if err = mongo.Teardown("datasets", "dimension.options", "instance_id", instanceID); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove dimension option resources", err, log.Data{"instance_id": instanceID})
					hasRemovedAllResources = false
				}
			}

			// remove edition
			if err = mongo.Teardown("datasets", "editions", "links.self.href", instanceResource.Links.Edition.HRef); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove edition resource", err, log.Data{"links.self.href": instanceResource.Links.Edition.HRef})
					hasRemovedAllResources = false
				}
			}

			// remove instance from neo4j
			datastore, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, "")
			if err != nil {
				log.ErrorC("Failed to connecton to neo4j database", err, nil)
				os.Exit(1)

			}

			if err = datastore.TeardownInstance(); err != nil {
				log.ErrorC("Failed to delete all instances in neo4j database", err, nil)
				hasRemovedAllResources = false
			}
		})

		// remove test file from s3
		if err := deleteS3File("eu-west-1", "ons-dp-cmd-test", "v4TestFile.csv"); err != nil {
			log.ErrorC("Failed to remove test file from s3", err, nil)
			hasRemovedAllResources = false
		} else {
			log.Info("successfully removed file from aws", nil)
		}

		if !hasRemovedAllResources {
			os.Exit(1)
		}
	})
}
