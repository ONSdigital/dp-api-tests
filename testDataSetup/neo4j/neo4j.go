package neo4j

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"github.com/ONSdigital/go-ns/log"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	. "github.com/smartystreets/goconvey/convey"
	"errors"
)

const ObservationTestData = "../../testDataSetup/neo4j/instance.cypher"
const HierarchyTestData = "../../testDataSetup/neo4j/hierarchy.cypher"
const GenericHierarchyCPIHTestData = "../testDataSetup/neo4j/genericHierarchyCPIH.cypher"

const codeListCPIHTestData = "../testDataSetup/neo4j/codeListCPIH.cypher"

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

// CreateGenericHierarchy the neo4j database
func (ds *Datastore) CreateGenericHierarchy(hierarchyCode string) error {
	_, err := ds.connection.ExecNeo("MATCH (n:`_generic_hierarchy_node_"+hierarchyCode+"`) DETACH DELETE n", nil)
	if err != nil {
		return err
	}

	file, err := os.Open(ds.testData)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.ErrorC("CreateGenericHierarchy", err, nil)
		}
	}()

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

func (ds *Datastore) CreateCPIHCodeList() error {
	_, err := ds.connection.ExecNeo("MATCH (n:`_code`) DETACH DELETE n", nil)
	if err != nil {
		return err
	}

	_, err = ds.connection.ExecNeo("MATCH (n:`_code_list`) DETACH DELETE n", nil)
	if err != nil {
		return err
	}

	file, err := os.Open(codeListCPIHTestData)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.ErrorC("CreateCPIHCodeList", err, nil)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, err = ds.connection.ExecNeo(line, nil)
		if err != nil {
			log.ErrorC("encountered error writing query to neo4j", err, log.Data{"cypher_file": codeListCPIHTestData, "cypher_line": scanner.Text()})
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	log.Info("successfully loaded data into neo4j", log.Data{"cypher_file": codeListCPIHTestData})

	return nil
}

// GetInstanceProperties returns a map of properties that are stored on the instance node
func (ds *Datastore) GetInstanceProperties(instanceID string) (map[string]interface{}, error) {

	query := fmt.Sprintf("MATCH (i:`_%s_Instance`) RETURN i", instanceID)
	stmt, err := ds.connection.PrepareNeo(query)
	So(err, ShouldBeNil)
	defer stmt.Close()

	rows, err := stmt.QueryNeo(nil)
	So(err, ShouldBeNil)
	defer rows.Close()

	data, _, err := rows.NextNeo()

	nodeData := data[0]
	graphNode, ok := nodeData.(graph.Node)
	if ok {
		return graphNode.Properties, nil
	}

	return nil, errors.New("failed to retrieve properties from neo4j instance node")
}

// CreateInstanceNode creates a new instance node for the given instance ID
func (ds *Datastore) CreateInstanceNode(instanceID string) (int64, error) {

	stmt, err := ds.connection.PrepareNeo(fmt.Sprintf("CREATE (i:`_%s_Instance`) RETURN i", instanceID))
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecNeo(nil)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CleanUpInstance removes the instance node for the given instance ID
func (ds *Datastore) CleanUpInstance(instanceID string) (error) {
	log.Info("cleaning up test instance", log.Data{"instanceID": instanceID})

	query := fmt.Sprintf("MATCH (i:`_%s_Instance`) DETACH DELETE i", instanceID)
	stmt, err := ds.connection.PrepareNeo(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecNeo(nil)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	So(count, ShouldEqual, int64(1))

	log.Info("cleaning up test instance complete", log.Data{"instanceID": instanceID})

	return nil
}
