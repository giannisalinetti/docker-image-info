# Makefile variables.
# Update VERSION on new releases
NAME=docker-image-info
VERSION=v0.3.3

.PHONY: all build tag_latest

all: build

build:
	docker build -t $(NAME):$(VERSION) .

tag_latest:
	docker tag $(NAME):$(VERSION) $(NAME):latest
