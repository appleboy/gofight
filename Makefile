.PHONY: all example test

export PROJECT_PATH = /go/src/github.com/appleboy/gofight

all: test

example:
	cd example && go test -v -cover .

test: example
	go test -v -cover .

docker_test: clean
	docker run --rm \
		-v $(PWD):$(PROJECT_PATH) \
		-w=$(PROJECT_PATH) \
		appleboy/golang-testing \
		sh -c "coverage all"

clean:
	rm -rf .cover
