#!/bin/sh

echo "Building bazaar.."

go build -o bin/bazaar cmd/bazaar/main.go

echo "Built to bin/bazaar"