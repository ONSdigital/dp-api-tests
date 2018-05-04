Dataset API Tests
================

### Getting started

This package will test all endpoints that exist within the Dataset API

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongo db
neo4j
zookeeper
kafka
dp-dataset-api
dp-auth-api-stub (mimics zebedee authentication)
```

`dp-dataset-api` should be run with `make acceptance`
