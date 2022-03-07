help: ## Show list of make targets and their description
	grep -E '^[/%.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
      | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint:
	./scripts/run_lint.sh

.PHONY: test
test:
	./scripts/run_test.sh

.PHONY: fmt
fmt:
	find . -name "*.go" | grep -v -E "(.*/proto/.*|./*/mock/.*)" | xargs -I '{}' gofmt -s -w '{}'

# Create a local kind cluster with everything needed for sisu.
.PHONY: kind-init-cluster
kind-init-cluster:
	@cd ./kind && $(MAKE) -j init-cluster

# Delete the cluster you made with kind-init-cluster.
.PHONY: kind-delete-cluster
kind-delete-cluster:
	@cd ./kind && $(MAKE) -j delete-cluster
