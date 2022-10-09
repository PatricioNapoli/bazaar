THIS := $(lastword $(MAKEFILE_LIST))

build:
	scripts/build.sh

run:
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
