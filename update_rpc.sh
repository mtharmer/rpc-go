#!/bin/sh
echo "Updating rpc files..."
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./rpcgo/rpcgo.proto
echo "Done!"
