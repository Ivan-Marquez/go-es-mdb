# es-mdb

> Complete integration between MongoDB and ElasticSearch using Go

## Overview

This project contains services to sync data between MongoDB and ElasticSearch. The following list contains all solution components and their description:

- `searchService`: this service integrates with ElasticSearch using basic search queries.
- `updateESService`: this service listens for changes on a MongoDB collection using change streams and updates the ElasticSearch index. Whenever an update to ElasticSearch fails, it will store the update on another collection so it can be processed at a later time.
- `retryESUpdate`: this cron job will retry all failed updates to ElasticSearch every two hours.

After updating a record in MongoDB, the `updateESService` will detect the change and will send the update to ElasticSearch.s

## Setup

Add a `.env` file with the following environment variables:

```bash
ELASTICSEARCH_URL=http://localhost:9200
MDB_URL=mongodb://mdb:27017/?replicaSet=rs0
DB_NAME=es_mdb
```

If on MacOS, add the following to your `/etc/hosts` file:

```bash
127.0.0.1  mdb
```

## Run

### MongoDB / ElasticSearch containers

**Docker Compose:**

```bash
# set up MongoDB container with replica set configuration
docker-compose up -d; sleep 5; docker exec mdb mongo ./scripts/rsInit.js;
```

### Mock Data

You can also add mock data to MongoDB and ElasticSearch by running the following the following script (optional):

```bash
# NOTE: MongoDB and ElasticSearch containers must be up and running before adding mock data
go run tools/mockData/main.go
```

### Services

**Build:**

```bash
# custom images not included in docker-compose. Replica set must be configured before running these images!
source .env;
docker build -t update-es:1.0 -f ./cmd/updateES/dockerfile . --build-arg ELASTICSEARCH_URL=${ELASTICSEARCH_URL} --build-arg MDB_URL=${MDB_URL} --build-arg DB_NAME=${DB_NAME};
docker build -t search-es:1.0 -f ./cmd/searchES/dockerfile . --build-arg ELASTICSEARCH_URL=${ELASTICSEARCH_URL} --build-arg MDB_URL=${MDB_URL} --build-arg DB_NAME=${DB_NAME};
```

**Run:**

```bash
# run images and add them to mdb container network
docker run -d --name update --env-file .env --network go-es-mdb_mdbn --rm update-es:1.0
docker run -d --name search --env-file .env --network go-es-mdb_mdbn --rm search-es:1.0
```

**Cleanup:**

```bash
docker-compose down -v;
```

## Example
