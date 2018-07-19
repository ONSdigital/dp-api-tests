package codeListAPI

import (
	"github.com/ONSdigital/dp-bolt/bolt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"github.com/pkg/errors"
	"fmt"
)

type DB struct {
	bolt *bolt.DB
	cfg  *config.Config
}

func NewDB() (*DB, error) {
	cfg, _ := config.Get()
	pool, err := golangNeo4jBoltDriver.NewClosableDriverPool(cfg.Neo4jAddr, 2)
	if err != nil {
		return nil, err
	}

	return &DB{
		cfg:  cfg,
		bolt: bolt.New(pool),
	}, nil
}

func (db *DB) setUp(stmts ... bolt.Stmt) error {
	if len(stmts) == 0 {
		return nil
	}

	for i, s := range stmts {
		_, err := db.bolt.Exec(s)
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("stmt index: %d", i))
		}
	}
	return nil
}

func (db *DB) tearDown() {
	_, err := db.bolt.Exec(tearItDown)
	if err != nil {
		log.ErrorC("tear down failure", err, nil)
		os.Exit(1)
	}
}
