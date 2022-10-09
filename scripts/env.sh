#!/bin/sh

echo "Please enter your INFURA API key:"
read input
echo "BAZAAR_INFURA_KEY=$input" > .env

echo "Do you wish to exclude dead tokens? 1 or 0"
read input
echo "BAZAAR_EXCLUDE_DEAD_TOKENS=$input" >> .env

echo "Do you wish to include fees when finding potential arbitrage paths? 1 or 0"
read input
echo "BAZAAR_INCLUDE_FEES=$input" >> .env
