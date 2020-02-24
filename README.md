# es-mdb

> Complete integration between MongoDB and ElasticSearch using Go

## Overview

This project contains services to sync data between MongoDB and ElasticSearch. The following list has the all the solution components and their description:

- `searchService`: this service integrates with ElasticSearch with basic search queries.
- `updateESService`: this service listens for changes on a MongoDB collection using change streams and updates the ElasticSearch index. Whenever an update to ElasticSearch fails, it will store the update on another collection so it can be processed at a later time.
- `retryESUpdate`: this cron job will retry all failed updates to ElasticSearch every two hours.

## Setup

Add a `.env` file with the following environment variables:

```bash
ELASTICSEARCH_URL=http://localhost:9200
MDB_URL=mongodb://mdb:27017/?replicaSet=rs0
DB_NAME=es_mdb
```

If on MacOS, add the following to your /etc/hosts file:

```bash
127.0.0.1  mdb
```

## Run

**Docker:**

```bash
# set up MongoDB container with replica set configuration
docker-compose up -d; sleep 5; docker exec mdb mongo ./scripts/rsInit.js;
```

**Services:**

**Build:**

```bash
# running custom images not in docker-compose. Replica set should be configured before running these images
source .env;
docker build -t update-es:1.0 --rm -f ./cmd/updateESService/dockerfile . --build-arg ELASTICSEARCH_URL=${ELASTICSEARCH_URL} --build-arg MDB_URL=${MDB_URL} --build-arg DB_NAME=${DB_NAME};
docker build -t search-es:1.0 --rm -f ./cmd/searchService/dockerfile . --build-arg ELASTICSEARCH_URL=${ELASTICSEARCH_URL} --build-arg MDB_URL=${MDB_URL} --build-arg DB_NAME=${DB_NAME};
```

**Run:**

```bash
docker run -d update-es:1.0 --env-file .env
docker run -d search-es:1.0 --env-file .env
```

**Cleanup:**

```bash
docker-compose down -v;
```

## WIP
