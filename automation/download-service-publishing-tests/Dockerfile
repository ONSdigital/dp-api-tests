FROM golang:latest

RUN go get github.com/ONSdigital/dp-api-tests/publishing/downloadService \
    && apt-get update \
    && apt-get install unzip \
    && wget https://releases.hashicorp.com/vault/0.9.5/vault_0.9.5_linux_amd64.zip \
    && unzip vault_0.9.5_linux_amd64.zip \
    && mv vault /bin/vault

WORKDIR $GOPATH/src/github.com/ONSdigital/dp-api-tests/publishing/downloadService 
