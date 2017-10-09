package config

func ImportAPIURL() string {
	return "http://localhost:21800"
}

func FilterAPIURL() string {
	return "http://localhost:22100"
}

func CodeListAPIURL() string {
	return "http://localhost:22400"
}

// Config values for the application.
// type Config struct {
// 	DatasetAPIURL  string `envconfig:"DATASET_API_URL"`
// 	ImportAPIURL   string `envconfig:"IMPORT_API_URL"`
// 	FilterAPIURL   string `envconfig:"FILTER_API_URL"`
// 	CodeListAPIURL string `envconfig:"CODELIST_API_URL"`
// }

// // Get the configuration values from the environment or provide the defaults.
// func Get() (*Config, error) {

// 	cfg := &Config{

// 		DatasetAPIURL:  "http://localhost:22000",
// 		ImportAPIURL:   "http://localhost:21800",
// 		FilterAPIURL:   "http://localhost:22100",
// 		CodeListAPIURL: "http://localhost:22400",
// 	}

// 	return cfg, envconfig.Process("", cfg)
// }
