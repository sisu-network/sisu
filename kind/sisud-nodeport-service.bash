#!/bin/bash

set -eu -o pipefail

jq -n --arg index "$SISU_INDEX" '{
  apiVersion: "v1",
  kind: "Service",
  metadata: {
    name: "sisud-nodeport",
  },
  spec: {
    type: "NodePort",
    selector: {app: "sisud"},
    ports: [
      {
	port: 25456,
        nodePort: (30000 + ($index | tonumber))
      }
    ]
  },
}'
