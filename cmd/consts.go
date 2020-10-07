package main

import (
	"errors"
)

const (
	appName = "dnq"
	//: env variables

	//: exit codes
	exitOk                = 0
	exitCodeParserFailure = 1
	exitCommandRequired   = 2
	exitCommandError      = 3
)

var (
	//: errors
	errWrongCompletionType = errors.New("unsupported completion type")
)
