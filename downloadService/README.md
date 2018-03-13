# Download Service Tests

### Getting Started
This package will test all endpoints that exist within the download service.

You should run `make test` which will setup valid vault credentials. If you run `go test ./...` then private file tests will be skipped as no vault credentials will exist.

### Services and software

dataset-api (`make acceptance`)
download-service (`make acceptance`)
vault