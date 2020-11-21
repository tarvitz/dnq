#!/bin/bash
ALLOWED_CYCLE_LEVEL=20

gocyclo -over ${ALLOWED_CYCLE_LEVEL} cmd pkg
golint ./cmd/... ./pkg/...
ineffassign ./cmd/* ./pkg/*
