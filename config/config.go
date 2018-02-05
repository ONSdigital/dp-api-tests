package config

import "github.com/kelseyhightower/envconfig"

// Config values for the application.
type Config struct {
	DatasetAPIURL             string   `envconfig:"DATASET_API_URL"`
	ImportAPIURL              string   `envconfig:"IMPORT_API_URL"`
	FilterAPIURL              string   `envconfig:"FILTER_API_URL"`
	HierarchyAPIURL           string   `envconfig:"HIERARCHY_API_URL"`
	CodeListAPIURL            string   `envconfig:"CODELIST_API_URL"`
	RecipeAPIURL              string   `envconfig:"RECIPE_API_URL"`
	SearchAPIURL              string   `envconfig:"SEARCH_API_URL"`
	ElasticSearchAPIURL       string   `envconfig:"ELASTIC_SEARCH_URL"`
	MongoAddr                 string   `envconfig:"MONGODB_BIND_ADDR"`
	MongoDB                   string   `envconfig:"MONGODB_DATABASE"`
	MongoFiltersDB            string   `envconfig:"MONGODB_FILTERS_DATABASE"`
	MongoImportsDB            string   `envconfig:"MONGODB_IMPORTS_DATABASE"`
	Neo4jAddr                 string   `envconfig:"NEO4J_BIND_ADDR"`
	Brokers                   []string `envconfig:"KAFKA_ADDR"`
	ObservationsInsertedTopic string   `envconfig:"IMPORT_OBSERVATIONS_INSERTED_TOPIC"`
	ObservationConsumerGroup  string   `envconfig:"OBSERVATION_CONSUMER_GROUP"`
	ObservationConsumerTopic  string   `envconfig:"OBSERVATION_CONSUMER_TOPIC"`
}

var cfg *Config

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		DatasetAPIURL:             "http://localhost:22000",
		ImportAPIURL:              "http://localhost:21800",
		FilterAPIURL:              "http://localhost:22100",
		CodeListAPIURL:            "http://localhost:22400",
		RecipeAPIURL:              "http://localhost:22300",
		HierarchyAPIURL:           "http://localhost:22600",
		SearchAPIURL:              "http://localhost:23100",
		ElasticSearchAPIURL:       "http://localhost:9200",
		MongoAddr:                 "localhost:27017",
		MongoDB:                   "test",
		MongoImportsDB:            "test",
		MongoFiltersDB:            "test",
		Neo4jAddr:                 "bolt://localhost:7687",
		Brokers:                   []string{"localhost:9092"},
		ObservationsInsertedTopic: "import-observations-inserted",
		ObservationConsumerGroup:  "observation-extracted",
		ObservationConsumerTopic:  "observation-extracted",
	}

	return cfg, envconfig.Process("", cfg)
}
