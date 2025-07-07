#!/bin/bash
go build -o bin/gobbench main.go
GOOS=linux GOARCH=s390x go build -o bin/gobbench-s390x main.go