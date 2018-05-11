Dataset API Tests
================

### Getting started

This package will test all endpoints that exist within the Dataset API when
running in publishing subnet.

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongodb
neo4j
zookeeper
kafka
dp-dataset-api
dp-auth-api-stub (mimics zebedee authentication)
```

`dp-dataset-api` should be run with `make acceptance-publishing`

#### Note

If an endpoint is only available on publishing, remember to add a test to
web/datasetAPI/hidden_endpoints_test.go to check request returns 404 in web subnet
