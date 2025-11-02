root_mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
REPO_ROOT_DIR := $(realpath $(dir $(root_mkfile_path)))

export PROJECT_NAME ?= $(shell basename $(REPO_ROOT_DIR))

DIST_DIR ?= $(REPO_ROOT_DIR)/dist
SOURCES := $(wildcard $(REPO_ROOT_DIR)/cmd/*/main.go)
COMMANDS := $(patsubst $(REPO_ROOT_DIR)/cmd/%/main.go,%,$(SOURCES))
BINS := $(addprefix $(DIST_DIR)/,$(COMMANDS))

ALL_GO_FILES := $(shell find $(REPO_ROOT_DIR) -type f -name '*.go')
ALL_GO_TEST_FILES := $(shell find $(REPO_ROOT_DIR) -type f -name '*_test.go')

export GO_LD_FLAGS ?= -v -s -w -X 'main.Version=$(GIT_SHA)'

$(DIST_DIR):
	mkdir -p $(DIST_DIR)
	echo $(BINS)

$(DIST_DIR)/%: $(REPO_ROOT_DIR)/cmd/%/main.go $(DIST_DIR) $(ALL_GO_FILES)
	go build -ldflags "$(GO_LD_FLAGS)" -o $@ $(dir $<)*.go

#? build: Compile project to binary files
.PHONY: build
build: $(BINS)

#? clean: Clean built files
.PHONY: clean
clean:
	rm -rf $(DIST_DIR)
	rm -rf coverage.out coverage.xml coverage.html

#? test: Execute tests
.PHONY: test
test:
	go test -v -race -coverprofile coverage.out $(REPO_ROOT_DIR)/...

#? lint: Execute golang linter
.PHONY: lint
lint:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.64.2 golangci-lint run --config=./.github/linters/.golangci.yml --timeout=5m

#? docker-server: Build server docker image
.PHONY: docker-server
docker-server:
	docker build -t svalasovich/server -f ./build/package/server.Dockerfile .

#? docker-client: Build client docker image
.PHONY: docker-client
docker-client:
	docker build -t svalasovich/client -f ./build/package/client.Dockerfile .

#? run: Run Docker Compose
.PHONY: run
run:
	docker compose -f deployments/docker-compose.yaml up

help: Makefile
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sort | sed -e 's/^/ /'
