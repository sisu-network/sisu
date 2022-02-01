# Use cluster name sisu by default, but accept overrides.
CLUSTER_NAME ?= sisu

# The first instance of sisu deployed into a kind cluster should just be named sisu
# for predictability, but subsequent instances should have a unique name.
SISU_NAMESPACE ?= sisu

ifneq ($(ASSUME_CLUSTER_EXISTS),)
# Adjust the namespace for a secondary sisu instance.
SISU_NAMESPACE := $(SISU_NAMESPACE)-$(shell openssl rand -hex 6)
$(info will use sisu namespace $(SISU_NAMESPACE))

# Assume the cluster is ready, so validate the exact cluster name is present in kind get clusters.
.PHONY: kind-cluster-ready
kind-cluster-ready:
	kind get clusters | grep -qx "$(CLUSTER_NAME)"

else

# Not assuming the cluster exists, so it is ready after we create it.
.PHONY: kind-cluster-ready
kind-cluster-ready: create-kind-cluster ;

endif

# Allow overriding the version of go (for go mod vendor).
GO ?= go

# Create a kind cluster, and build, load, and deploy all required images.
.PHONY: init-cluster
init-cluster: create-kind-cluster create-clusterrole deploy-ganache deploy-mysql deploy-sisu
	$(info Access the cluster's sisu instance with:)
	$(info $(SPACE) $(SPACE) kubectl --context "kind-$(CLUSTER_NAME)" -n "$(SISU_NAMESPACE)")

# Delete the cluster.
.PHONY: delete-cluster
delete-cluster: require-kind
	kind delete cluster --name "$(CLUSTER_NAME)"

# Just create the cluster without loading anything.
.PHONY: create-kind-cluster
create-kind-cluster: require-kind
	kind create cluster --name $(CLUSTER_NAME)

# Ensure we have the kind command; if we don't, the error message should be clear enough.
.PHONY: require-kind
require-kind:
	kind --version

.PHONY: require-jq
require-jq:
	jq --version

.PHONY: go-mod-vendor
go-mod-vendor:
	$(GO) mod vendor

.PHONY: create-clusterrole
create-clusterrole: create-kind-cluster
	kubectl --context "kind-$(CLUSTER_NAME)" apply -f ./namespace-lister-clusterrole.yml

########### GANACHE ##########
# There is only one ganache statefulset in a kind cluster, as it represents an external resource.
# It is deployed into a fixed "ganache" namespace.

# Build the image for ganache (there wasn't an official image for ganache 7.0 alpha at time of writing).
.PHONY: build-docker-ganache
build-docker-ganache:
	docker build --quiet -t sisu.test/ganache:latest -f ./Dockerfile.ganache /dev/null

# Load the ganache image into the kind cluster, after the cluster is ready and we have built the image.
.PHONY: deploy-ganache
deploy-ganache: kind-cluster-ready build-docker-ganache
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/ganache:latest
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace ganache
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace ganache apply -f ./ganache.yml
	@# Block until ganache statefulset is ready.
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace ganache rollout status --timeout=30s statefulset/ganache

########## MySQL ##########
# There is a single MySQL instance in the kind cluster, to be shared by all sisu instances.
# It is deployed into a fixed "mysql" namespace.

# Ensure the docker image for mysql:8 is present on the host machine.
.PHONY: prepare-mysql
prepare-mysql:
	test -n "$$(docker image ls --quiet mysql:8)" || docker image pull --quiet mysql:8

# Deploy mysql once the cluster is ready.
# The mysql image isn't small, so we pre-load it in to the cluster to save a large network transfer.
.PHONY: deploy-mysql
deploy-mysql: kind-cluster-ready prepare-mysql
	kind load docker-image --name "$(CLUSTER_NAME)" mysql:8
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace mysql
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace mysql apply -f ./mysql.yml
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace mysql rollout status --timeout=45s deployment/mysql

########## SISU ##########
# Deploying one "instance" of sisu deploys deyes, dheart, and sisud to its own namespace.

# Meta-target to deploy all the individual services.
.PHONY: deploy-sisu
deploy-sisu: deploy-deyes deploy-dheart deploy-sisud
	$(info Added sisu instance into namespace $(SISU_NAMESPACE))

# Build the all-in-one sisu docker image.
# TODO: the "go install" step of this takes about two minutes on my machine.
#       It should be possible to optimize that with a buildkit volume for a Go build cache.
.PHONY: build-docker-sisu
build-docker-sisu: go-mod-vendor
	cd .. && docker build --quiet -t sisu.test/sisu-all:latest -f ./kind/Dockerfile.sisu-all .

# Load sisu all-in-one image after the cluster is ready and the sisu all-in-one image is built.
.PHONY: load-sisu
load-sisu: kind-cluster-ready build-docker-sisu
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/sisu-all:latest
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace "$(SISU_NAMESPACE)"
	@# Add arbitrary label so that sisu instances can discover other sisu instances in the same kind cluster.
	kubectl --context "kind-$(CLUSTER_NAME)" label namespace  "$(SISU_NAMESPACE)" sisu.network/purpose=test

.PHONY: deploy-deyes
deploy-deyes: load-sisu
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace "$(SISU_NAMESPACE)" apply -f ./deyes.yml
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace "$(SISU_NAMESPACE)" rollout status --timeout=30s deployment/deyes

.PHONY: deploy-dheart
deploy-dheart: load-sisu
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace "$(SISU_NAMESPACE)" apply -f ./dheart.yml
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace "$(SISU_NAMESPACE)" rollout status --timeout=30s deployment/dheart

.PHONY: deploy-sisud
deploy-sisud: load-sisu load-bitnami-kubectl require-jq
	@# Patch the clusterrolebinding so the init container can list namespaces.
	@# We have to include all the namespaces every time we patch it.
	kubectl --context "kind-$(CLUSTER_NAME)" patch clusterrolebinding namespace-lister --patch "$$( \
		kubectl --context "kind-$(CLUSTER_NAME)" get ns -l sisu.network/purpose=test -oname | \
		jq -R '{kind: "ServiceAccount", name: "namespace-lister", namespace: (. | ltrimstr("namespace/"))}' | \
		jq -c --slurp '{subjects: .}' \
	)"
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace "$(SISU_NAMESPACE)" apply -f ./sisud.yml
	@# TODO: sisud will fail currently because it has incorrect configuration.
	@#kubectl --context "kind-$(CLUSTER_NAME)" --namespace "$(SISU_NAMESPACE)" rollout status --timeout=30s deployment/sisud

# The bitnami/kubectl image is used for the init container on sisud to discover the other sisu instances in the kind cluster.
.PHONY: prepare-bitnami-kubectl
prepare-bitnami-kubectl:
	test -n "$$(docker image ls --quiet bitnami/kubectl:1.22)" || docker image pull --quiet bitnami/kubectl:1.22

.PHONY: load-bitnami-kubectl
load-bitnami-kubectl: kind-cluster-ready prepare-bitnami-kubectl
	kind load docker-image --name "$(CLUSTER_NAME)" bitnami/kubectl:1.22
