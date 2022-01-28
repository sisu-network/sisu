# Use cluster name sisu by default, but accept overrides.
CLUSTER_NAME ?= sisu

# Allow overriding the version of go (for go mod vendor).
GO ?= go

# Create a kind cluster, and build, load, and deploy all required images.
.PHONY: create-cluster
create-cluster: create-kind-cluster deploy
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

.PHONY: go-mod-vendor
go-mod-vendor:
	$(GO) mod vendor

# Load all the docker images into the kind cluster.
.PHONY: deploy
deploy: deploy-ganache deploy-mysql deploy-deyes deploy-dheart ;

# Build the image for ganache (there wasn't an official image for ganache 7.0 alpha at time of writing).
.PHONY: build-docker-ganache
build-docker-ganache:
	docker build --quiet -t sisu.test/ganache:latest -f ./Dockerfile.ganache /dev/null

# Load the ganache image into the kind cluster, after the cluster is ready and we have built the image.
.PHONY: deploy-ganache
deploy-ganache: create-kind-cluster build-docker-ganache
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/ganache:latest
	kubectl --context "kind-$(CLUSTER_NAME)" apply -f ./ganache.yml
	@# Block until ganache statefulset is ready.
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu rollout status --timeout=30s statefulset/ganache

# Build the all-in-one sisu docker image.
# TODO: the "go install" step of this takes about two minutes on my machine.
#       It should be possible to optimize that with a buildkit volume for a Go build cache.
.PHONY: build-docker-sisu
build-docker-sisu: go-mod-vendor
	cd .. && docker build --quiet -t sisu.test/sisu-all:latest -f ./kind/Dockerfile.sisu-all .

# Load sisu all-in-one image after the cluster is ready and the sisu all-in-one image is built.
.PHONY: load-sisu
load-sisu: create-kind-cluster build-docker-sisu
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/sisu-all:latest

.PHONY: deploy-deyes
deploy-deyes: load-sisu
	kubectl --context "kind-$(CLUSTER_NAME)" apply -f ./deyes.yml
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu rollout status --timeout=30s deployment/deyes

.PHONY: deploy-dheart
deploy-dheart: load-sisu
	kubectl --context "kind-$(CLUSTER_NAME)" apply -f ./dheart.yml
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu rollout status --timeout=30s deployment/dheart

# Ensure the docker image for mysql:8 is present on the host machine.
.PHONY: prepare-mysql
prepare-mysql:
	test -n "$$(docker image ls --quiet mysql:8)" || docker image pull --quiet mysql:8

# Deploy mysql once the cluster is ready.
# The mysql image isn't small, so we pre-load it in to the cluster to save a large network transfer.
.PHONY: deploy-mysql
deploy-mysql: create-kind-cluster prepare-mysql
	kind load docker-image --name "$(CLUSTER_NAME)" mysql:8
	kubectl --context "kind-$(CLUSTER_NAME)" apply -f ./mysql.yml
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu rollout status --timeout=45s deployment/mysql
