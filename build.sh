#!/bin/bash

GOOS=darwin GOARCH=arm64 go build -o kubectl-cleaner main.go

tar czf kubectl-cleaner_darwin_arm64.tar.gz kubectl-cleaner
