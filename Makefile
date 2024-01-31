# See https://golangci-lint.run/usage/install/
LINTER_VERSION = v1.55.2

# Variables needed when building binaries
VERSION := $(shell grep -oE -m 1 '([0-9]+)\.([0-9]+)\.([0-9]+)' CHANGELOG.md )

# To be used for dependencies not installed with gomod
LOCAL_DEPS_INSTALL_LOCATION = /usr/local/bin

# Target for Pull Request checks
.PHONY: pr-checks
pr-checks: clean build unit-test lint coverage-total

.PHONY: clean
clean:
	rm -rf build
	mkdir -p build

.PHONY: deps
deps:
	go env -w "GOPRIVATE=github.com/ildomm/*"
	go mod download

.PHONY: build
build: deps build-server

.PHONY: build-server
build-server: deps
	# Build the http server binary
	go build -ldflags="-X main.semVer=${VERSION}" \
        -o build/http_server

.PHONY: unit-test
unit-test: deps
	go test -tags=testing -count=1 ./...

.PHONY: lint-install
lint-install:
	[ -e ${LOCAL_DEPS_INSTALL_LOCATION}/golangci-lint ] || \
	wget -O- -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b ${LOCAL_DEPS_INSTALL_LOCATION} ${LINTER_VERSION}

.PHONY: lint
lint: deps lint-install
	golangci-lint run

.PHONY: coverage-report
coverage-report: clean deps
	go test -tags=testing ./... \
		-coverprofile=build/cover.out github.com/ildomm/ssccg/...
	go tool cover -html=build/cover.out -o build/coverage.html
	echo "** Coverage is available in build/coverage.html **"

.PHONY: coverage-total
coverage-total: clean deps
	go test -tags=testing ./... \
		-coverprofile=build/cover.out github.com/ildomm/ssccg/...
	go tool cover -func=build/cover.out | grep total