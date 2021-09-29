#!/usr/bin/env bash

echo "128 bytes, count 1"
go run ./cmd -s 128 -c 1
echo
echo "512 bytes, count 10"
go run ./cmd -s 512 -c 10
echo
echo "1024 bytes, count 100"
go run ./cmd -s 1024 -c 100
echo
echo "4096 bytes, count 1000"
go run ./cmd -s 4096 -c 1000



