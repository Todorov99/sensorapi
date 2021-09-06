#!/bin/sh

docker run --publish 5432:5432  postres-db:1.0.1
docker run --rm -p 5050:5050 thajeztah/pgadmin4

docker run -p 8086:8086 influxdb

docker run -p 8086:8086 influxdb