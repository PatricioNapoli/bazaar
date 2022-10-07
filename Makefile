THIS := $(lastword $(MAKEFILE_LIST))

INFURAFILE := .infura
INFURAKEY :=$(file < $(INFURAFILE))

build:
	scripts/build.sh

run:
	export INFURA_KEY=$(INFURAKEY)
	scripts/run.sh

go:
	@$(MAKE) -f $(THIS) build
	@$(MAKE) -f $(THIS) run

env:
	scripts/env.sh

sol:
	solc --abi assets/reserves.sol -o assets --overwrite
	abigen --abi=assets/UniswapView.abi --pkg=chain --out=pkg/chain/reserves.go
	go mod tidy
