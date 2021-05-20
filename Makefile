PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.PHONY: build
build:  ## Build application
	go build -o build/application $(PROJECT_PATH)/cmd/server/main.go