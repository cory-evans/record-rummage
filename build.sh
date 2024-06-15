#!/bin/bash

docker buildx build \
	--platform linux/amd64,linux/arm64 \
	-t ghcr.io/cory-evans/record-rummage:latest \
	--push \
	.