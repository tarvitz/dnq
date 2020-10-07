package telegram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// setups http server for api requests mocking.
func setup() (*http.ServeMux, *httptest.Server) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()
	// server is a test HTTP server used to provide mock DomainAPI responses.
	server := httptest.NewServer(mux)
	return mux, server
}

// teardown closes the test HTTP server.
func teardown(server *httptest.Server) {
	server.Close()
}

// Tests a method send by a client.
func testClientCallMethod(t *testing.T, r *http.Request, want string) {
	if r.Method != want {
		t.Errorf("Request method: `%s`, want `%s`", r.Method, want)
	}
}

// returns endpoint and remote urls based on test server URL
func url(endpoint string, srv *httptest.Server) (string, string) {
	return endpoint, fmt.Sprintf("%s%s", srv.URL, endpoint)
}

// mustOpen opens a file or panics.
func mustOpen(filename string) *os.File {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return fd
}

// mustReadFile reads file content
func mustReadFile(filename string) []byte {
	fd := mustOpen(filename)
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	return content
}
