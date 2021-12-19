#!/bin/sh -eux
trap 'docker-compose down -v' EXIT
docker-compose down -v
docker-compose up -d
sleep 10
CGO_ENABLED=0 TZ='Asia/Hong_Kong' go test
