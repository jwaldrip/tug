NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

REVISION=$(shell git rev-parse --short HEAD)
BASE_VERSION=$(shell cat VERSION)
VERSION=$(BASE_VERSION)-$(REVISION)

all: build

build:
	go install .

test:
	@echo "$(OK_COLOR)==> Testing...$(NO_COLOR)"
	@script/test $(TEST)

release:
	script/cross-compile $(VERSION)
	cd build && tar czf tug.tgz tug
	s3cmd put -P build/tug.tgz s3://tug-binaries/tug.tgz

.PHONY: build test upload
