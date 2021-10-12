
.DEFAULT_GOAL:=help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run: ## Run the generate-package-repository program for the specified channel
ifeq ($(origin CHANNEL),undefined && $(origin TAG),undefined)
	@echo "Error! CHANNEL or TAG env var not set"
else
	go run create-package-repo.go $(CHANNEL) $(TAG)
endif

validate-metadata-cr : ## Validate metadata.yml
ifeq ($(origin METADATAFILE),undefined)
	@echo "Error! METADATAFILE env var not set"
else
	go run scripts/cmd/metadata/main.go $(METADATAFILE)
endif

validate-package-cr : ## Validate package.yml
ifeq ($(origin PACKAGEFILE),undefined)
	@echo "Error! PACKAGEFILE env var not set"
else
	go run scripts/cmd/package/main.go $(PACKAGEFILE)
endif