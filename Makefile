.PHONY: test build clean

name = fhir-server
build_dir = bin
os = $(shell go env GOOS)
arch = $(shell go env GOARCH)
port = 9090
image_version = v1.0.0
image_tag = $(name)-image:$(image_version)
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

dev-server:
	if [ ! -f $(air_hot_reload) ]; then \
		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s; \
	fi

	$(air_hot_reload) -c .air.toml

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