dp-api-tests
================

### Getting started

These tests should be run against empty databases. To ensure this, your applications should be configured to use a different database.
Do this by running the following and then restarting your applications:
```
export MONGODB_DATABASE=test
export MONGODB_IMPORTS_DATABASE=test
```

`go test ./...` - run all the tests

### Testing standards

All tests should be written in the following structure:
* Teardown all data related to test
* Setup all data related to test, do NOT use API's to setup this data, if this fails then the test fails for the wrong reason.
* Run test
* Teardown all data related to test

### Configuration

An overview of the configuration options available, either as a table of
environment variables, or with a link to a configuration guide.

| Environment variable       | Default                              | Description
| -------------------------- | -------------------------------------| -----------
| CODELIST_API_URL           | http://localhost:22400               | The host name for the Codelist API
| DATASET_API_URL            | http://localhost:22000               | The host name for the Dataset API
| FILTER_API_URL             | http://localhost:22100               | The host name for the Filter API
| IMPORT_API_URL             | http://localhost:21800               | The host name for the Import API
| MONGODB_BIND_ADDR          | localhost:27017                      | The MongoDB bind address

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

####

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
