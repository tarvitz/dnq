package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewHTTPTestServer(t *testing.T) {
	server, closeFunc := NewHTTPTestServer(Route{
		`^/$`:        &Response{Contents: []byte(`this is test`)},
		`^/options$`: &Response{Options: &Options{Status: http.StatusForbidden}},
	})
	defer closeFunc()

	t.Run("ok", func(in *testing.T) {
		var content []byte

		response, err := http.Get(server.URL)
		defer func() { _ = response.Body.Close() }()

		if err != nil {
			in.Errorf("got error: %v", err)
		}
		if content, err = ioutil.ReadAll(response.Body); err != nil {
			in.Errorf("got error: %v", err)
		}
		expected := []byte(`this is test`)
		if !bytes.Equal(content, expected) {
			in.Errorf("expected: `%s`, got: `%s`", expected, content)
		}
	})

	t.Run("ok/not-found", func(in *testing.T) {
		url := fmt.Sprintf("%s/%s", server.URL, "this/is/the-test")
		response, _ := http.Get(url)
		defer func() { _ = response.Body.Close() }()

		if response.StatusCode != http.StatusNotFound {
			in.Errorf("expected: %v, got: %v", http.StatusNotFound, response.StatusCode)
		}
	})

	t.Run("ok/options", func(in *testing.T) {
		url := fmt.Sprintf("%s/%s", server.URL, "options")
		response, _ := http.Get(url)
		defer func() { _ = response.Body.Close() }()

		if response.StatusCode != http.StatusForbidden {
			in.Errorf("expected: %v, got: %v", http.StatusForbidden, response.StatusCode)
		}
	})

	// wrong regex will be skipped, default status is 404
	t.Run("ok/wrong-regex", func(in *testing.T) {
		route := Route{
			//: this should not work due to wrong regex, fallback to default
			`^/test/.?+*$`: &Response{Options: &Options{Status: http.StatusCreated}},
		}
		server, rollback := NewHTTPTestServer(route)
		defer rollback()

		url := fmt.Sprintf("%s/%s", server.URL, "test/check")
		response, _ := http.Get(url)
		defer func() { _ = response.Body.Close() }()

		if response.StatusCode != http.StatusNotFound {
			in.Errorf("expected: %v, got: %v", http.StatusNotFound, response.StatusCode)
		}
	})
}

func TestNewHTTPBuffer(t *testing.T) {
	buffer := NewHTTPBuffer()
	buffer.WriteHeader(http.StatusOK)
	//: just to register in a coverage
	buffer.Header()

	expected := "this is a test"
	_, _ = fmt.Fprintf(buffer, expected)
	if result := buffer.String(); result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}
