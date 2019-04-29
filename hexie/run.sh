#!/bin/bash

go build -o hexie main.go
dev_appserver.py app.yaml &

sleep 5

curl http://localhost:8080/enter
curl http://localhost:8080/finish

sleep 1
echo killing dev_appserver.py
pkill python
