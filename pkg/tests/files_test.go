package tests

import (
	"bytes"
	"io/ioutil"
	"testing"
)

const testFile = "../../go.mod"

// just register it in a coverage as far as `void` does nothing.
func Test_void(t *testing.T) {
	void()
}

func TestOpenFile(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		fd, closer := OpenFile(testFile)
		defer closer()
		if fd == nil {
			in.Errorf("expected *os.File, got nil instead")
		}
	})

	t.Run("ok/cant-open-file", func(in *testing.T) {
		fd, _ := OpenFile("non existent")
		if fd != nil {
			in.Errorf("expected nil got: %v", fd)
		}
	})
}

func TestMustReadFile(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		defer func() {
			r := recover()
			if r != nil {
				in.Errorf("paniced")
			}
		}()
		result := MustReadFile(testFile)
		expected, err := ioutil.ReadFile(testFile)
		if !(err == nil && bytes.Equal(result, expected)) {
			in.Errorf("err: `%v`\nexp: `%s`\ngot: `%s`", err, expected, result)
		}
	})

	t.Run("panic", func(in *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				in.Errorf("didn't panic")
			}
		}()
		MustReadFile("non existent file")

	})
}
