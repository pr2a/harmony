#!/bin/bash

export GO111MODULE=on
ACTION=${1:-build}

function build_only
{
   VERSION=$(git rev-list --count HEAD)
   COMMIT=$(git describe --always --long --dirty)
   BUILTAT=$(date +%FT%T%z)
   BUILTBY=${USER}@

   go build -ldflags="-X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILTAT} -X main.builtBy=${BUILTBY}" -o hexie main.go
}

case "$ACTION" in
   "build")
      build_only
      ./hexie -version
      ;;
   "test")
      dev_appserver.py app.yaml
      ;;
esac

sleep 1
echo killing dev_appserver.py
pkill python
