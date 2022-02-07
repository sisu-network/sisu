// This file is a Go convention for tracking tooling dependencies:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
//
// It allows us to depend on the main packages of dheart and deyes,
// so that we can easily "go run" or "go build" those tools during local development.

//go:build tools
// +build tools

package tools

import (
	_ "github.com/sisu-network/dheart"
	_ "github.com/sisu-network/deyes"
)
