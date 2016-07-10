.PHONY: all example test

TARGETS_NOVENDOR := $(shell glide novendor)

all: install test

install:
	glide install

update:
	glide update --all-dependencies --resolve-current

example:
	cd example && go test -v -cover -covermode=count -coverprofile=coverage.out .

example_html: example
	cd example && go tool cover -html=coverage.out && unlink coverage.out

test: example
	go test -v -cover -covermode=count -coverprofile=coverage.out .

test_html: test
	go tool cover -html=coverage.out && unlink coverage.out

lint:
	golint example
	golint .
