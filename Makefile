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

.PHONY: dev-mysql-up
dev-mysql-up:
	@docker compose \
		-f docker/docker-compose-mysql.dev.yml \
		up -d

.PHONY: dev-mysql-down
dev-mysql-down:
	@docker compose \
		-f docker/docker-compose-mysql.dev.yml \
		down -v --rmi local

# Create a local kind cluster with everything needed for sisu.
.PHONY: create-kind-cluster
create-kind-cluster:
	@cd ./kind && $(MAKE) -j create-cluster

# Delete the cluster you made with create-kind-cluster.
.PHONY: delete-kind-cluster
delete-kind-cluster:
	@cd ./kind && $(MAKE) -j delete-cluster

# Set global git configuration to only replace our private dependencies with the SSH URL.
.PHONY: configure-git
configure-git:
	git config --global url."git@github.com:sisu-network/deyes".insteadOf 'https://github.com/sisu-network/deyes'
	git config --global url."git@github.com:sisu-network/dheart".insteadOf 'https://github.com/sisu-network/dheart'
	git config --global url."git@github.com:sisu-network/lib".insteadOf 'https://github.com/sisu-network/lib'
	git config --global url."git@github.com:sisu-network/tss-lib".insteadOf 'https://github.com/sisu-network/tss-lib'