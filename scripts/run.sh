#!/bin/sh

export $(egrep -v '^#' .env | xargs) > /dev/null
go run cmd/bazaar/main.go
