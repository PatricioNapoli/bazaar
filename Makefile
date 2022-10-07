THIS := $(lastword $(MAKEFILE_LIST))

build:
	scripts/build.sh

run:
	scripts/run.sh

go:
	@$(MAKE) -f $(THIS) build
	@$(MAKE) -f $(THIS) run
