#!/bin/bash
# build the local thing
go build -o bin/gobbench main.go

# build for AWS Graviton2
GOOS=linux GOARCH=arm64 go build -o bin/gobbench-graviton2 main.go

# build for IBM z/Architecture (LinuxONE)
GOOS=linux GOARCH=s390x go build -o bin/gobbench-s390x main.go