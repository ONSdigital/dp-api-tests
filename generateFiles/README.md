Generate Files Tests
================

### Getting started

This package will test a sunny day scenario of the backend and api services.

Make sure all API's have access to the following environment variables before
starting the services.

```
export MONGODB_DATABASE=test
export MONGODB_FILTERS_DATABASE=test
export MONGODB_IMPORTS_DATABASE=test
```

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongo db
zookeeper
kafka
neo4j
recipe API
import API
dataset API
filter API - Not yet
dimension extractor
dimension import
observation extractor
observation importer
import tracker
import reporter
dataset exporter
xlsx dataset exporter
```
