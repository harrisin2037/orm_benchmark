#!/bin/sh -eux
trap 'docker-compose down -v' EXIT
docker-compose down -v
docker-compose up -d
sleep 10
CGO_ENABLED=0 go test -run=^$ -bench=. -benchmem -count=10 | tee bench.txt
benchstat bench.txt
