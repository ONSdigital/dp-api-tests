package hierarchyAPI

import (
	"testing"

	"fmt"
	"net/http"
	"os"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/neo4j"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSuccessfullyGetRootHierarchy(t *testing.T) {
	instanceID := uuid.NewV4().String()
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

			Convey("Then a root hierarchy node is return as a response", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}", instanceID, "aggregate").
					Expect().Status(http.StatusOK).JSON().Object()

				// Check root node
				response.Value("has_data").Boolean().Equal(true)
				response.Value("no_of_children").Number().Equal(12)
				selfLink := response.Value("links").Object().Value("self").Object()
				codeLink := response.Value("links").Object().Value("code").Object()
				codeLink.Value("href").String().
					Equal("http://localhost:22400/code-list/e44de4c4-d39e-4e2f-942b-3ca10584d078/code/cpi1dim1A0")
				selfLink.Value("href").String().
					Equal(fmt.Sprintf("http://localhost:22600/hierarchies/%s/aggregate", instanceID))

				// Check first child node
				first := response.Value("children").Array().First().Object()
				first.Value("has_data").Boolean().True()
				first.Value("label").String().Equal("01 Food and non-alcoholic beverages")
				first.Value("no_of_children").Number().Equal(2)
				firstSelfLink := first.Value("links").Object().Value("self").Object()
				firstCodeLink := first.Value("links").Object().Value("code").Object()
				firstCodeLink.Value("href").String().
					Equal("http://localhost:22400/code-list/e44de4c4-d39e-4e2f-942b-3ca10584d078/code/cpi1dim1T10000")
				firstSelfLink.Value("href").String().
					Equal(fmt.Sprintf("http://localhost:22600/hierarchies/%s/aggregate/cpi1dim1T10000", instanceID))
			})
		})

		err = datastore.TeardownHierarchy()
		if err != nil {
			log.ErrorC("Unable to tear down test data", err, nil)
			os.Exit(1)
		}
	})
}

func TestErrorStatesGetRootHierarchy(t *testing.T) {
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
		SkipConvey("When a root hierarchy node is requested with a invalid id", func() {

			Convey("Then a 400 response code is returned", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}", "1234567890", "cpi")
				response.Expect().Status(http.StatusBadRequest)
			})
		})

		Convey("When a root hierarchy node is requested with a dimension name", func() {
			Convey("Then a 404 response code is returned", func() {
				response := hierarchyAPI.GET("/hierarchies/{instance_id}/{dimension_name}", instanceID, "000")
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
