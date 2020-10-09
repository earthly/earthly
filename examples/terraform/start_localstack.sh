#!/usr/bin/env bash

docker run -d -p 4566:4566 -e SERVICES=s3 --name localstack --network="host" localstack/localstack:0.11.5