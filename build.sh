#!/bin/bash

# Store influxdata gopath
idp=$GOPATH/src/github.com/influxdata

# Create build directory
mkdir -p build/ > /dev/null 2>&1

# Build 0.11 binaries
echo "Restoring 0.11 dependencies..."
INFLUX_VERSION="0.11"
cd $idp/influxdb
git co v0.11.1 > /dev/null 2>&1
gdm restore > /dev/null 2>&1
cd $idp/schema-shape
echo "Building 0.11 binaries..."
go build -o build/schema-shape_darwin_0.11
GOOS=linux GOARCH=amd64 go build -o build/schema-shape_linux_amd64_0.11

# Build 0.12 binaries
echo "Restoring 0.12 dependencies..."
INFLUX_VERSION="0.12"
cd $idp/influxdb
git co v0.12.2 > /dev/null 2>&1
gdm restore > /dev/null 2>&1
cd $idp/schema-shape
echo "Building 0.12 binaries..."
go build -o build/schema-shape_darwin_0.12
GOOS=linux GOARCH=amd64 go build -o build/schema-shape_linux_amd64_0.12

# Build 0.13 binaries
echo "Restoring 0.13 dependencies..."
INFLUX_VERSION="0.13"
cd $idp/influxdb
git co v0.13.0 > /dev/null 2>&1
gdm restore > /dev/null 2>&1
cd $idp/schema-shape
echo "Building 0.13 binaries..."
go build -o build/schema-shape_darwin_0.13
GOOS=linux GOARCH=amd64 go build -o build/schema-shape_linux_amd64_0.13

# Build 1.0 binaries
echo "Restoring 1.0 dependencies..."
INFLUX_VERSION="1.0"
cd $idp/influxdb
git co v1.0.0-beta1 > /dev/null 2>&1
gdm restore > /dev/null 2>&1
cd $idp/schema-shape
echo "Building 1.0 binaries..."
go build -o build/schema-shape_darwin_1.0
GOOS=linux GOARCH=amd64 go build -o build/schema-shape_linux_amd64_1.0