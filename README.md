dp-api-tests
================

### Getting started

These tests should be run against empty databases. To ensure this, your
applications should be configured to use a different database.
Do this by running the following and then restarting your applications:
```
export MONGODB_DATABASE=test
export MONGODB_FILTERS_DATABASE=test
export MONGODB_IMPORTS_DATABASE=test
```
or by running `make acceptance` instead of `make debug` to start up the
following services:
```
dp-dataset-api
dp-filter-api
dp-import-api
```

If you're running applications using websysd ensure `InheritEnvironment`
is set to true in the `websysd.json` file.

The acceptance tests are split between running in publishing or web and result
in different results/responses; by default you should be running the services in
publishing.

To test services in publishing run the following command:

`cd publishing; go test ./...`

To run the web tests, make sure all instances of your API services are running
in web and then run the following command:

`cd web; go test ./...`

### Testing standards

All tests should be written in the following structure:
* Teardown all data related to test
* Setup all data related to test, do NOT use API's to setup this data, if this fails then the test fails for the wrong reason.
* Run test
* Teardown all data related to test

### Configuration

An overview of the configuration options available, either as a table of
environment variables, or with a link to a configuration guide.

| Environment variable               | Default                      | Description
| ---------------------------------- | ---------------------------- | -----------
| CODELIST_API_URL                   | http://localhost:22400       | The host name for the Codelist API
| DATASET_API_URL                    | http://localhost:22000       | The host name for the Dataset API
| DOWNLOAD_SERVICE_URL               | http://localhost:23600       | The host name for the Download Service
| FILTER_API_URL                     | http://localhost:22100       | The host name for the Filter API
| HIERARCHY_API_URL                  | http://localhost:22600       | The host name for the Hierarchy API
| IMPORT_API_URL                     | http://localhost:21800       | The host name for the Import API
| RECIPE_API_URL                     | http://localhost:22300       | The host name for the Recipe API
| SEARCH_API_URL                     | http://localhost:23100       | The host name for the Search API
| ELASTIC_SEARCH_URL                 | http://localhost:9200        | The host name for elasticsearch
| MONGODB_BIND_ADDR                  | localhost:27017              | The MongoDB bind address
| MONGODB_DATABASE                   | test                         | The Dataset API mongo database
| MONGODB_FILTERS_DATABASE           | test                         | The Filter API mongo database
| MONGODB_IMPORTS_DATABASE           | test                         | The Import API mongo database
| NEO4J_BIND_ADDR                    | bolt://localhost:7687        | The Neo4j bind address
| KAFKA_ADDR                         | localhost:9092               | The list of kafka hosts
| IMPORT_OBSERVATIONS_INSERTED_TOPIC | import-observations-inserted | The Kafka topic to produce events for the number of inserted observations
| OBSERVATION_CONSUMER_GROUP         | observation-extracted        | The Kafka consumer group to consume observation extracted events from
| OBSERVATION_CONSUMER_TOPIC         | observation-extracted        | The Kafka topic to consume observation extracted events from
| ENCRYPTION_DISABLED                | false                        | A boolean flag to identify if encryption of files is disabled or not
| VAULT_ADDR                         | http://localhost:8200        | The vault address
| VAULT_TOKEN                        | -                            | Vault token required for the client to talk to vault. (Use `make debug` to create a vault token)
| VAULT_PATH                         | secret/shared/psk            | The path where the psks will be stored in for vault

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

####

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
