#!/bin/sh

echo "Please enter your INFURA API key:"
read -s input
echo "BAZAAR_INFURA_KEY=$input" > .env

echo "Do you wish to exclude dead tokens? (default is 1)"
read input
input=${input:-1}
echo "BAZAAR_EXCLUDE_DEAD_TOKENS=$input" >> .env

echo "Do you wish to include fees when finding potential arbitrage paths? (default is 0)"
read input
input=${input:-0}
echo "BAZAAR_INCLUDE_FEES=$input" >> .env

echo "Do you wish to include gas when finding potential arbitrage paths? (default is 0)"
read input
input=${input:-0}
echo "BAZAAR_INCLUDE_GAS=$input" >> .env

echo "Please enter the output dir+filename: (default is output/output.json)"
read input
input=${input:-output/output.json}
echo "BAZAAR_OUTPUT_FILENAME=$input" >> .env
