#!/bin/bash

echo "Downloading Go packages..."
go mod download

echo "Building..."
go build -v ./...

echo "Installing..."
# Install the Go program
go install ./...

echo "Done."