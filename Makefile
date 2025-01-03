PKGS               = ./...
BINARY_NAME        = beanstalkd_exporter
DOCKER_IMAGE_NAME ?= beanstalkd_exporter
DOCKER_IMAGE_TAG  ?= $(shell git describe --tags --abbrev=0)

.PHONY: all
all: dep vet staticcheck lint clean test build

.PHONY: devtools
devtools:
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.7
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2

.PHONY: fmt
fmt:
	go fmt $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: staticcheck
staticcheck:
	staticcheck $(PKGS)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dep
dep:
	go mod download

.PHONY: test
test:
	go test $(PKGS)

.PHONY: build
build:
	go build -o $(BINARY_NAME) -v

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

.PHONY: docker
docker:
	docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .
