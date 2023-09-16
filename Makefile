.PHONY: test build clean

name = fhird
build_dir = build
os = $(shell go env GOOS)
arch = $(shell go env GOARCH)
port = 9090
image_version = 1.0.0
image_tag = $(name):v$(image_version)
container = $(name)

export DOCKER_BUILDKIT=1

test:
	go test -v ./...

clean:
	rm -rf $(build_dir)/*
	go clean

run-server:
	$(build_dir)/$(name) --verbose

build-server:
	GOOS=$(os) GOARCH=amd64 CGO_ENABLED=0 go build -o $(build_dir)/$(name) .

build-image:
	docker build -t $(image_tag) -f Dockerfile.multistage .

start-container:
	docker run -p $(port):$(port) --name $(container) --rm $(image_tag) 

stop-container:
	docker stop $(container)

start: build-image start-container

stop: stop-container

push-image:
	docker tag $(image_tag) registry.digitalocean.com/fructose/$(image_tag)
	docker push registry.digitalocean.com/fructose/$(image_tag)