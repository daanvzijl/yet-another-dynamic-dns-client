#!/usr/bin/env -S just --justfile

default:
    just --list

install:
    hk install

lint:
    hk check --all

fix:
    hk fix --all

test:
    go test ./...

test-race:
    go test -race ./...

build:
    go build -o yaddc ./main.go

run:
    go run ./main.go

docker-build:
    docker build -t yaddc .
