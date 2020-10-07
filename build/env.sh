#!/bin/bash -x
APP_BIN="${APP_BIN:-dnq}"
TS=$(date -u +%Y%m%d.%H%M%S)
VERSION=${APP_VERSION:-"dev-${TS}"}

if [[ ! -d bin ]]; then
    mkdir bin
fi
