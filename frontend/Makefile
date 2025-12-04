IMAGE_NAME := asia-northeast1-docker.pkg.dev/glossy-attic-443507-u6/docker/frontend
TAG := $(shell git rev-parse --short HEAD)

.PHONY: build push

build:
	docker build --platform linux/amd64 -t $(IMAGE_NAME):$(TAG) .
	docker tag $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest

push:
	docker push $(IMAGE_NAME):$(TAG)
	docker push $(IMAGE_NAME):latest
