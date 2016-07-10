.PHONY: all example test

TARGETS_NOVENDOR := $(shell glide novendor)

all: install test

install:
	glide install

update:
	glide update --all-dependencies --resolve-current

example:
	cd example && go test -v -cover -covermode=count -coverprofile=coverage.txt .

example_html: example
	cd example && go tool cover -html=coverage.txt && unlink coverage.txt

test: example
	go test -v -cover -covermode=count -coverprofile=coverage.txt .

test_html: test
	go tool cover -html=coverage.txt && unlink coverage.txt

lint:
	golint example
	golint .
