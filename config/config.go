package config

import "github.com/kelseyhightower/envconfig"

// Config values for the application.
type Config struct {
	CodeListAPIURL            string   `envconfig:"CODELIST_API_URL"`
	DatasetAPIURL             string   `envconfig:"DATASET_API_URL"`
	DownloadServiceURL        string   `envconfig:"DOWNLOAD_SERVICE_URL"`
	FilterAPIURL              string   `envconfig:"FILTER_API_URL"`
	HierarchyAPIURL           string   `envconfig:"HIERARCHY_API_URL"`
	ImportAPIURL              string   `envconfig:"IMPORT_API_URL"`
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
	EncryptionDisabled        bool     `envconfig:"ENCRYPTION_DISABLED"`
	VaultToken                string   `envconfig:"VAULT_TOKEN"`
	VaultAddress              string   `envconfig:"VAULT_ADDR"`
	VaultPath                 string   `envconfig:"VAULT_PATH"`
}

var cfg *Config

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		CodeListAPIURL:            "http://localhost:22400",
		DatasetAPIURL:             "http://localhost:22000",
		DownloadServiceURL:        "http://localhost:23600",
		FilterAPIURL:              "http://localhost:22100",
		HierarchyAPIURL:           "http://localhost:22600",
		ImportAPIURL:              "http://localhost:21800",
		RecipeAPIURL:              "http://localhost:22300",
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
		EncryptionDisabled:        false,
		VaultAddress:              "http://localhost:8200",
		VaultToken:                "",
		VaultPath:                 "secret/shared/psk",
	}

	return cfg, envconfig.Process("", cfg)
}
