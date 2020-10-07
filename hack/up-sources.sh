#!/bin/sh
#: very straight forward way, don't be disappointed to much :)
ls cmd/*.go|grep -v test.go | grep -v *_test.go
