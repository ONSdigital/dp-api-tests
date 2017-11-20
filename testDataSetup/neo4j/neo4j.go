package neo4j

import (
	"bytes"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"html/template"
)

const ObservationTestData = "testDataSetup/neo4j/data.cypher"

type Datastore struct {
	instance   string
	testData   string
	connection bolt.Conn
}

type CypherTemplate struct {
	Instance string
}

func NewDatastore(uri, instance, testdata string) *Datastore {
	driver, err := bolt.NewDriver().OpenNeo(uri)
	if err != nil {
		panic(err)
	}
	return &Datastore{connection: driver, instance: instance, testData: testdata}
}

// Teardown ...
func (ds *Datastore) TeardownObservation() error {
	query := fmt.Sprintf("MATCH (n:`_%s_Instance`)-[r]-(m)-[t]-(o) detach delete n,m,o;", ds.instance)
	_, err := ds.connection.QueryNeo(query, nil)
	if err != nil {
		return err
	}
	return ds.connection.Close()
}

// Setup ...
func (ds *Datastore) Setup() error {
	t, err := template.ParseFiles("testDataSetup/neo4j/data.cypher")
	if err != nil {
		return err
	}
	query := new(bytes.Buffer)
	t.Execute(query, CypherTemplate{Instance: ds.instance})
	_, err = ds.connection.QueryNeo(query.String(), nil)
	if err != nil {
		return err
	}

	return err
}

//func main() {
//	ds := NewDatastore("bolt://localhost:7687", "123-234", ObservationTestData)
//	ds.Setup()
//	ds.TeardownObservation()
//}
