Search API Tests
================

### Getting started

This package will test all publishing endpoints that exist within the Search API

#### Services and software

The following software needs to be running for acceptance tests to be able to
pass:

```text
mongo db
zookeeper
kafka
elasticsearch 5.x
dataset API
search API
```

If elasticsearch is running on version 6.x then the highlighting tests will fail,
as this is a breaking change from elasticsearch version 5 to 6.

Both API's will need to have the following configuration:

```
export MONGODB_DATABASE=test ENABLE_PRIVATE_ENDPOINTS=true
```

Or one can run make make acceptance-publishing
