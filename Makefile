# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
BINARY_NAME := nsproxy
BINARY_PATH := bin/$(BINARY_NAME)
CONFIG_FILE := config.json
SYSTEMD_FILE := init/nsproxy.service

all: test build

build:
	$(GOBUILD) -o $(BINARY_PATH) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

run:
	$(GOBUILD) -o $(BINARY_PATH) -v
	./$(BINARY_PATH)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_PATH) -v

install: build install-binary install-config install-service

install-binary:
	cp $(BINARY_PATH) /usr/bin/$(BINARY_NAME)

install-config:
	cp $(CONFIG_FILE) /etc/nsproxy.json

install-service:
	cp $(SYSTEMD_FILE) /etc/systemd/system/
	systemctl enable nsproxy

uninstall: uninstall-service uninstall-binary uninstall-config

uninstall-binary:
	rm -f /usr/bin/$(BINARY_NAME)

uninstall-config:
	rm -f /etc/nsproxy.json

uninstall-service:
	systemctl stop nsproxy
	systemctl disable nsproxy
	rm -f /etc/systemd/system/nsproxy.service

.PHONY: all build test clean run build-linux install uninstall install-binary install-config install-service uninstall-binary uninstall-config uninstall-service
