THIS := $(lastword $(MAKEFILE_LIST))

.PHONY: test

build:
	scripts/build.sh

run:
	scripts/run.sh

test:
	scripts/test.sh

go:
	@$(MAKE) -f $(THIS) build
	@$(MAKE) -f $(THIS) run

env:
	scripts/env.sh

sol:
	solc --abi assets/reserves.sol -o assets --overwrite
	abigen --abi=assets/UniswapView.abi --pkg=chain --out=pkg/chain/reserves.go
	go mod tidy
