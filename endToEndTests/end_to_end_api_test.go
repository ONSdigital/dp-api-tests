package generateFiles

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

var timeout = time.Duration(30 * time.Second)

// TODO Once export services have been updated with encryption and decryption
// remove decrypt boolean flag from all setup functions
func TestSuccessfulEndToEndProcess(t *testing.T) {

	importAPI := httpexpect.New(t, cfg.ImportAPIURL)
	recipeAPI := httpexpect.New(t, cfg.RecipeAPIURL)
	datasetAPI := httpexpect.New(t, cfg.DatasetAPIURL)
	filterAPI := httpexpect.New(t, cfg.FilterAPIURL)
	hierarchyAPI := httpexpect.New(t, cfg.HierarchyAPIURL)
	searchAPI := httpexpect.New(t, cfg.SearchAPIURL)

	hasRemovedAllResources := true
	filename := "v4TestFile.csv"
	recipe := "2943f3c5-c3f1-4a9a-aa6e-14d21c33524c"

	// Get dataset ID from recipe API
	recipeResponse := recipeAPI.GET("/recipes/{recipe}", recipe).
		Expect().Status(http.StatusOK).JSON().Object()

	recipeResponse.Value("id").NotNull()

	Convey("Given a v4 file exists in aws", t, func() {
		// Send v4 file to aws
		err := sendV4FileToAWS(region, bucketName, filename, true)
		if err != nil {
			log.ErrorC("failed to load in v4 to aws, discontinue with test", err, nil)
			os.Exit(1)
		}

		// import API expects a s3 url as the location of the file
		location := "s3://" + bucketName + "/" + filename

		Convey("When a job is imported and the version of the dataset is published", func() {

			// STEP 1 - Create job with state created
			postJobResponse := importAPI.POST("/jobs").WithBytes([]byte(createValidJobJSON(recipe, location))).
				WithHeader(internalTokenHeader, importAPIInternalToken).Expect().Status(http.StatusCreated).JSON().Object()

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
			instanceResource, err := mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
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
			importAPI.PUT("/jobs/{id}", jobID).WithHeader(internalTokenHeader, importAPIInternalToken).WithBytes([]byte(`{"state":"submitted"}`)).
				Expect().Status(http.StatusOK)

			// Check import job state is completed or submitted
			jobResource, err := mongo.GetJob(cfg.MongoImportsDB, "imports", "id", jobID)
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
			totalObservations := int64(1510)

			tryAgain := true

			exitObservationsCompleteLoop := make(chan bool)

			go func() {
				time.Sleep(timeout)
				close(exitObservationsCompleteLoop)
			}()

		observationsCompleteLoop:
			for tryAgain {
				select {
				case <-exitObservationsCompleteLoop:
					break observationsCompleteLoop
				default:
					instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
					if err != nil {
						log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
						os.Exit(1)
					}
					if instanceResource.ImportTasks.ImportObservations.State == "completed" {
						tryAgain = false
					} else {
						So(instanceResource.State, ShouldEqual, "submitted")
						So(instanceResource.ImportTasks.ImportObservations.State, ShouldEqual, "created")
						time.Sleep(time.Millisecond * 100) // Relax the continuous battering of mongo database
					}
				}
			}

			if tryAgain != false {
				err = errors.New("timed out")
				log.ErrorC("Timed out - failed to get instance document to a state of completed", err, log.Data{"instance_id": instanceID, "state": instanceResource.State, "timeout": timeout})
				os.Exit(1)
			}

			So(instanceResource.Headers, ShouldResemble, &[]string{"V4_0", "time", "time", "uk-only", "geography", "cpih1dim1aggid", "aggregate"})
			So(instanceResource.State, ShouldEqual, "submitted")
			So(instanceResource.ImportTasks.ImportObservations.State, ShouldEqual, "completed")
			So(instanceResource.ImportTasks.ImportObservations.InsertedObservations, ShouldResemble, totalObservations)
			So(instanceResource.TotalObservations, ShouldResemble, totalObservations)

			// Check dimension options
			count, err := mongo.CountDimensionOptions(cfg.MongoDB, "dimension.options", "instance_id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve dimension option resources", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(count, ShouldEqual, 156)

			// Check hierarchies have been built
			tryAgain = true

			exitHierarchiesCompleteLoop := make(chan bool)

			go func() {
				time.Sleep(timeout)
				close(exitHierarchiesCompleteLoop)
			}()

		hierarchiesCompleteLoop:
			for tryAgain {
				select {
				case <-exitHierarchiesCompleteLoop:
					break hierarchiesCompleteLoop
				default:
					instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
					if err != nil {
						log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
						os.Exit(1)
					}

					if instanceResource.ImportTasks.BuildHierarchyTasks == nil ||
						len(instanceResource.ImportTasks.BuildHierarchyTasks) < 1 {

						log.ErrorC("no build hierarchy tasks found", err, log.Data{"instance_id": instanceID})
						os.Exit(1)
					}

					if instanceResource.ImportTasks.BuildHierarchyTasks[0].State == "completed" {
						tryAgain = false
					} else {
						So(instanceResource.State, ShouldEqual, "submitted")
						So(instanceResource.ImportTasks.BuildHierarchyTasks[0].State, ShouldEqual, "created")
						time.Sleep(time.Millisecond * 100)
					}
				}
			}

			if tryAgain != false {
				err = errors.New("timed out")
				log.ErrorC("Timed out - failed to get instance document to have hierarchy tasks with states of completed", err, log.Data{"instance_id": instanceID, "hierarchy_tasks": instanceResource.ImportTasks.BuildHierarchyTasks, "timeout": timeout})
				os.Exit(1)
			}

			// Check hierarchies exist by calling the hierarchy api
			getHierarchyParentDimensionResponse := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension}", instanceID, "aggregate").WithHeader(internalTokenHeader, internalTokenID).
				Expect().Status(http.StatusOK).JSON().Object()

			getHierarchyParentDimensionResponse.Value("has_data").Equal(true)
			getHierarchyParentDimensionResponse.Value("label").Equal("Overall Index")
			getHierarchyParentDimensionResponse.Value("no_of_children").Equal(12)
			getHierarchyParentDimensionResponse.Value("links").Object().Value("code").Object().Value("href").Equal("http://localhost:22400/code-list/cpih1dim1aggid/code/cpih1dim1A0")
			getHierarchyParentDimensionResponse.Value("links").Object().Value("code").Object().Value("id").Equal("cpih1dim1A0")
			getHierarchyParentDimensionResponse.Value("links").Object().Value("self").Object().Value("href").Equal("http://localhost:22600/hierarchies/" + instanceID + "/aggregate")
			getHierarchyParentDimensionResponse.Value("links").Object().Value("self").Object().Value("id").Equal("cpih1dim1A0")

			numberOfChildren := getHierarchyParentDimensionResponse.Value("no_of_children").Raw()
			getHierarchyParentDimensionResponse.Value("children").Array().Length().Equal(numberOfChildren)

			// Reset tryAgain for next loop
			tryAgain = true

			exitElasticSearchCompleteLoop := make(chan bool)

			go func() {
				time.Sleep(timeout)
				close(exitElasticSearchCompleteLoop)
			}()

			// Check elastic search tasks have completed against instance
		elasticSearchCompleteLoop:
			for tryAgain {
				select {
				case <-exitElasticSearchCompleteLoop:
					break elasticSearchCompleteLoop
				default:
					instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
					if err != nil {
						log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
						os.Exit(1)
					}

					if instanceResource.ImportTasks.SearchTasks == nil ||
						len(instanceResource.ImportTasks.SearchTasks) < 1 {

						log.ErrorC("no build hierarchy tasks found", err, log.Data{"instance_id": instanceID})
						os.Exit(1)
					}

					if instanceResource.ImportTasks.SearchTasks[0].State == "completed" {
						tryAgain = false
					} else {
						So(instanceResource.State, ShouldEqual, "submitted")
						So(instanceResource.ImportTasks.SearchTasks[0].State, ShouldEqual, "created")
						time.Sleep(time.Millisecond * 100)
					}
				}
			}

			if tryAgain != false {
				err = errors.New("timed out")
				log.ErrorC("Timed out - failed to get instance document to have search tasks with states of completed", err, log.Data{"instance_id": instanceID, "search_tasks": instanceResource.ImportTasks.SearchTasks, "timeout": timeout})
				os.Exit(1)
			}

			// todo, add retry loop to check when the instance is set to complete.
			time.Sleep(time.Second * 5)

			// get the instance again now the tracker has had change to set the instance status to complete
			instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			// Check instance state is completed
			So(instanceResource.State, ShouldEqual, "completed")
			So(instanceResource.ImportTasks.SearchTasks[0].DimensionName, ShouldEqual, "aggregate")

			// STEP 4 - Update instance with meta data and change state to `edition-confirmed`
			datasetAPI.PUT("/instances/{instance_id}", instanceID).WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(validPUTInstanceMetadataJSON)).Expect().Status(http.StatusOK)

			// Check instance has updated
			instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
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
			editionResource, err := mongo.GetEdition(cfg.MongoDB, "editions", "links.self.href", instanceResource.Links.Edition.HRef)
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

			versionResource, err := mongo.GetVersion(cfg.MongoDB, "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(versionResource.CollectionID, ShouldEqual, "308064B3-A808-449B-9041-EA3A2F72CFAC")
			So(versionResource.State, ShouldEqual, "associated")

			// Check dataset has updated
			datasetResource, err := mongo.GetDataset(cfg.MongoDB, "datasets", "_id", datasetName)
			if err != nil {
				log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
				os.Exit(1)
			}

			So(datasetResource.Current, ShouldBeNil)
			So(datasetResource.Next.CollectionID, ShouldEqual, "308064B3-A808-449B-9041-EA3A2F72CFAC")
			So(datasetResource.Next.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
			So(datasetResource.Next.Links.LatestVersion.ID, ShouldEqual, "1")
			So(datasetResource.Next.State, ShouldEqual, "associated")

			exitHasDownloadsLoop := make(chan bool)

			go func() {
				time.Sleep(timeout)
				close(exitHasDownloadsLoop)
			}()

			// Waiting for version to have downloads before updating state to published
			hasDownloads := false
			var XLSSize int
		hasDownloadsLoop:
			for !hasDownloads {
				select {
				case <-exitHasDownloadsLoop:
					break hasDownloadsLoop
				default:

					instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
					if err != nil {
						log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
						os.Exit(1)
					}
					if instanceResource.Downloads != nil {
						if instanceResource.Downloads.XLS != nil {
							if instanceResource.Downloads.XLS.URL != "" {
								XLSSize, err = strconv.Atoi(instanceResource.Downloads.XLS.Size)
								if err != nil {
									log.ErrorC("cannot convert xls size of type string to integer", err, log.Data{"xls_size": instanceResource.Downloads.XLS.Size})
									os.Exit(1)
								}
								So(XLSSize, ShouldBeBetweenOrEqual, 19000, 20000)
								So(instanceResource.Downloads.XLS.URL, ShouldNotBeEmpty)
								CSVSize, err := strconv.Atoi(instanceResource.Downloads.CSV.Size)
								if err != nil {
									log.ErrorC("cannot convert csv size of type string to integer", err, log.Data{"csv_size": instanceResource.Downloads.CSV.Size})
									os.Exit(1)
								}
								So(CSVSize, ShouldBeBetweenOrEqual, 137000, 139000)
								So(instanceResource.Downloads.CSV.URL, ShouldNotBeEmpty)
								hasDownloads = true
							}
						}
					} else {
						So(instanceResource.State, ShouldEqual, "associated")
					}
				}
			}

			if hasDownloads == false {
				err := errors.New("timed out")
				log.ErrorC("Timed out - failed to get instance document with available downloads", err, log.Data{"instance_id": instanceID, "downloads": instanceResource.Downloads, "timeout": timeout})
				os.Exit(1)
			}

			// STEP 6 -  Update version to a state of published
			datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetName, "2017", "1").WithHeader(internalTokenHeader, internalTokenID).
				WithBytes([]byte(`{"state":"published"}`)).Expect().Status(http.StatusOK)

			versionResource, err = mongo.GetVersion(cfg.MongoDB, "instances", "id", instanceID)
			if err != nil {
				log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
				os.Exit(1)
			}

			So(versionResource.State, ShouldEqual, "published")

			// Check edition has updated
			editionResource, err = mongo.GetEdition(cfg.MongoDB, "editions", "links.self.href", instanceResource.Links.Edition.HRef)
			if err != nil {
				log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
				os.Exit(1)
			}

			So(editionResource.State, ShouldEqual, "published")

			// Check dataset has updated
			datasetResource, err = mongo.GetDataset(cfg.MongoDB, "datasets", "_id", datasetName)
			if err != nil {
				log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
				os.Exit(1)
			}

			So(datasetResource.Current, ShouldNotBeNil)
			So(datasetResource.Current.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
			So(datasetResource.Current.State, ShouldEqual, "published")

			// Check data exists in elaticsearch by calling search API to find dimension option
			getSearchResponse := searchAPI.GET("/search/datasets/{id}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetName, versionResource.Edition, strconv.Itoa(versionResource.Version), "aggregate").
				WithQuery("q", "Overall Index").Expect().Status(http.StatusOK).JSON().Object()

			getSearchResponse.Value("count").Equal(1)
			getSearchResponse.Value("items").Array().Length().Equal(1)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("code").Equal("cpih1dim1A0")
			getSearchResponse.Value("items").Array().Element(0).Object().Value("dimension_option_url").Equal("http://localhost:22400/code-list/cpih1dim1aggid/code/cpih1dim1A0")
			getSearchResponse.Value("items").Array().Element(0).Object().Value("has_data").Equal(true)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("label").Equal("Overall Index")
			getSearchResponse.Value("items").Array().Element(0).Object().Value("matches").Object().NotContainsKey("code")
			getSearchResponse.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Length().Equal(2)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("start").Equal(1)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(0).Object().Value("end").Equal(7)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("start").Equal(9)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("matches").Object().Value("label").Array().Element(1).Object().Value("end").Equal(13)
			getSearchResponse.Value("items").Array().Element(0).Object().Value("number_of_children").Equal(12)
			getSearchResponse.Value("limit").Equal(20)
			getSearchResponse.Value("offset").Equal(0)

			Convey("Then an api customer should be able to get a csv and xls download link", func() {
				// Get downloads link from version document
				bucket := "csv-exported"
				csvURL := versionResource.Downloads.CSV.URL
				csvFilename := strings.TrimPrefix(csvURL, "https://"+bucket+".s3."+region+".amazonaws.com/")

				// read csv download from s3
				csvFile, err := getS3File(region, bucket, csvFilename, false)
				if err != nil {
					log.ErrorC("unable to find csv full download in s3", err, log.Data{"csv_url": csvURL, "csv_filename": csvFilename})
					os.Exit(1)
				}
				defer csvFile.Close()

				csvReader := csv.NewReader(csvFile)
				headerRow, err := csvReader.Read()
				if err != nil {
					log.ErrorC("unable to read header row", err, log.Data{"csv_url": csvURL, "csv_filename": csvFilename})
				}

				So(len(headerRow), ShouldEqual, 7)

				// check the number of rows and anything else (e.g. meta data)
				numberOfCSVRows := 0
				for {
					_, err = csvReader.Read()
					if err != nil {
						if err == io.EOF {
							break
						}
						log.ErrorC("unable to read row", err, log.Data{"csv_url": csvURL, "csv_filename": csvFilename})
						os.Exit(1)
					}
					numberOfCSVRows++
				}
				So(numberOfCSVRows, ShouldEqual, 1510)

				xlsURL := versionResource.Downloads.XLS.URL
				xlsFilename := strings.TrimPrefix(xlsURL, "https://"+bucket+".s3-"+region+".amazonaws.com/")

				// read xls download from s3
				xlsFile, err := getS3File(region, bucket, xlsFilename, false)
				if err != nil {
					log.ErrorC("unable to find xls full download in s3", err, log.Data{"xls_url": xlsURL, "csv_filename": xlsFilename})
					os.Exit(1)
				}
				defer xlsFile.Close()

				So(xlsFile, ShouldNotBeEmpty)

				xlsFileSize, err := getS3FileSize(region, bucket, xlsFilename, false)
				if err != nil {
					log.ErrorC("unable to extract size of xls full download in s3", err, log.Data{"xls_url": xlsURL, "csv_filename": xlsFilename})
					os.Exit(1)
				}

				expectedXLSFileSize := int64(XLSSize)
				So(xlsFileSize, ShouldResemble, &expectedXLSFileSize)

				Convey("Then an api customer should be able to filter a dataset and be able to download a csv and xlsx download of the data", func() {
					json := GetValidPOSTCreateFilterJSON(datasetName, "2017", "1")
					log.Info("ajhgjlabfjlarebvjkrbvqlj", log.Data{"json": json})
					filterBlueprintResponse := filterAPI.POST("/filters").WithQuery("submitted", "true").
						WithBytes([]byte(json)).Expect().Status(http.StatusCreated).JSON().Object()

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

					log.Info("filter response", log.Data{"resp": filterBlueprintResponse.Raw()})

					filterOutputID := filterBlueprintResponse.Value("links").Object().Value("filter_output").Object().Value("id").String().Raw()

					filterOutputResource, err := mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
					if err != nil {
						log.ErrorC("Unable to retrieve filter output document", err, log.Data{"filter_output_id": filterOutputID})
						os.Exit(1)
					}

					So(filterOutputResource.FilterID, ShouldEqual, filterOutputID)
					So(filterOutputResource.InstanceID, ShouldEqual, instanceID)
					So(filterOutputResource.State, ShouldEqual, "created")

					filterOutputResourceCompleted := true
					for filterOutputResourceCompleted {
						filterOutputResource, err = mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
						if err != nil {
							log.ErrorC("Unable to retrieve filter output document", err, log.Data{"filter_output_id": filterOutputID})
							os.Exit(1)
						}
						if filterOutputResource.State == "completed" {
							filterOutputResourceCompleted = false
						}
					}

					So(filterOutputResource.FilterID, ShouldEqual, filterOutputID)
					So(filterOutputResource.InstanceID, ShouldEqual, instanceID)
					So(filterOutputResource.State, ShouldEqual, "completed")
					So(filterOutputResource.Downloads.CSV, ShouldNotBeNil)
					So(filterOutputResource.Downloads.XLS, ShouldNotBeNil)

					filteredCSVURL := filterOutputResource.Downloads.CSV.URL
					filteredCSVFilename := strings.TrimPrefix(filteredCSVURL, "https://"+bucket+".s3."+region+".amazonaws.com/")

					// read csv download from s3
					filteredCSVFile, err := getS3File(region, bucket, filteredCSVFilename, false)
					if err != nil {
						log.ErrorC("unable to find filtered csv download in s3", err, log.Data{"filtered_csv_url": filteredCSVURL, "filtered_csv_filename": filteredCSVFilename})
						os.Exit(1)
					}

					filteredCSVReader := csv.NewReader(filteredCSVFile)

					if err = checkFileRowCount(filteredCSVReader, 39); err != nil {
						log.ErrorC("unable to check file row count", err, nil)
						os.Exit(1)
					}

					filteredXLSURL := filterOutputResource.Downloads.XLS.URL
					filteredXLSFilename := strings.TrimPrefix(filteredXLSURL, "https://"+bucket+".s3-"+region+".amazonaws.com/")

					// read xls download from s3
					filteredXLSFile, err := getS3File(region, bucket, filteredXLSFilename, false)
					if err != nil {
						log.ErrorC("unable to find filtered xls download in s3", err, log.Data{"filtered_xls_url": filteredXLSURL, "filtered_xls_filename": filteredXLSFilename})
						os.Exit(1)
					}

					So(filteredXLSFile, ShouldNotBeEmpty)

					filteredXLSFileSize, err := getS3FileSize(region, bucket, filteredXLSFilename, false)
					if err != nil {
						log.ErrorC("unable to extract size of filtered xls download in s3", err, log.Data{"filtered_xls_url": filteredXLSURL, "filtered_xls_filename": filteredXLSFilename})
						os.Exit(1)
					}

					minExpectedXLSFileSize := int64(7432)
					maxExpectedXLSFileSize := int64(7436)
					So(*filteredXLSFileSize, ShouldBeBetweenOrEqual, minExpectedXLSFileSize, maxExpectedXLSFileSize)

					filterBlueprint := &mongo.Doc{
						Database:   cfg.MongoFiltersDB,
						Collection: "filters",
						Key:        "filter_id",
						Value:      filterBlueprintID,
					}

					filterOutput := &mongo.Doc{
						Database:   cfg.MongoFiltersDB,
						Collection: "filterOutputs",
						Key:        "filter_id",
						Value:      filterOutputID,
					}

					// remove filter output
					if err = mongo.Teardown(filterBlueprint, filterOutput); err != nil {
						if err != mgo.ErrNotFound {
							log.ErrorC("failed to remove filter output resource", err, log.Data{"filter_output_id": filterOutputID})
							hasRemovedAllResources = false
						}
					}

					if err = deleteS3File("eu-west-1", "csv-exported", filteredCSVFilename); err != nil {
						log.ErrorC("Failed to remove filtered csv file from s3", err, log.Data{"filename": filteredCSVFilename})
						hasRemovedAllResources = false
					}

					if err = deleteS3File("eu-west-1", "csv-exported", filteredXLSFilename); err != nil {
						log.ErrorC("Failed to remove filtered xls file from s3", err, log.Data{"filename": filteredXLSFilename})
						hasRemovedAllResources = false
					}
				})

				// remove test file from s3
				if err := deleteS3File(region, bucket, csvFilename); err != nil {
					log.ErrorC("Failed to remove full downloadable test csv file from s3", err, nil)
					hasRemovedAllResources = false
				}

				// remove test file from s3
				if err := deleteS3File(region, bucket, xlsFilename); err != nil {
					log.ErrorC("Failed to remove full downloadable test xls file from s3", err, nil)
					hasRemovedAllResources = false
				}
			})

			var docs []*mongo.Doc

			dataset := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "datasets",
				Key:        "_id",
				Value:      datasetName,
			}

			importJob := &mongo.Doc{
				Database:   cfg.MongoImportsDB,
				Collection: "imports",
				Key:        "id",
				Value:      jobID,
			}

			instance := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "instances",
				Key:        "id",
				Value:      instanceID,
			}

			edition := &mongo.Doc{
				Database:   cfg.MongoDB,
				Collection: "editions",
				Key:        "links.self.href",
				Value:      instanceResource.Links.Edition.HRef,
			}

			docs = append(docs, dataset, importJob, instance, edition)

			// remove all mongo documents created in the test
			if err = mongo.Teardown(docs...); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove edition resource", err, log.Data{"links.self.href": instanceResource.Links.Edition.HRef})
					hasRemovedAllResources = false
				}
			}

			// remove all dimension options from mongo collection
			if err = mongo.TeardownAll(cfg.MongoDB, "dimension.options"); err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to remove edition resource", err, log.Data{"links.self.href": instanceResource.Links.Edition.HRef})
					hasRemovedAllResources = false
				}
			}

			// remove elasticsearch index for instance and dimension (call elasticsearch directly)
			if _, err = elasticsearch.DeleteIndex(cfg.ElasticSearchAPIURL + "/" + instanceID + "_aggregate"); err != nil {
				log.ErrorC("Failed to delete index from elasticsearch", err, nil)
				hasRemovedAllResources = false
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

func checkFileRowCount(csvReader *csv.Reader, expectedCount int64) error {
	numberOfRows := int64(0)
	// Iterate over file counting the number of rows that exist
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.ErrorC("encountered error reading csv", err, log.Data{"csv_line": line})
			return err
		}
		numberOfRows++
	}

	So(numberOfRows, ShouldEqual, expectedCount)

	return nil
}
