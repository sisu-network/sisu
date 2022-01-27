# Use cluster name sisu by default, but accept overrides.
CLUSTER_NAME ?= sisu

# Create a kind cluster, and build, load, and deploy all required images.
.PHONY: create-cluster
create-cluster: create-kind-cluster load-docker-images
	$(info Access the cluster with:)
	$(info $(SPACE) $(SPACE) kubectl --context "kind-$(CLUSTER_NAME)" -n sisu)

# Delete the cluster.
.PHONY: delete-cluster
delete-cluster: require-kind
	kind delete cluster --name "$(CLUSTER_NAME)"

# Just create the cluster without loading anything.
.PHONY: create-kind-cluster
create-kind-cluster: require-kind
	kind create cluster --name $(CLUSTER_NAME)
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace sisu

# Ensure we have kind; if we don't, the error message should be clear enough.
.PHONY: require-kind
require-kind:
	kind --version

# Load all the docker images into the kind cluster.
.PHONY: load-docker-images
load-docker-images: load-docker-ganache ;

# Build the image for ganache (there wasn't an official image for ganache 7.0 alpha at time of writing).
.PHONY: build-docker-ganache
build-docker-ganache:
	docker build --quiet -t sisu.test/ganache:latest -f ./Dockerfile.ganache /dev/null

# Load the ganache image into the kind cluster, after the cluster is ready and we have built the image.
.PHONY: load-docker-ganache
load-docker-ganache: create-kind-cluster build-docker-ganache
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/ganache:latest
	kubectl --context "kind-$(CLUSTER_NAME)" apply -f ./ganache.yml
	@# Block until ganache statefulset is ready.
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu rollout status --timeout=30s statefulset/ganache
