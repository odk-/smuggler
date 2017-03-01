IMAGE_NAME=odkurzacz/smuggler

.DEFAULT_GOAL := default

default: build build-docker

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smuggler .
	strip smuggler
build-docker:
	docker build -t $(IMAGE_NAME) .
start:
	docker run -d -P --name="smuggler" $(IMAGE_NAME)
stop:
	docker rm -f smuggler
