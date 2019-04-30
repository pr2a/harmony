#!/bin/bash

export GO111MODULE=on

function build_only
{
   VERSION=$(git rev-list --count HEAD)
   COMMIT=$(git describe --always --long --dirty)
   BUILTAT=$(date +%FT%T%z)
   BUILTBY=${USER}@

   go build -ldflags="-X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILTAT} -X main.builtBy=${BUILTBY}" -o hexie main.go
}

build_only
./hexie -version

dev_appserver.py app.yaml &

sleep 5

curl http://localhost:8080/enter
curl http://localhost:8080/finish

sleep 1
echo killing dev_appserver.py
pkill python
