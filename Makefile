.PHONY: test build clean

name = fhird
build_dir = bin
os = $(shell go env GOOS)
arch = $(shell go env GOARCH)
port = 9090
image_version = 0.0.1
image_tag = $(name)-image:v$(image_version)
container = $(name)-container
air_hot_reload = bin/air

export DOCKER_BUILDKIT=1

test:
	go test -v ./...

clean:
	rm -rf $(build_dir)/*
	go clean

run-server:
	$(build_dir)/$(name) --verbose

build-server:
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -o $(build_dir)/$(name) .

build-image:
	docker build -t $(image_tag) -f Dockerfile.multistage .

start-container:
	docker run -p $(port):$(port) --name $(container) --rm $(image_tag) 

stop-container:
	docker stop $(container)

start: build-image start-container

stop: stop-container

fhir-examples:
	wget https://www.hl7.org/fhir/R4/examples-json.zip -O examples-json.zip
	unzip examples-json.zip -d examples-json
	rm -rf examples-json.zip