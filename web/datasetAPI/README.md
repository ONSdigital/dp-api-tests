Dataset API Tests
================

### Getting started

This package will test all endpoints that exist within the Dataset API when
running in web subnet

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongodb
neo4j
zookeeper
kafka
dp-dataset-api
```

`dp-dataset-api` should be run with `make acceptance-web`
