#!/bin/bash

set -e

DOCKER_IMAGE_NAME="docker-machine-driver-sakuracloud-build"
DOCKER_CONTAINER_NAME="docker-machine-driver-sakuracloud-build-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker build -t $DOCKER_IMAGE_NAME .

docker run --name $DOCKER_CONTAINER_NAME \
  -e DEBUG \
  -e STATIC \
  -e VERBOSE \
  -e BUILDTAGS \
  -e PARALLEL \
  -e COVERAGE_DIR \
  -e TARGET_OS \
  -e TARGET_ARCH \
  -e PREFIX \
  -e TRAVIS_JOB_ID \
  -e TRAVIS_PULL_REQUEST \
  -e SAKURACLOUD_ACCESS_TOKEN \
  -e SAKURACLOUD_ACCESS_TOKEN_SECRET \
  $DOCKER_IMAGE_NAME \
  make "$@"

if [[ -d bin ]]; then
  rm -Rf bin
fi

if [[ "$@" == *"build"* ]]; then
  docker cp $DOCKER_CONTAINER_NAME:/go/src/github.com/yamamoto-febc/docker-machine-sakuracloud/bin bin
fi
