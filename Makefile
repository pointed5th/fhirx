.PHONY: build test clean

NAME=fhird
BUILD_DIR=build
TARGET=$(BUILD_DIR)/$(NAME)
SRC=./cmd/server/main.go
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
VERSION=1.0.0
DOCKER_IMAGE=$(NAME)
DOCKER_IMAGE_TAG=$(NAME):v$(VERSION)
DOCKER_CONTAINER_NAME=$(NAME)-container

build:
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -o $(TARGET) $(SRC)

docker-build:
	docker build -t $(DOCKER_IMAGE_TAG) -f Dockerfile.multistage . --no-cache

test:
	go test -v ./...

clean:
	rm -rf $(build_dir)/*
	go clean

run: build
	$(TARGET)


# .PHONY: test build clean

# NAME=fructose
# BUILD_DIR=bin
# PORT = 9090
# DOCKER_CONTAINER_NAME = $(name)-container
# DOCKER_IMAGE_VERSION = 0.0.1
# DOCKER_IMAGE_TAGE = $(name):v$(image_version)
# AIR_HOTRELOAD = bin/air
# GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(NAME) .

# export DOCKER_BUILDKIT=1

# test:
# 	go test -v ./...

# test-db:
# 	go test -timeout 30s -run ^TestNewDbHandler -v ./...

# clean:
# 	rm -rf $(build_dir)/*
# 	go clean

# dev:
# 	if [ ! -f $(air_hotreload) ]; then \
# 		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s; \
# 	fi

# 	$(air_hotreload) -c .air.toml

# build:
# 	GOOS=$(os) GOARCH=${arch} CGO_ENABLED=0 go build -o $(build_dir)/$(name) .	
	
# build-image:
# 	docker build -t $(image_tag) -f Dockerfile.multistage . --no-cache

# start:
# 	docker run -p $(port):$(port) --name $(container) --rm $(image_tag) 

# stop:
# 	docker stop $(container)

# push:
# 	docker tag $(image_tag) registry.digitalocean.com/fructose/$(image_tag)
# 	docker push registry.digitalocean.com/fructose/$(image_tag)
