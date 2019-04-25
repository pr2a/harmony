#!/usr/bin/env bash

set -x

export GO111MODULE=on

function build_only
{
   VERSION=$(git rev-list --count HEAD)
   COMMIT=$(git describe --always --long --dirty)
   BUILTAT=$(date +%FT%T%z)
   BUILTBY=${USER}@

   go build -ldflags="-X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILTAT} -X main.builtBy=${BUILTBY}" -o backend main.go
}

build_only

./backend -version

# ./backend
