package tests

import (
	"log"
	"os"
)

func void() {}

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
