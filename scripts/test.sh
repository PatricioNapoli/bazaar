#!/bin/sh

ENVFILE=.env
if [ -f "$ENVFILE" ]; then
     export $(egrep -v '^#' .env | xargs) > /dev/null
fi

mkdir -p output
go test -v cmd/bazaar/main_test.go
