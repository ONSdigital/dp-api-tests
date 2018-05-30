Import API Tests
================

### Getting started

This package will test all endpoints that exist within the Import API

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongo db
zookeeper
kafka
Dataset API
Recipe API
```

When running the dataset API, one should use a publishing instance of the
service, this can be done by running `make acceptance-publishing`
