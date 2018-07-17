package codeListAPI

import (
	"github.com/ONSdigital/dp-bolt/bolt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/go-ns/log"
	"os"
)

var (
	tearItDown = bolt.Stmt{Query: "MATCH(n:`_api_test`) DETACH DELETE n", Params: nil}

	setupStmts = []bolt.Stmt{
		{Query: "CREATE (node:`_code_list`:`_name_gender`:`_api_test` { label:'sex', edition:'one-off' })", Params: nil},
		{Query: "MATCH (parent:`_code_list`:`_name_gender`) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)", Params: map[string]interface{}{"v": "2", "l": "Female"}},
		{Query: "MATCH (parent:`_code_list`:`_name_gender`) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)", Params: map[string]interface{}{"v": "1", "l": "Male"}},
		{Query: "MATCH (parent:`_code_list`:`_name_gender`) WITH parent CREATE (node:`_code`:`_api_test` { value:{v}})-[:usedBy { label: {l}}]->(parent)", Params: map[string]interface{}{"v": "0", "l": "All"}},
	}

	codeListName = "sex"
	codeListID   = "gender"
	edition      = "one-off"
)

type DB struct {
	bolt *bolt.DB
	cfg  *config.Config
}

func NewDB() (*DB, error) {
	cfg, _ := config.Get()
	pool, err := golangNeo4jBoltDriver.NewClosableDriverPool(cfg.Neo4jAddr, 4)
	if err != nil {
		return nil, err
	}

	return &DB{
		cfg:  cfg,
		bolt: bolt.New(pool),
	}, nil
}

func (db *DB) setUp() error {
	for _, stmt := range setupStmts {
		_, err := db.bolt.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) tearDown() {
	_, err := db.bolt.Exec(tearItDown)
	if err != nil {
		log.ErrorC("test teardown failure", err, nil)
		os.Exit(1)
	}
}
