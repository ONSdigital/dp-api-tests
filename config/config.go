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
	MongoAddr                 string   `envconfig:"MONGODB_BIND_ADDR"`
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
		MongoAddr:                 "localhost:27017",
		HierarchyAPIURL:           "http://localhost:22600",
		Neo4jAddr:                 "bolt://localhost:7687",
		Brokers:                   []string{"localhost:9092"},
		ObservationsInsertedTopic: "import-observations-inserted",
		ObservationConsumerGroup:  "observation-extracted",
		ObservationConsumerTopic:  "observation-extracted",
	}

	return cfg, envconfig.Process("", cfg)
}
