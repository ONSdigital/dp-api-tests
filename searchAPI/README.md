Search API Tests
================

### Getting started

This package will test all endpoints that exist within the Filter API

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongo db
zookeeper
kafka
elasticsearch
dataset API
search API
```

Both API's will need to have the following configuration:

```
export MONGODB_DATABASE=test
export MONGODB_FILTERS_DATABASE=test
export MONGODB_IMPORTS_DATABASE=test
```
