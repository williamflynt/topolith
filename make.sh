#! /usr/bin/env sh

peg -inline -switch -strict -output pkg/grammar/grammar.peg.go pkg/grammar/grammar.peg && \
go vet fmt ./... && \
go test ./... && \
# Build for desktop CLI.
go build -o ./bin/topolith ./cmd/cli/main.go && \
chmod +x ./bin/topolith && \
# Build for web stack or other WASM compatible platforms.
GOOS=js GOARCH=wasm go build -o ./bin/topolith.wasm ./cmd/wasm && \
# Copy files for our web application.
cp "$GOROOT/misc/wasm/wasm_exec.js" ./web/public/ && \
cp ./bin/topolith.wasm ./web/public/ && \
# Build the web application.
cd ./web && npm run build && cd ..
