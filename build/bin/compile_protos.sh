#!/bin/sh

proto_imports="./test:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src"

protoc -I=$proto_imports --gogofaster_out=import_path=onos/test/grpc,plugins=grpc:test test/grpc/*.proto