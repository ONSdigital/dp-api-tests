package neo4j

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/ONSdigital/go-ns/log"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

const ObservationTestData = "../testDataSetup/neo4j/instance.cypher"
const HierarchyTestData = "../testDataSetup/neo4j/hierarchy.cypher"
const GenericHierarchyCPIHTestData = "../testDataSetup/neo4j/genericHierarchyCPIH.cypher"

// Datastore used to setup data within neo4j
type Datastore struct {
	instance   string
	testData   string
	connection bolt.Conn
}

// CypherTemplate allows cypher queries to be updated with new ID
type CypherTemplate struct {
	Instance string
}

// NewDatastore creates a new datastore for a test
func NewDatastore(uri, instance, testdata string) (*Datastore, error) {
	driver, err := bolt.NewDriver().OpenNeo(uri)
	if err != nil {
		return nil, err
	}
	return &Datastore{connection: driver, instance: instance, testData: testdata}, nil
}

// TeardownInstance removes all instance nodes within neo4j
func (ds *Datastore) TeardownInstance() error {
	query := fmt.Sprintf("MATCH (n:`_%s_Instance`)-[r]-(m)-[t]-(o) detach delete n,m,o;", ds.instance)
	results, err := ds.connection.QueryNeo(query, nil)
	if err != nil {
		return err
	}
	results.Close()
	return ds.connection.Close()
}

// DropDatabases cleans out all data that exists on neo4j
func DropDatabases(uri string) error {
	log.Info("dropping neo4j database", log.Data{"uri": uri})
	pool, err := bolt.NewDriverPool(uri, 1)
	if err != nil {
		return err
	}

	conn, err := pool.OpenPool()
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecNeo("MATCH(n) DETACH DELETE n", nil); err != nil {
		return err
	}

	log.Info("dropping databases complete", nil)

	return nil
}

// TeardownHierarchy removes all hierarchy nodes within neo4j
func (ds *Datastore) TeardownHierarchy() error {
	_, err := ds.connection.ExecNeo("MATCH (n:`_generic_hierarchy_node_e44de4c4-d39e-4e2f-942b-3ca10584d078`) DETACH DELETE n", nil)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("MATCH (n:`_hierarchy_node_%s_aggregate`) DETACH DELETE n", ds.instance)
	_, err = ds.connection.ExecNeo(query, nil)
	if err != nil {
		return err
	}
	return ds.connection.Close()
}

// Setup the neo4j database
func (ds *Datastore) Setup() error {
	t, err := template.ParseFiles(ds.testData)
	if err != nil {
		return err
	}
	query := new(bytes.Buffer)
	t.Execute(query, CypherTemplate{Instance: ds.instance})
	_, err = ds.connection.ExecNeo(query.String(), nil)
	if err != nil {
		return err
	}
	return err
}

// CreateGenericHierarchy the neo4j database
func (ds *Datastore) CreateGenericHierarchy(hierarchyCode string) error {
	_, err := ds.connection.ExecNeo("MATCH (n:`_generic_hierarchy_node_"+hierarchyCode+"`) DETACH DELETE n", nil)
	if err != nil {
		return err
	}

	//b, err := ioutil.ReadFile(ds.testData)
	file, err := os.Open(ds.testData)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, err = ds.connection.ExecNeo(line, nil)
		if err != nil {
			log.ErrorC("encountered error writing query to neo4j", err, log.Data{"cypher_file": ds.testData, "cypher_line": scanner.Text()})
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	log.Info("successfully loaded data into neo4j", log.Data{"cypher_file": ds.testData})
	return nil
}
