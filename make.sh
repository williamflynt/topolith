#! /usr/bin/env sh

peg -inline -switch -strict -output pkg/grammar/grammar.peg.go pkg/grammar/grammar.peg && \
go vet fmt ./... && \
go test ./... && \
go build ./...
