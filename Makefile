BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}\
{{end}}' ./...)

VERSION := 0.1.0b
DATE_FMT := +%Y-%m-%d
BUILD_DATE ?= $(shell date "$(DATE_FMT)")

LDFLAGS := -X thunderatz.org/thor/core.Version=$(VERSION) $(LDFLAGS)
LDFLAGS := -X thunderatz.org/thor/core.BuildDate=$(BUILD_DATE) $(LDFLAGS)

EXEC_FILE := thor

all: bin/$(EXEC_FILE)

bin/$(EXEC_FILE): $(BUILD_FILES) | bin
	@echo Version: $(VERSION)
	@echo BuildDate: $(BUILD_DATE)
	@go build -trimpath -ldflags "$(LDFLAGS)" -o $@ ./cmd/thor

build: clean all

bin:
	@mkdir -p bin

clean:
	@-rm -rf bin

.PHONY: build
