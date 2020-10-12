package tests

import (
	"io/ioutil"
	"log"
	"os"
)

// noop function
func void() {
	return
}

// OpenFile opens a file and returns its descriptor and callback function which
// closes file. callback function has to be used in defer instruction
func OpenFile(filename string) (*os.File, func()) {
	fd, err := os.Open(filename)
	if err != nil {
		log.Printf("there's an issue with open file: %v, `%v`\n", filename, err)
		return nil, void
	}

	return fd, func() {
		_ = fd.Close()
	}
}

func MustReadFile(filename string) (content []byte) {
	var err error
	if content, err = ioutil.ReadFile(filename); err == nil {
		return content
	}
	panic(err)
}
