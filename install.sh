#!/usr/bin/env sh

prefix=${HOME}/.local
mkdir -p ${prefix}/bin
go build -o ${prefix}/bin/remotectlc cmd/remotectlc/main.go
go build -o ${prefix}/bin/remotectld cmd/remotectld/main.go
