.PHONY: all example test

all: test

example:
	cd example && go test -v -cover .

test: example
	go test -v -cover .
