#!/bin/sh -eux
trap 'docker-compose down -v' EXIT

docker-compose down -v
docker-compose up -d

sleep 10

CGO_ENABLED=0 TZ='Asia/Hong_Kong' go test -run=^$ -bench=. -benchmem -count=1 -cpu 4 | tee bench-purego.txt
CGO_ENABLED=1 TZ='Asia/Hong_Kong' go test -run=^$ -bench=. -benchmem -count=1 -cpu 4 | tee bench-cgo.txt

benchstat bench-cgo.txt bench-purego.txt
