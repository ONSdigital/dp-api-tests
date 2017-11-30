package neo4j

import (
	"bytes"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"html/template"
)

const ObservationTestData = "../testDataSetup/neo4j/instance.cypher"
const HierarchyTestData = "../testDataSetup/neo4j/hierarchy.cypher"

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
