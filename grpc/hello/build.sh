#!/bin/sh
set -e
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  idl/greet.proto
go build -o bin/server cmd/server/main.go
go build -o bin/client cmd/client/main.go
exit 0