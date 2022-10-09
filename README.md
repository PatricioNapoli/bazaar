
<p align="center">
    <img alt="Grand Bazaar" src="assets/bazaar.jpg" width="400px"/>
</p>

<div align="center">

  <a style="margin-right:15px" href="#"><img src="https://forthebadge.com/images/badges/made-with-go.svg" alt="Made with Go"/></a>
  <a style="margin-right:15px" href="#"><img src="https://forthebadge.com/images/badges/powered-by-black-magic.svg" alt="Made with Go"/></a>
  <a href="https://www.paradigm.xyz/2020/08/ethereum-is-a-dark-forest"><img src="assets/dark-forest.svg" alt="Ethereum is a dark forest"/></a>


  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-brightgreen.svg" alt="License MIT"/></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/go-1.18-blue.svg" alt="Go 1.18"/></a>
  <a href="https://github.com/PatricioNapoli/bazaar/actions/workflows/build.yml"><img src="https://github.com/PatricioNapoli/bazaar/actions/workflows/build.yml/badge.svg" alt="build"/></a>

</div>


# Bazaar

## Overview

Go console application that searches for possible market arbitration paths utilizing pre-approved coin pair swapping paths.  
Possible markets are two, UniswapV2 and SushiSwap.  
Always starts at 1 WETH and attempts to exit with same coin.

## Prerequisites

If you are not using docker, these are required:  

make  
go 1.18+  
solc  
abigen (geth)  

## Environment

For getting a prompt for the required INFURA key and other configs, and creating .env file.  

`make env`  

Alternatively, you may export manually:  

`export BAZAAR_INFURA_KEY=YOURKEY`

## Build & Run

### Make

`make go`  

Alternatively, you may run `make build` and `make run` separately.  
Or run the scripts in `scripts/`.  

Contract ABI for pair reserves is already built but can be recompiled with:  

`make sol`

### Docker

#### Building

`docker build -t patricionapoli/bazaar .` 

#### Running 

First, either create the env file through:

`make env`  

And then:  

`docker run --env-file .env -v "$(pwd)"/output:/go/src/bazaar/output patricionapoli/bazaar`  

Or set the infura API key directly:  

`docker run -e BAZAAR_INFURA_KEY=key -v "$(pwd)"/output:/go/src/bazaar/output patricionapoli/bazaar`  

Please note that the output folder is being mapped to `/output`

#### Running release

`docker run --env-file .env -v "$(pwd)"/output:/go/src/bazaar/output ghcr.io/patricionapoli/bazaar:master`