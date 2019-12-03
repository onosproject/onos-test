#!/bin/sh

proto_imports="./pkg:./test:./benchmark:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src"

protoc -I=$proto_imports --gogofaster_out=import_path=onos/benchmark/grpc,plugins=grpc:benchmark benchmark/grpc/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=onos/benchmark,plugins=grpc:pkg pkg/benchmark/*.proto