generate:
	@go generate ./...
.PHONY: generate

install:
	@go get -u -v github.com/golang/dep/cmd/dep
	@dep ensure
.PHONY: install

ci:
	@go get -u -v github.com/golang/dep/cmd/dep
	@dep ensure -vendor-only
.PHONY: ci

build: generate
	@echo "====> Build note cli"
	@go build -o ./bin/note main.go
.PHONY: build

local: build
	@echo "====> Add to GOPATH"
	@cp ./bin/note ${GOPATH}/bin
.PHONY: local

release:
	@echo "====> Build and release"
	@go get github.com/goreleaser/goreleaser
	@goreleaser
.PHONY: release

test:
	@go test ./...
.PHONY: test

test.cov:
	@go test ./... -coverprofile=coverage.txt -covermode=atomic
.PHONY: test.cov

docs:
	@echo "====> Build docs"
	@vuepress build docs
.PHONY: docs

docs.dev:
	@echo "====> Dev docs"
	@vuepress dev docs
.PHONY: docs.dev

deploydocs: docs
	@echo "====> Deploy docs to netlify"
	@netlify deploy -p docs/.vuepress/dist
.PHONY: deploydocs
