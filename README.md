dp-api-tests
================


### Getting started

`go test ./...` - run all the tests

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

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
