#!/bin/bash -xe
functions="$(dirname "$0")/env.sh"
if [ -f "$functions" ]; then
  # shellcheck disable=SC1090
  source "$functions"
fi

# shellcheck disable=SC2046
go build -ldflags "-X main.appVersion=${VERSION} -w -s" -o bin/ci $(cat sources.txt)
