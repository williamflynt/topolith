#! /usr/bin/env sh

go vet fmt ./... && \
go test ./... && \
peg -inline -switch -strict -output pkg/grammar/grammar.peg.go pkg/grammar/grammar.peg && \
go build -o ./bin/topolith ./cmd/repl/main.go
