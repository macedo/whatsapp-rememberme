TARGET_FILE:=${shell head -n1 go.mod | sed -r 's/.*\/(.*)/\1/g' }
BUILD_DIR=.build

.PHONY: target download install-tools

target: build-app

download: ## Download go.mod dependencies
	echo Download go.mod dependencies
	go mod download

install-tools: download ## Install tools
	echo Installing tools from tools/tools.go
	go list -f '{{range .Imports}}{{.}} {{end}}' tools/tools.go | xargs go install
