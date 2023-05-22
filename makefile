.DEFAULT_GOAL = all

numcpus  := $(shell cat /proc/cpuinfo | grep '^processor\s*:' | wc -l)
version  := $(shell git rev-list --count HEAD).$(shell git rev-parse --short HEAD)

name     := parse
package  := github.com/corpix/$(name)

.PHONY: all
all:

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: bench
bench:
	go test -bench=. -v ./...

.PHONY: lint
lint:
	golangci-lint --color=always --timeout=120s run ./...

.PHONY: doc
doc:
	gomarkdoc --include-unexported -o doc/doc.md .
