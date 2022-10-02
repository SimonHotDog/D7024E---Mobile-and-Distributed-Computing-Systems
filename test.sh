#!/bin/bash

go test -coverpkg=./... -coverprofile=coverage.out ./...

# Remove all internal test utilities that should not be tested
sed -i "/d7024e\/internal\/test/d" "coverage.out"

# Generate HTML report and display total root package code coverage
go tool cover -html=coverage.out -o coverage.html
go tool cover -func coverage.out | grep total
