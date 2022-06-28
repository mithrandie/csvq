GOPATH := $(shell pwd)/build
BINARY := csvq
RELEASE_ARCH := darwin/amd64 darwin/arm64 linux/amd64 linux/386 linux/arm linux/arm64 freebsd/amd64 freebsd/386 freebsd/arm netbsd/amd64 netbsd/386 netbsd/arm openbsd/amd64 openbsd/386 windows/amd64 windows/386
PRERELEASE_ARCH := darwin/amd64 darwin/arm64 linux/amd64 windows/amd64
BUILD_TAGS := -tags urfave_cli_no_docs

ifneq ($(shell command -v git && git remote -v 2>/dev/null | grep mithrandie/csvq.git && echo true),true)
	VERSION := $(shell git describe --tags --always 2>/dev/null)
endif

ifdef VERSION
	LDFLAGS := -ldflags="-X github.com/mithrandie/csvq/lib/query.Version=$(VERSION) -s -w"
endif

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL: $(BINARY)

$(BINARY): build

.PHONY: build
build:
	go build $(BUILD_TAGS) -trimpath $(LDFLAGS) -o $(GOPATH)/bin/

.PHONY: install
install:
	go install $(BUILD_TAGS) -trimpath $(LDFLAGS)

.PHONY: clean
clean:
	go clean -i -cache -modcache

.PHONY: build-all
build-all:
	IFS='/'; \
	for TARGET in $(RELEASE_ARCH); \
	do \
		set -- $$TARGET; \
		GOOS=$$1 GOARCH=$$2 go build $(BUILD_TAGS) -trimpath $(LDFLAGS) -o "dist/$(BINARY)-$(VERSION)-$${1}-$${2}/"; \
	done

.PHONY: build-pre-release
build-pre-release:
	IFS='/'; \
	for TARGET in $(PRERELEASE_ARCH); \
	do \
		set -- $$TARGET; \
		GOOS=$$1 GOARCH=$$2 go build $(BUILD_TAGS) -trimpath $(LDFLAGS) -o "dist/$(BINARY)-$(VERSION)-$${1}-$${2}/"; \
	done

.PHONY: dist -
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../BINARY_CODE_LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) cp ../CHANGELOG.md {} \; && \
	$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
	cd ..

.PHONY: release
release:
ifeq ($(shell git tag --contains 2>/dev/null),)
	$(error HEAD commit is not tagged)
else
	git push origin $(VERSION)
endif

.PHONY: install-goyacc
install-goyacc:
ifeq ($(shell command -v goyacc 2>/dev/null),)
	go get github.com/cznic/goyacc
endif

.PHONY: yacc
yacc: install-goyacc
	cd lib/parser && \
	goyacc -o parser.go -v parser.output parser.y && \
	cd ../../lib/json && \
	goyacc -p jq -o query_parser.go -v query_parser.output query_parser.y && \
	goyacc -p jp -o path_parser.go -v path_parser.output path_parser.y && \
	cd ../..

