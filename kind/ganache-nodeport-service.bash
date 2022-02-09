#!/bin/bash

set -eu -o pipefail

jq -n --arg index "$GANACHE_INDEX" '{
  apiVersion: "v1",
  kind: "Service",
  metadata: {
    name: @text "ganache-\($index)-nodeport",
  },
  spec: {
    type: "NodePort",
    selector: {app: "ganache", "statefulset.kubernetes.io/pod-name": @text "ganache-\($index)"},
    ports: [
      {
	port: 7545,
	nodePort: (32000+($index | tonumber)),
      }
    ]
  },
}'
