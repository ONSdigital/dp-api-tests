package generateFiles

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"

	"net/url"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/rchttp"
)

var timeout = time.Duration(30 * time.Second)

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
			t.FailNow()
		}

		// import API expects a s3 url as the location of the file
		location := "s3://" + bucketName + "/" + filename

		log.Info("Create job with state created", nil)
		postJobResponse := importAPI.POST("/jobs").WithBytes([]byte(createValidJobJSON(recipe, location))).
			WithHeaders(headers).Expect().Status(http.StatusCreated).JSON().Object()

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
			t.FailNow()
		}

		So(len(instanceResource.Dimensions), ShouldEqual, 3)
		So(instanceResource.Links.Job.ID, ShouldEqual, jobID)
		So(instanceResource.Links.Job.HRef, ShouldEqual, cfg.ImportAPIURL+"/jobs/"+jobID)
		So(instanceResource.Links.Dataset.ID, ShouldEqual, datasetName)
		So(instanceResource.Links.Dataset.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName)
		So(instanceResource.Links.Self.HRef, ShouldEqual, cfg.DatasetAPIURL+"/instances/"+instanceID)
		So(instanceResource.State, ShouldEqual, "created")

		log.Info("Create dataset with dataset id from previous response", nil)
		postDatasetResponse := datasetAPI.POST("/datasets/{id}", datasetName).WithHeaders(headers).WithBytes([]byte(validPOSTCreateDatasetJSON)).
			Expect().Status(http.StatusCreated).JSON().Object()

		postDatasetResponse.Value("next").Object().Value("links").Object().Value("self").Object().Value("href").String().Match(cfg.DatasetAPIURL + "/datasets/" + datasetName + "$")
		postDatasetResponse.Value("next").Object().Value("state").Equal("created")

		log.Info("Update job state to submitted", nil)
		importAPI.PUT("/jobs/{id}", jobID).WithHeaders(headers).WithBytes([]byte(`{"state":"submitted"}`)).
			Expect().Status(http.StatusOK)

		// Check import job state is completed or submitted
		jobResource, err := mongo.GetJob(cfg.MongoImportsDB, "imports", "id", jobID)
		if err != nil {
			log.ErrorC("Unable to retrieve job resource", err, log.Data{"job_id": jobID})
			t.FailNow()
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
					t.FailNow()
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
			t.FailNow()
		}

		So(instanceResource.Headers, ShouldResemble, &[]string{"V4_0", "mmm-yy", "time", "uk-only", "geography", "cpih1dim1aggid", "aggregate"})
		So(instanceResource.State, ShouldEqual, "submitted")
		So(instanceResource.ImportTasks.ImportObservations.State, ShouldEqual, "completed")
		So(instanceResource.ImportTasks.ImportObservations.InsertedObservations, ShouldResemble, totalObservations)
		So(instanceResource.TotalObservations, ShouldResemble, totalObservations)

		// Check dimension options
		count, err := mongo.CountDimensionOptions(cfg.MongoDB, "dimension.options", "instance_id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve dimension option resources", err, log.Data{"instance_id": instanceID})
			t.FailNow()
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
					t.FailNow()
				}

				if instanceResource.ImportTasks.BuildHierarchyTasks == nil ||
					len(instanceResource.ImportTasks.BuildHierarchyTasks) < 1 {

					log.ErrorC("no build hierarchy tasks found", err, log.Data{"instance_id": instanceID})
					t.FailNow()
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
			t.FailNow()
		}

		// Check hierarchies exist by calling the hierarchy api
		getHierarchyParentDimensionResponse := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension}", instanceID, "aggregate").WithHeaders(headers).
			Expect().Status(http.StatusOK).JSON().Object()

		getHierarchyParentDimensionResponse.Value("has_data").Equal(true)
		getHierarchyParentDimensionResponse.Value("label").Equal("Overall Index")
		getHierarchyParentDimensionResponse.Value("no_of_children").Equal(12)
		getHierarchyParentDimensionResponse.Value("links").Object().Value("code").Object().Value("href").Equal(cfg.CodeListAPIURL + "/code-lists/cpih1dim1aggid/codes/cpih1dim1A0")
		getHierarchyParentDimensionResponse.Value("links").Object().Value("code").Object().Value("id").Equal("cpih1dim1A0")
		getHierarchyParentDimensionResponse.Value("links").Object().Value("self").Object().Value("href").Equal(cfg.HierarchyAPIURL + "/hierarchies/" + instanceID + "/aggregate")

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
					t.FailNow()
				}

				if instanceResource.ImportTasks.SearchTasks == nil ||
					len(instanceResource.ImportTasks.SearchTasks) < 1 {

					log.ErrorC("no build hierarchy tasks found", err, log.Data{"instance_id": instanceID})
					t.FailNow()
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
			t.FailNow()
		}

		// todo, add retry loop to check when the instance is set to complete.
		time.Sleep(time.Second * 5)

		// get the instance again now the tracker has had change to set the instance status to complete
		instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
			t.FailNow()
		}

		// Check instance state is completed
		So(instanceResource.State, ShouldEqual, "completed")
		So(instanceResource.ImportTasks.SearchTasks[0].DimensionName, ShouldEqual, "aggregate")

		log.Info("Update instance with meta data and change state to `edition-confirmed`", nil)
		datasetAPI.PUT("/instances/{instance_id}", instanceID).WithHeaders(headers).
			WithBytes([]byte(validPUTInstanceMetadataJSON)).Expect().Status(http.StatusOK)

		// Check instance has updated
		instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve instance resource", err, log.Data{"instance_id": instanceID})
			t.FailNow()
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
		editionResource, err := mongo.GetEdition(cfg.MongoDB, "editions", "next.links.self.href", instanceResource.Links.Edition.HRef)
		if err != nil {
			log.ErrorC("Unable to retrieve edition resource", err, log.Data{"links.self.href": instanceResource.Links.Edition.HRef})
			t.FailNow()
		}

		So(editionResource.Next.Edition, ShouldEqual, "2017")
		So(editionResource.Next.Links.Dataset.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName)
		So(editionResource.Next.Links.Dataset.ID, ShouldEqual, datasetName)
		So(editionResource.Next.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
		So(editionResource.Next.Links.LatestVersion.ID, ShouldEqual, "1")
		So(editionResource.Next.Links.Self.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017")
		So(editionResource.Next.Links.Versions.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions")
		So(editionResource.Next.State, ShouldEqual, "edition-confirmed")

		log.Info("Get single observation from pre-published version", nil)
		observationsResource := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/observations", datasetName, "2017", "1").
			WithQueryString("time=Apr-05&geography=K02000001&aggregate=cpih1dim1G50100").
			WithHeaders(headers).
			Expect().Status(http.StatusOK).JSON().Object()

		observationsResource.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("href").String().Match("/code-lists/cpih1dim1aggid/codes/cpih1dim1G50100$")
		observationsResource.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("id").Equal("cpih1dim1G50100")
		observationsResource.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("href").String().Match("/code-lists/uk-only/codes/K02000001$")
		observationsResource.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("id").Equal("K02000001")
		observationsResource.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("href").String().Match("/code-lists/mmm-yy/codes/Apr-05$")
		observationsResource.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("id").Equal("Apr-05")
		observationsResource.Value("limit").Equal(10000)
		observationsResource.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetName + "/editions/2017/versions/1/metadata$")
		observationsResource.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetName + "/editions/2017/versions/1/observations\\?aggregate=cpih1dim1G50100&geography=K02000001&time=Apr-05$")
		observationsResource.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetName + "/editions/2017/versions/1$")
		observationsResource.Value("links").Object().Value("version").Object().Value("id").Equal("1")
		observationsResource.Value("observations").Array().Length().Equal(1)
		observationsResource.Value("observations").Array().Element(0).Object().Value("observation").Equal("81.7")
		observationsResource.Value("offset").Equal(0)
		observationsResource.Value("total_observations").Equal(1)
		observationsResource.Value("unit_of_measure").Equal("Pounds Sterling")

		log.Info("Update version with collection_id and change state to associated", nil)
		datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetName, "2017", "1").WithHeaders(headers).
			WithBytes([]byte(validPUTUpdateVersionToAssociatedJSON)).Expect().Status(http.StatusOK)

		versionResource, err := mongo.GetVersion(cfg.MongoDB, "instances", "id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
			t.FailNow()
		}

		So(versionResource.CollectionID, ShouldEqual, "308064B3-A808-449B-9041-EA3A2F72CFAC")
		So(versionResource.State, ShouldEqual, "associated")

		// Check dataset has updated
		datasetResource, err := mongo.GetDataset(cfg.MongoDB, "datasets", "_id", datasetName)
		if err != nil {
			log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
			t.FailNow()
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
		var csvSize int
		var csvwSize int
		var xlsSize int
	hasDownloadsLoop:
		for !hasDownloads {
			select {
			case <-exitHasDownloadsLoop:
				break hasDownloadsLoop
			default:

				instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
				if err != nil {
					log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
					t.FailNow()
				}

				if instanceResource.Downloads != nil {
					if instanceResource.Downloads.XLS != nil {
						if instanceResource.Downloads.XLS.Private != "" {
							xlsSize, err = strconv.Atoi(instanceResource.Downloads.XLS.Size)
							if err != nil {
								log.ErrorC("cannot convert xls size of type string to integer", err, log.Data{"xls_size": instanceResource.Downloads.XLS.Size})
								t.FailNow()
							}
							So(xlsSize, ShouldBeBetweenOrEqual, 19000, 20000)
							So(instanceResource.Downloads.XLS.Private, ShouldNotBeEmpty)

							csvSize, err = strconv.Atoi(instanceResource.Downloads.CSV.Size)
							if err != nil {
								log.ErrorC("cannot convert csv size of type string to integer", err, log.Data{"csv_size": instanceResource.Downloads.CSV.Size})
								t.FailNow()
							}
							So(csvSize, ShouldBeBetweenOrEqual, 137000, 139000)
							So(instanceResource.Downloads.CSV.URL, ShouldNotBeEmpty)

							csvwSize, err = strconv.Atoi(instanceResource.Downloads.CSVW.Size)
							if err != nil {
								log.ErrorC("cannot convert csvw size of type string to integer", err, log.Data{"csvw_size": instanceResource.Downloads.CSVW.Size})
								t.FailNow()
							}
							So(csvwSize, ShouldBeBetweenOrEqual, 1400, 1600)
							So(instanceResource.Downloads.CSVW.URL, ShouldNotBeEmpty)
							hasDownloads = true
						}
					}
				} else {
					So(instanceResource.State, ShouldEqual, "associated")
				}

				time.Sleep(time.Millisecond * 200)
			}
		}

		if hasDownloads == false {
			err = errors.New("timed out")
			log.ErrorC("Timed out - failed to get instance document with available downloads", err, log.Data{"instance_id": instanceID, "downloads": instanceResource.Downloads, "timeout": timeout})
			t.FailNow()
		}

		log.Info("attempting to read private full file download from S3", nil)
		instanceResource, err = mongo.GetInstance(cfg.MongoDB, "instances", "id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve instance document", err, log.Data{"instance_id": instanceID})
			t.FailNow()
		}

		logData := log.Data{
			"private_csv_link":  instanceResource.Downloads.CSV.Private,
			"private_csv_size":  instanceResource.Downloads.CSV.Size,
			"private_csvw_link": instanceResource.Downloads.CSVW.Private,
			"private_csvw_size": instanceResource.Downloads.CSVW.Size,
			"private_xls_link":  instanceResource.Downloads.XLS.Private,
			"private_xls_size":  instanceResource.Downloads.XLS.Size,
		}
		log.Debug("Pre publish full downloads have been generated", logData)

		privateCSVURL, err := url.Parse(instanceResource.Downloads.CSV.Private)
		if err != nil {
			log.ErrorC("failed to parse private CSV URL", err, log.Data{"url": instanceResource.Downloads.CSV.Private})
			t.FailNow()
		}

		privateCSVFilename := privateCSVURL.Path

		// read csv download from s3
		privateCSVFile, err := getS3File(region, bucket, privateCSVFilename, true)
		if err != nil {
			log.ErrorC("unable to find csv download in s3", err, nil)
			t.Error()
			t.FailNow()
		}
		privateCSVReader := csv.NewReader(privateCSVFile)
		if err = checkFileRowCount(privateCSVReader, 1511); err != nil {
			log.ErrorC("unable to check file row count", err, nil)
			t.FailNow()
		}

		log.Debug("Pre publish - attempting to download full CSV file from the download service", logData)
		testFileDownload(instanceResource.Downloads.CSV.URL, csvSize, false)

		log.Debug("Pre publish - attempting to download full CSVW file from the download service", logData)
		testFileDownload(instanceResource.Downloads.CSVW.URL, csvwSize, false)

		log.Debug("Pre publish - attempting to download full XLS file from the download service", logData)
		testFileDownload(instanceResource.Downloads.XLS.URL, xlsSize, false)

		log.Info("Then an authenticated user should be able to filter a dataset", nil)

		prePublishFilterBlueprintID, prePublishFilterOutputID := testFiltering(t, filterAPI, instanceID, false)

		log.Info("pre publish filter test passed", log.Data{
			"filter_blueprint": prePublishFilterBlueprintID,
			"filter_output":    prePublishFilterOutputID,
		})

		prePublishFilterBlueprint := &mongo.Doc{
			Database:   cfg.MongoFiltersDB,
			Collection: "filters",
			Key:        "filter_id",
			Value:      prePublishFilterBlueprintID,
		}

		prePublishFilterOutput := &mongo.Doc{
			Database:   cfg.MongoFiltersDB,
			Collection: "filterOutputs",
			Key:        "filter_id",
			Value:      prePublishFilterOutputID,
		}

		//remove filter output
		if err = mongo.Teardown(prePublishFilterBlueprint, prePublishFilterOutput); err != nil {
			if err != mgo.ErrNotFound {
				log.ErrorC("failed to remove filter output resource", err, log.Data{"filter_output_id": prePublishFilterOutputID})
				hasRemovedAllResources = false
			}
		}

		log.Info("STEP 6 - Update version to a state of published", nil)
		datasetAPI.PUT("/datasets/{id}/editions/{edition}/versions/{version}", datasetName, "2017", "1").WithHeaders(headers).
			WithBytes([]byte(`{"state":"published"}`)).Expect().Status(http.StatusOK)

		versionResource, err = mongo.GetVersion(cfg.MongoDB, "instances", "id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
			t.FailNow()
		}

		log.Info("Check edition has updated", nil)
		editionResource, err = mongo.GetEdition(cfg.MongoDB, "editions", "current.links.self.href", instanceResource.Links.Edition.HRef)
		if err != nil {
			log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
			t.FailNow()
		}

		So(editionResource.Next.State, ShouldEqual, "published")
		So(editionResource.Current, ShouldNotBeNil)
		So(editionResource.Current.State, ShouldEqual, "published")

		log.Info("Check dataset has updated", nil)
		datasetResource, err = mongo.GetDataset(cfg.MongoDB, "datasets", "_id", datasetName)
		if err != nil {
			log.ErrorC("Unable to retrieve dataset resource", err, log.Data{"dataset_id": datasetName})
			t.FailNow()
		}

		So(datasetResource.Current, ShouldNotBeNil)
		So(datasetResource.Current.Links.LatestVersion.HRef, ShouldEqual, cfg.DatasetAPIURL+"/datasets/"+datasetName+"/editions/2017/versions/1")
		So(datasetResource.Current.State, ShouldEqual, "published")

		log.Info("Check data exists in elaticsearch by calling search API to find dimension option", nil)
		getSearchResponse := searchAPI.GET("/search/datasets/{id}/editions/{edition}/versions/{version}/dimensions/{dimension}", datasetName, versionResource.Edition, strconv.Itoa(versionResource.Version), "aggregate").
			WithQuery("q", "Overall Index").WithHeader(authorizationTokenHeader, authorizationToken).Expect().Status(http.StatusOK).JSON().Object()

		getSearchResponse.Value("count").Equal(1)
		getSearchResponse.Value("items").Array().Length().Equal(1)
		getSearchResponse.Value("items").Array().Element(0).Object().Value("code").Equal("cpih1dim1A0")
		getSearchResponse.Value("items").Array().Element(0).Object().Value("dimension_option_url").Equal("http://localhost:22400/code-lists/cpih1dim1aggid/codes/cpih1dim1A0")
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

		log.Info("Get single observation post-published version", nil)
		postObservationsResource := datasetAPI.GET("/datasets/{id}/editions/{edition}/versions/{version}/observations", datasetName, "2017", "1").
			WithQueryString("time=Apr-05&geography=K02000001&aggregate=cpih1dim1G50100").WithHeader(authorizationTokenHeader, authorizationToken).
			Expect().Status(http.StatusOK).JSON().Object()

		postObservationsResource.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("href").String().Match("/code-lists/cpih1dim1aggid/codes/cpih1dim1G50100$")
		postObservationsResource.Value("dimensions").Object().Value("aggregate").Object().Value("option").Object().Value("id").Equal("cpih1dim1G50100")
		postObservationsResource.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("href").String().Match("/code-lists/uk-only/codes/K02000001$")
		postObservationsResource.Value("dimensions").Object().Value("geography").Object().Value("option").Object().Value("id").Equal("K02000001")
		postObservationsResource.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("href").String().Match("/code-lists/mmm-yy/codes/Apr-05$")
		postObservationsResource.Value("dimensions").Object().Value("time").Object().Value("option").Object().Value("id").Equal("Apr-05")
		postObservationsResource.Value("limit").Equal(10000)
		postObservationsResource.Value("links").Object().Value("dataset_metadata").Object().Value("href").String().Match("/datasets/" + datasetName + "/editions/2017/versions/1/metadata$")
		postObservationsResource.Value("links").Object().Value("self").Object().Value("href").String().Match(".+/datasets/" + datasetName + "/editions/2017/versions/1/observations\\?aggregate=cpih1dim1G50100&geography=K02000001&time=Apr-05$")
		postObservationsResource.Value("links").Object().Value("version").Object().Value("href").String().Match("/datasets/" + datasetName + "/editions/2017/versions/1$")
		postObservationsResource.Value("links").Object().Value("version").Object().Value("id").Equal("1")
		postObservationsResource.Value("observations").Array().Length().Equal(1)
		postObservationsResource.Value("observations").Array().Element(0).Object().Value("observation").Equal("81.7")
		postObservationsResource.Value("offset").Equal(0)
		postObservationsResource.Value("total_observations").Equal(1)
		postObservationsResource.Value("unit_of_measure").Equal("Pounds Sterling")

		versionResourcePostPublish, err := mongo.GetVersion(cfg.MongoDB, "instances", "id", instanceID)
		if err != nil {
			log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
			t.FailNow()
		}
		logData["after_loop_public_csv_link"] = versionResourcePostPublish.Downloads.CSV.Public
		logData["after_loop_public_csvw_link"] = versionResourcePostPublish.Downloads.CSVW.Public
		logData["after_loop_public_xls_link"] = versionResourcePostPublish.Downloads.XLS.Public
		log.Debug("Pre publish full downloads have been generated", logData)

		log.Info("Get downloads link from version document", nil)
		csvURL := versionResource.Downloads.CSV.URL

		log.Debug("getting downloads from the version",
			log.Data{
				"csv_link":  versionResource.Downloads.CSV.URL,
				"csvw_link": versionResource.Downloads.CSVW.URL,
				"xls_link":  versionResource.Downloads.XLS.URL,
			})

		var req *http.Request
		req, err = http.NewRequest("GET", csvURL, nil)
		req.Header.Set(authorizationTokenHeader, authorizationToken)

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error(err, nil)
		}
		log.Info("get csv response", log.Data{
			"response_status": response.StatusCode,
		})
		defer func() {
			if err = response.Body.Close(); err != nil {
				log.ErrorC("get downloads link body", err, logData)
			}
		}()

		csvReader := csv.NewReader(response.Body)

		headerRow, err := csvReader.Read()
		if err != nil {
			log.ErrorC("unable to read header row", err, log.Data{"csv_url": csvURL})
		}

		So(len(headerRow), ShouldEqual, 7)

		log.Info("check the number of rows and anything else (e.g. meta data)", nil)
		numberOfCSVRows := 0
		for {
			_, err = csvReader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.ErrorC("unable to read row", err, log.Data{"csv_url": csvURL})
				t.FailNow()
			}
			numberOfCSVRows++
		}
		So(numberOfCSVRows, ShouldEqual, 1510)

		testFileDownload(versionResource.Downloads.CSV.URL, csvSize, true)
		testFileDownload(versionResource.Downloads.CSVW.URL, csvwSize, true)
		testFileDownload(versionResource.Downloads.XLS.URL, xlsSize, true)

		exitHasPublicDownloadsLoop := make(chan bool)

		go func() {
			time.Sleep(timeout)
			close(exitHasPublicDownloadsLoop)
		}()

		// Waiting for version to have downloads before updating state to published
	hasPublicDownloadsLoop:
		for {
			select {
			case <-exitHasPublicDownloadsLoop:
				err = errors.New("timed out")
				log.ErrorC("timeout waiting for public full download links to be generated", err, nil)
				t.FailNow()
			default:
				var versionResourcePostPublish mongo.Version
				versionResourcePostPublish, err = mongo.GetVersion(cfg.MongoDB, "instances", "id", instanceID)
				if err != nil {
					log.ErrorC("Unable to retrieve version resource", err, log.Data{"instance_id": instanceID})
					t.FailNow()
				}

				if versionResourcePostPublish.Downloads.XLS.Public != "" {
					break hasPublicDownloadsLoop
				}
			}
		}

		log.Info("Then an api customer should be able to filter a dataset and be able to download", nil)

		filterBlueprintID, filterOutputID := testFiltering(t, filterAPI, instanceID, true)

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
			Key:        "current.links.self.href",
			Value:      instanceResource.Links.Edition.HRef,
		}

		docs = append(docs, dataset, importJob, instance, edition)

		log.Debug("tearing down", nil)

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
			t.FailNow()
		}

		if err = datastore.TeardownInstance(); err != nil {
			log.ErrorC("Failed to delete all instances in neo4j database", err, nil)
			hasRemovedAllResources = false
		}

		// remove test file from s3
		if err := deleteS3File("eu-west-1", "ons-dp-cmd-test", "v4TestFile.csv"); err != nil {
			log.ErrorC("Failed to remove test file from s3", err, nil)
			hasRemovedAllResources = false
		} else {
			log.Info("successfully removed file from aws", nil)
		}

		if !hasRemovedAllResources {
			t.FailNow()
		}
	})
}

func testFileDownload(url string, expectedSize int, isPublished bool) {

	log.Info("attempting to download file from the download service", log.Data{
		"url":           url,
		"expected_size": expectedSize,
		"is_published":  isPublished,
	})

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err, nil)
	}

	req.Header.Add(authorizationTokenHeader, authorizationToken)

	response, err := rchttp.DefaultClient.Do(context.Background(), req)
	if err != nil {
		log.Error(err, nil)
	}

	So(response.StatusCode, ShouldEqual, http.StatusOK)

	log.Info("got response from the download service", log.Data{
		"response_status": response.StatusCode,
	})

	file := response.Body
	So(file, ShouldNotBeEmpty)

	b, _ := ioutil.ReadAll(file)

	actualFileSize := len(b)
	So(actualFileSize, ShouldEqual, expectedSize)
}

func testFiltering(t *testing.T, filterAPI *httpexpect.Expect, instanceID string, isPublished bool) (string, string) {

	json := GetValidPOSTCreateFilterJSON(datasetName, "2017", "1")

	log.Info("creating filter", log.Data{"json": json})

	filterBlueprintRequest := filterAPI.POST("/filters").
		WithQuery("submitted", "true").
		WithHeader(authorizationTokenHeader, authorizationToken).
		WithBytes([]byte(json))

	if !isPublished {
		filterBlueprintRequest = filterBlueprintRequest.WithHeaders(headers)
	}

	filterBlueprintResponse := filterBlueprintRequest.
		Expect().Status(http.StatusCreated).
		JSON().Object()

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
		t.FailNow()
	}
	So(filterOutputResource.FilterID, ShouldEqual, filterOutputID)
	So(filterOutputResource.InstanceID, ShouldEqual, instanceID)
	So(filterOutputResource.State, ShouldEqual, "created")

	log.Info("waiting for filter to be set to complete", nil)
	for i := 0; i < 10; i++ {

		filterOutputResource, err = mongo.GetFilter(cfg.MongoFiltersDB, "filterOutputs", "filter_id", filterOutputID)
		if err != nil {
			log.ErrorC("Unable to retrieve filter output document", err, log.Data{"filter_output_id": filterOutputID})
			t.FailNow()
		}
		if filterOutputResource.State == "completed" {
			break
		}

		time.Sleep(time.Millisecond * 200)
	}

	So(filterOutputResource.FilterID, ShouldEqual, filterOutputID)
	So(filterOutputResource.InstanceID, ShouldEqual, instanceID)
	So(filterOutputResource.State, ShouldEqual, "completed")
	So(filterOutputResource.Downloads.CSV, ShouldNotBeNil)
	So(filterOutputResource.Downloads.XLS, ShouldNotBeNil)

	log.Info("filter is complete, checking csv download",
		log.Data{
			"public_link":  filterOutputResource.Downloads.CSV.Public,
			"private_link": filterOutputResource.Downloads.CSV.Private,
			"href":         filterOutputResource.Downloads.CSV.HRef,
		})

	filteredCSVURLRaw := filterOutputResource.Downloads.CSV.Public
	if !isPublished {
		filteredCSVURLRaw = filterOutputResource.Downloads.CSV.Private
	}

	filteredCSVURL, err := url.Parse(filteredCSVURLRaw)
	if err != nil {
		log.ErrorC("failed to parse private filteredCSV URL", err, log.Data{"url": filteredCSVURLRaw})
		t.FailNow()
	}

	filteredCSVFilename := filteredCSVURL.Path

	log.Info("attempting to get filtered CSV file from S3",
		log.Data{
			"filtered_csv_filename": filteredCSVFilename,
			"filtered_csv_url":      filteredCSVURLRaw,
		})

	// read csv download from s3
	filteredCSVFile, err := getS3File(region, bucket, filteredCSVFilename, !isPublished)
	if err != nil {
		log.ErrorC("unable to find filtered csv download in s3", err, log.Data{"filtered_csv_url": filteredCSVURL, "filtered_csv_filename": filteredCSVFilename})
		t.Error()
		t.FailNow()
	}
	filteredCSVReader := csv.NewReader(filteredCSVFile)
	if err = checkFileRowCount(filteredCSVReader, 39); err != nil {
		log.ErrorC("unable to check file row count", err, nil)
		t.FailNow()
	}

	log.Debug("checking xlsx download",
		log.Data{
			"public_link":  filterOutputResource.Downloads.XLS.Public,
			"private_link": filterOutputResource.Downloads.XLS.Private,
			"href":         filterOutputResource.Downloads.XLS.HRef,
		})

	filteredXLSURLRaw := filterOutputResource.Downloads.XLS.Public
	if !isPublished {
		filteredXLSURLRaw = filterOutputResource.Downloads.XLS.Private
	}

	filteredXLSURL, err := url.Parse(filteredXLSURLRaw)
	if err != nil {
		log.ErrorC("failed to parse private filtered XLS URL", err, log.Data{"url": filteredXLSURLRaw})
		t.FailNow()
	}

	filteredXLSFilename := filteredXLSURL.Path
	log.Info("attempting to get filtered XLS file from S3",
		log.Data{
			"filtered_xls_filename": filteredXLSFilename,
			"filtered_xls_url":      filteredXLSURLRaw,
		})

	// read xls download from s3
	filteredXLSFile, err := getS3File(region, bucket, filteredXLSFilename, !isPublished)
	if err != nil {
		log.ErrorC("unable to find filtered xls download in s3", err, log.Data{"filtered_xls_url": filteredXLSURL, "filtered_xls_filename": filteredXLSFilename})
		t.FailNow()
	}
	So(filteredXLSFile, ShouldNotBeEmpty)
	filteredXLSFileSize, err := getS3FileSize(region, bucket, filteredXLSFilename, !isPublished)
	if err != nil {
		log.ErrorC("unable to extract size of filtered xls download in s3", err, log.Data{"filtered_xls_url": filteredXLSURL, "filtered_xls_filename": filteredXLSFilename})
		t.FailNow()
	}

	minExpectedXLSFileSize := int64(7465)
	maxExpectedXLSFileSize := int64(7469)
	So(*filteredXLSFileSize, ShouldBeBetweenOrEqual, minExpectedXLSFileSize, maxExpectedXLSFileSize)

	expectedCSVSize, _ := strconv.Atoi(filterOutputResource.Downloads.CSV.Size)
	testFileDownload(filterOutputResource.Downloads.CSV.HRef, expectedCSVSize, isPublished)

	expectedXLSSize, _ := strconv.Atoi(filterOutputResource.Downloads.XLS.Size)
	testFileDownload(filterOutputResource.Downloads.XLS.HRef, expectedXLSSize, isPublished)

	return filterBlueprintID, filterOutputID

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
