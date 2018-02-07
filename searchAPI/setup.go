package searchAPI

import (
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-api-tests/testDataSetup/elasticsearch"
	"github.com/ian-kent/go-log/log"
)

func createSearchIndex(url, instanceID, dimension string) error {
	log.Info("nath got here", nil)

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	directory := ".."
	if filepath.Base(currentPath) == "dp-api-tests" {
		directory = "."
	}

	index := elasticsearch.Index{
		InstanceID:   instanceID,
		Dimension:    dimension,
		TestDataFile: directory + "/testDataSetup/elasticsearch/testData.json",
		URL:          url,
		MappingsFile: directory + "/testDataSetup/elasticsearch/mappings.json",
	}

	if err := index.CreateSearchIndex(); err != nil {
		return err
	}

	return nil
}
