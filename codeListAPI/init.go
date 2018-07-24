package codeListAPI

import (
	"github.com/ONSdigital/dp-bolt/bolt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/ONSdigital/dp-api-tests/config"
	"os"
	"github.com/pkg/errors"
	"fmt"
	"testing"
)

type DB struct {
	bolt *bolt.DB
	cfg  *config.Config
	t    *testing.T
}

func NewDB(t *testing.T) *DB {
	cfg, err := config.Get()
	if err != nil {
		t.Fatal(errors.New("error getting test config"))
		return nil
	}

	pool, err := golangNeo4jBoltDriver.NewClosableDriverPool(cfg.Neo4jAddr, 2)
	if err != nil {
		t.Fatal(errors.New("error creating neo4j driver pool"))
		return nil
	}

	return &DB{
		cfg:  cfg,
		bolt: bolt.New(pool),
		t:    t,
	}
}

func (db *DB) Setup(stmts ... bolt.Stmt) {
	if len(stmts) == 0 {
		return
	}

	for i, s := range stmts {
		_, _, err := db.bolt.Exec(s)
		if err != nil {
			db.t.Fatal(errors.WithMessage(err, fmt.Sprintf("stmt index: %d", i)))
			os.Exit(1) // don't think this is required as Fail(...) will call runtime.GoExit()
		}
	}
}

func (db *DB) TearDown() {
	db.t.Log("attempting tear down")
	_, _, err := db.bolt.Exec(tearItDown)
	if err != nil {
		db.t.Fatal(errors.WithMessage(err, "tear down failure"))
		os.Exit(1) // don't think this is required as Fail(...) will call runtime.GoExit()
	}
	db.t.Log("tear down completed successfully")
}
