#!/usr/bin/env sh
#
# run tests.
#

docker-compose run --rm vhstatus go test -v \
    ./cmd/... \
    ./internal/...

