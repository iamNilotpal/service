# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Define dependencies
GOLANG          := golang:1.22
ALPINE          := alpine:3.20.2
POSTGRES        := postgres:15.7

APP             := sales
SERVICE_NAME    := sales-api
BASE_IMAGE_NAME := iamNilotpal/service
VERSION       	:= "$(shell git rev-parse --short HEAD)"
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)


# Install dependencies
dev-docker:
	docker pull $(GOLANG)
	docker pull $(ALPINE)
	docker pull $(POSTGRES)

# Building containers
all: service metrics

build-service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# Administration
migrate:
	go run apps/tooling/sales-admin/main.go migrate

seed: migrate
	go run apps/tooling/sales-admin/main.go seed

liveness:
	curl -il http://localhost:3000/v1/liveness

readiness:
	curl -il http://localhost:3000/v1/readiness

# Running tests within the local computer
test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-race lint vuln-check

# Modules support
deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
run-local:
	go run apps/services/sales/main.go

run-help:
	go run apps/services/sales/main.go --help
