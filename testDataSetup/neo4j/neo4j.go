package neo4j

import "github.com/ONSdigital/dp-api-tests/config"

// Teardown ...
func Teardown() error {
  config, err := config.Get()
	if err != nil {
		return err
	}

	return nil
}

// Setup ...
func Setup() error {
  config, err := config.Get()
	if err != nil {
		return err
	}

	return nil
}
