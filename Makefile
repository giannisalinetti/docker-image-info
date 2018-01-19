# Makefile variables.
# Update VERSION on new releases
NAME=docker_image_info
VERSION=v0.2

.PHONY: all build tag_latest

all: build

build:
	docker build -t $(NAME):$(VERSION) .

tag_latest:
	docker tag $(NAME):$(VERSION) $(NAME):latest
