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
elasticsearch (version can be either 5.x or 6.x)
dataset API
search API
```

Both API's will need to have the following configuration:

```
export MONGODB_DATABASE=test
```
