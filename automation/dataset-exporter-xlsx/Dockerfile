FROM maven:3.5.4-jdk-8-alpine

RUN apk update && apk upgrade \
    && apk add --no-cache bash git openssh unzip wget make fontconfig ttf-dejavu \
    && git clone https://github.com/ONSdigital/dp-dataset-exporter-xlsx.git \
    && wget https://releases.hashicorp.com/vault/0.9.5/vault_0.9.5_linux_amd64.zip \
    && unzip vault_0.9.5_linux_amd64.zip \
    && mv vault /bin/vault \
    && cd dp-dataset-exporter-xlsx \
    && make build

WORKDIR dp-dataset-exporter-xlsx