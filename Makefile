.PHONY: all example test

all: install test

install:
	glide install

update:
	glide update --all-dependencies --resolve-current

example:
	cd example && go test -v -cover -covermode=count -coverprofile=coverage.out .

test: example
	go test -v -cover -covermode=count -coverprofile=coverage.out .

lint:
	golint .
