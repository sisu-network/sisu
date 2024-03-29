# Use cluster name sisu by default, but accept overrides.
CLUSTER_NAME ?= sisu

# Not assuming the cluster exists, so it is ready after we create it.
.PHONY: kind-cluster-ready
kind-cluster-ready: create-kind-cluster ;

# Allow overriding the version of go (for go mod vendor).
GO ?= go

# Create a kind cluster, and build, load, and deploy all required images.
.PHONY: init-cluster
init-cluster: create-kind-cluster ready-config-server ready-ganache ready-mysql
	$(info Access the cluster's sisu instance with:)
	$(info $(SPACE) $(SPACE) kubectl --context "kind-$(CLUSTER_NAME)" -n sisu-0)

ifeq ($(SISU_COUNT),)
$(warn $$SISU_COUNT not set; defaulting to 2 instances)
SISU_COUNT := 2
endif

# init-cluster depends on deploy-sisu-0 ... deploy-sisu-$(SISU_COUNT - 1)
init-cluster: $(foreach n,$(shell seq 0 $$(( $(SISU_COUNT) - 1)) ),ready-sisu-$(n))

# Delete the cluster.
.PHONY: delete-cluster
delete-cluster: require-kind
	kind delete cluster --name "$(CLUSTER_NAME)"

# Just create the cluster without loading anything.
# This configures the single kind node to expose 7545 and 7546 as node ports for the two ganache instances,
# and 25456+ as node ports for sisu-0 and up.
.PHONY: create-kind-cluster
create-kind-cluster: require-kind require-jq
	jq -n '{ \
	  kind: "Cluster", \
	  apiVersion: "kind.x-k8s.io/v1alpha4", \
	  nodes: [ \
	    { \
	      role: "control-plane", \
	      extraPortMappings: ( \
	        [range(2) | {listenAddress: "127.0.0.1", containerPort:(32000+.), hostPort:(7545+.)}] + \
	        [range($(SISU_COUNT)) | {listenAddress: "127.0.0.1", containerPort:(30000+.), hostPort:(25456+.)}] \
	      ) \
	    } \
	  ] \
	}' | \
	    kind create cluster --name $(CLUSTER_NAME) --config -


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

########### GANACHE ##########
# There is only one ganache statefulset in a kind cluster, as it represents an external resource.
# It is deployed into a fixed "ganache" namespace.

# Build the image for ganache (there wasn't an official image for ganache 7.0 alpha at time of writing).
.PHONY: build-docker-ganache
build-docker-ganache:
	docker build --quiet -t sisu.test/ganache:latest -f ./Dockerfile.ganache /dev/null

# Load the ganache image into the kind cluster, after the cluster is ready and we have built the image.
.PHONY: deploy-ganache
deploy-ganache: kind-cluster-ready build-docker-ganache require-jq
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/ganache:latest
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace ganache
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace ganache apply -f ./ganache.yml
	@# There is no obvious way to kubectl apply a multidocument JSON, so this is two separate calls.
	GANACHE_INDEX=0 ./ganache-nodeport-service.bash | \
		    kubectl --context "kind-$(CLUSTER_NAME)" --namespace ganache apply -f -
	GANACHE_INDEX=1 ./ganache-nodeport-service.bash | \
		    kubectl --context "kind-$(CLUSTER_NAME)" --namespace ganache apply -f -

.PHONY: ready-ganache
ready-ganache: deploy-ganache
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

.PHONY: ready-mysql
ready-mysql: deploy-mysql
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace mysql rollout status --timeout=45s deployment/mysql

########## SISU ##########
# Deploying one "instance" of sisu deploys deyes, dheart, and sisud to its own namespace.

# Meta-target to deploy all the individual services, using pattern matching.
.PHONY: ready-sisu-%
ready-sisu-%: ready-deyes-% ready-dheart-% ready-sisud-%
	$(info Deployed sisu instance into namespace $(@:ready-%=%))

# Build the all-in-one sisu docker image.
.PHONY: build-docker-sisu
build-docker-sisu: go-mod-vendor
	cd .. && docker build --quiet -t sisu.test/sisu-all:latest -f ./kind/Dockerfile.sisu-all .

# Load the sisu image into the kind cluster, after it has been built and after the kind cluster is ready.
.PHONY: preload-sisu-image
preload-sisu-image: kind-cluster-ready build-docker-sisu
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/sisu-all:latest

# Create the namespace for the sisu instance,
# and create the configmap that reports the total node count and the index of this namespace.
.PHONY: create-sisu-ns-%
create-sisu-ns-%: kind-cluster-ready preload-sisu-image
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace sisu-$(@:create-sisu-ns-%=%)
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:create-sisu-ns-%=%) create configmap sisu-cluster \
		--from-literal node-count=$(SISU_COUNT) --from-literal index=$(@:create-sisu-ns-%=%)

.PHONY: deploy-deyes-%
deploy-deyes-%: create-sisu-ns-% deploy-config-server deploy-mysql
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:deploy-deyes-%=%) apply -f ./deyes.yml

.PHONY: ready-deyes-%
ready-deyes-%: deploy-deyes-% ready-mysql ready-config-server
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:ready-deyes-%=%) rollout status --timeout=60s deployment/deyes

.PHONY: deploy-dheart-%
deploy-dheart-%: create-sisu-ns-%
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:deploy-dheart-%=%) apply -f ./dheart.yml

.PHONY: ready-dheart-%
ready-dheart-%: deploy-dheart-% ready-config-server
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:ready-dheart-%=%) rollout status --timeout=30s deployment/dheart

.PHONY: deploy-sisud-%
deploy-sisud-%: create-sisu-ns-% require-jq
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:deploy-sisud-%=%) apply -f ./sisud.yml
	SISU_INDEX=$(@:deploy-sisud-%=%) ./sisud-nodeport-service.bash | \
		    kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:deploy-sisud-%=%) apply -f -

.PHONY: ready-sisud-%
ready-sisud-%: deploy-sisud-% ready-config-server
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace sisu-$(@:ready-sisud-%=%) rollout status --timeout=30s deployment/sisud

########## config-server ##########
.PHONY: build-docker-config-server
build-docker-config-server:
	docker build --quiet -t sisu.test/config-server:latest -f ./Dockerfile.config-server /dev/null

.PHONY: load-config-server
load-config-server: kind-cluster-ready build-docker-config-server
	kind load docker-image --name "$(CLUSTER_NAME)" sisu.test/config-server:latest

.PHONY: deploy-config-server
deploy-config-server: load-config-server preload-sisu-image
	kubectl --context "kind-$(CLUSTER_NAME)" create namespace config-server
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace config-server create configmap sisu-cluster \
		--from-literal node-count=$(SISU_COUNT)
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace config-server create -f ./config-server.yml

.PHONY: ready-config-server
ready-config-server: deploy-config-server
	kubectl --context "kind-$(CLUSTER_NAME)" --namespace config-server rollout status --timeout=30s statefulset/config-server
