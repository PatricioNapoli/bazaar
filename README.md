
<p align="center">
    <img alt="Grand Bazaar" src="assets/bazaar.jpg" width="400px"/>
</p>

<div align="center">

  <a style="margin-right:15px" href="#"><img src="https://forthebadge.com/images/badges/made-with-go.svg" alt="Made with Go"/></a>
  <a style="margin-right:15px" href="#"><img src="https://forthebadge.com/images/badges/powered-by-black-magic.svg" alt="Made with Go"/></a>
  <a href="https://www.paradigm.xyz/2020/08/ethereum-is-a-dark-forest"><img src="assets/dark-forest.svg" alt="Ethereum is a dark forest"/></a>


  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-brightgreen.svg" alt="License MIT"/></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/go-1.18-blue.svg" alt="Go 1.18"/></a>
</div>


# Bazaar

## Overview

Go console application that searches for possible market arbitration paths utilizing pre-approved coin pair swapping paths.  
Possible markets are two, UniswapV2 and SushiSwap.  
Always starts at 1 WETH and attempts to exit with same coin.

## Prerequisites

make  
go 1.18+  
solc  
abigen  
Infura API Key

## Dependencies

`go-ethereum`

## Environment

`make env`  

For getting a prompt for the INFURA key.  
Alternatively, you may export manually:

`export INFURA_KEY=YOURKEY`

## Build & Run

`make go`  

Alternatively, you may run `make build` and `make run` separately.  
Or run the scripts in `scripts/`.  

Contract ABI for pair reserves is already built but can be recompiled with:  

`make sol`