#!/bin/bash

go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
go tool cover -func coverage.out | grep total
