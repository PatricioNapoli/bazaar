#!/bin/sh

ENVFILE=.env
if [ -f "$ENVFILE" ]; then
     export $(egrep -v '^#' .env | xargs) > /dev/null
fi

mkdir -p output
./bin/bazaar
