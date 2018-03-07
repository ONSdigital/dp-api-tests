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
export RSA_PRIVATE_KEY=<location of private key>
```

To create your private key, you can do the following:

1) Run `openssl genrsa -out private.pem<number_of_bytes>` where `number_of_bytes`
   can be 1024 or 2048
2) Add `export RSA_PRIVATE_KEY=$(cat $HOME/<file_location>/<filename>.pem)`
   to your .bashrc or .bash_profile
3) Run `exec $SHELL -l` in terminal or open new terminal window

The endToEndTest can be run with or without decryption/encryption by exporting
the `ENCRYPTION_DISABLED` environment variable to `true` or `false` respectively.

This environment variable will need to be set for backend services if set to `false`.

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
