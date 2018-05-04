BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
GOVENDOR := $(BIN_DIR)/govendor
PACKR := $(BIN_DIR)/packr

RELEASE_DIR := "release/"
BINARY := "compare"

VERSION := $(shell cat VERSION)
GIT_REV := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.version=$(VERSION) -X main.hash=$(GIT_REV) -X main.date=`date -u +%Y.%m.%d`"

PLATFORMS := linux darwin windows

os = $(word 1, $@)

.PHONY: test lint clean update assets release $(PLATFORMS)

release: update assets $(PLATFORMS)

update: $(GOVENDOR)
	@echo "Updating dependancies"
	@govendor sync

assets:
	@echo "Packing ressources"
	@packr

windows: EXT = ".exe"
linux darwin: EXT = ""
$(PLATFORMS):
	@echo "Building $(os)"
	@mkdir -p $(RELEASE_DIR)
	@GOOS=$(os) GOARCH=amd64 go build -o $(RELEASE_DIR)$(BINARY)-$(VERSION)-$(os)$(EXT) -ldflags $(LDFLAGS)

clean:
	@echo "Cleaning workspace"
	@packr clean
	@rm -rf $(RELEASE_DIR)

test: $(GOVENDOR)
	@echo "Running tests"
	@govendor test -cover ./... +local

lint: $(GOMETALINTER)
	@echo "Running lint"
	@gometalinter --vendor ./...

$(GOMETALINTER):
	@echo "Installing gometalinter"
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install &> /dev/null

$(GOVENDOR):
	@echo "Installing govendor"
	@go get -u github.com/kardianos/govendor
