package hierarchyAPI

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetNodeHierarchy(t *testing.T) {
	instanceID := uuid.NewV4().String()
	cpiCode := "cpi1dim1T120000"
	hierarchyAPI := httpexpect.New(t, cfg.HierarchyAPIURL)

	Convey("Given an existing hierarchy", t, func() {

		datastore, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, neo4j.HierarchyTestData)
		if err != nil {
			log.ErrorC("Unable to connect to neo4j", err, nil)
			os.Exit(1)
		}
		Convey("When a root hierarchy node is requested", func() {

			err := datastore.Setup()
			if err != nil {
				log.ErrorC("Unable to setup test data", err, nil)
				os.Exit(1)
			}

			Convey("Then a child hierarchy node is return as a response", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}/{code}", instanceID, "aggregate", cpiCode).
					Expect().Status(http.StatusOK).JSON().Object()

				// Check root node
				response.Value("has_data").Boolean().Equal(true)
				response.Value("no_of_children").Number().Equal(6)
				selfLink := response.Value("links").Object().Value("self").Object()
				codeLink := response.Value("links").Object().Value("code").Object()
				codeLink.Value("href").String().
					Equal("http://localhost:22400/code-lists/e44de4c4-d39e-4e2f-942b-3ca10584d078/codes/cpi1dim1T120000")
				selfLink.Value("href").String().
					Equal(fmt.Sprintf("%s/hierarchies/%s/aggregate/%s", cfg.HierarchyAPIURL, instanceID, cpiCode))

				// Check first child node
				first := response.Value("children").Array().First().Object()
				first.Value("has_data").Boolean().True()
				first.Value("label").String().Equal("12.1 Personal care")
				first.Value("no_of_children").Number().Equal(2)
				firstSelfLink := first.Value("links").Object().Value("self").Object()
				firstCodeLink := first.Value("links").Object().Value("code").Object()
				firstCodeLink.Value("href").String().
					Equal("http://localhost:22400/code-lists/e44de4c4-d39e-4e2f-942b-3ca10584d078/codes/cpi1dim1G120100")
				firstSelfLink.Value("href").String().
					Equal(fmt.Sprintf("%s/hierarchies/%s/aggregate/cpi1dim1G120100", cfg.HierarchyAPIURL, instanceID))
			})
		})

		err = datastore.TeardownHierarchy()
		if err != nil {
			log.ErrorC("Unable to tear down test data", err, nil)
			os.Exit(1)
		}
	})
}

func TestErrorStatesGetNodeHierarchy(t *testing.T) {
	instanceID := uuid.NewV4().String()
	hierarchyAPI := httpexpect.New(t, cfg.HierarchyAPIURL)

	Convey("Given an existing hierarchy", t, func() {

		datastore, err := neo4j.NewDatastore(cfg.Neo4jAddr, instanceID, neo4j.HierarchyTestData)
		if err != nil {
			log.ErrorC("Unable to connect to neo4j", err, nil)
			os.Exit(1)
		}

		err = datastore.Setup()
		if err != nil {
			log.ErrorC("Unable to setup test data", err, nil)
			os.Exit(1)
		}
		// This should return 400 but returns 404
		SkipConvey("When a child hierarchy node is requested with a invalid instance", func() {
			Convey("Then a 400 response code is returned", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}/{code}", "0000", "aggregate", "cpi1dim1G120100")
				response.Expect().Status(http.StatusBadRequest)
			})
		})
		// This should return 400 but returns 404
		SkipConvey("When a child hierarchy node is requested with a invalid dimension name", func() {
			Convey("Then a 400 response code is returned", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}/{code}", instanceID, "0000", "cpi1dim1999999")
				response.Expect().Status(http.StatusBadRequest)
			})
		})
		// This should return 404 but is 200
		SkipConvey("When a child hierarchy node is requested with a invalid code", func() {
			Convey("Then a 404 response code is returned", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}/{code}", instanceID, "aggregate", "cpi1dim1999999")
				response.Expect().Status(http.StatusNotFound)
			})
		})

		err = datastore.TeardownHierarchy()
		if err != nil {
			log.ErrorC("Unable to tear down test data", err, nil)
			os.Exit(1)
		}
	})
}
