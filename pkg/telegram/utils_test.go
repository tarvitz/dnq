package telegram

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func Test_orInt(t *testing.T) {
	for idx, entry := range []struct {
		left     int
		right    int
		expected int
	}{
		{0, 1337, 1337},
		{1, 1337, 1},
		{1337, 0, 1337},
	} {
		t.Run(fmt.Sprintf("test-%d", idx), func(in *testing.T) {
			if result := orInt(entry.left, entry.right); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func Test_orString(t *testing.T) {
	for idx, entry := range []struct {
		left     string
		right    string
		expected string
	}{
		{"", "test", "test"},
		{" ", "test", " "},
		{"1", "2", "1"},
		{"", "", ""},
	} {
		t.Run(fmt.Sprintf("test-%d", idx), func(in *testing.T) {
			if result := orString(entry.left, entry.right); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestClose(t *testing.T) {
	body := bytes.NewBuffer([]byte{}) // nil body will panic
	request, _ := http.NewRequest("GET", "http://localhost", body)
	Close(request.Body)
}

// just a register a call in coverage, as far as there's nothing much to test.
func Test_noop(t *testing.T) {
	noop()
}

func Test_safeClose(t *testing.T) {
	t.Run("ok/only-close-having", func(in *testing.T) {
		request, _ := http.NewRequest(
			"GET", "http://localhost", bytes.NewBuffer([]byte{}))
		call := safeClose([]interface{}{request.Body})
		call()
	})

	t.Run("ok/mixed", func(in *testing.T) {
		request, _ := http.NewRequest(
			"GET", "http://localhost", bytes.NewBuffer([]byte{}))
		//: only request.Body is Close compatible, other items will be omitted.
		call := safeClose([]interface{}{request.Body, request, request.Host})
		call()
	})

	t.Run("ok/no-values", func(in *testing.T) {
		call := safeClose([]interface{}{})
		call()
	})
}

func Test_uuid(t *testing.T) {
	once := uuid()
	twice := uuid()

	if once == twice {
		t.Errorf("uuid() is expected to return unique values, got: `%v` = `%v`", once, twice)
	}
}
