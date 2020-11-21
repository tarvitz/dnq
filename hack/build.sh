#!/bin/bash
DEFAULT_APP_VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))-dev
APP_VERSION=${APP_VERSION:-$DEFAULT_APP_VERSION}
docker build --build-arg=APP_VERSION=${APP_VERSION} -f Dockerfile . \
    -t nfox/dnq:${APP_VERSION}
