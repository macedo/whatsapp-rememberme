TARGET_FILE:=${shell head -n1 go.mod | sed -r 's/.*\/(.*)/\1/g' }
BUILD_DIR=.build

.PHONY: target clear download install-tools

target: dev

clear:
	rm -rf ./esbuild ./server

dev: clear
	go build -o esbuild cmd/esbuild/*.go
	go build -o server cmd/web/*.go
	./esbuild &./server

download: ## Download go.mod dependencies
	echo Download go.mod dependencies
	go mod download

install-tools: download ## Install tools
	echo Installing tools from tools/tools.go
	go list -f '{{range .Imports}}{{.}} {{end}}' tools/tools.go | xargs go install
