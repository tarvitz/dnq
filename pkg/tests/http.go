package tests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
)

// Route is a map of regular expression-like strings and response objects
// basically keys of route matches over regular expression with request url path and if
// match occurs, response will be used. A tiny plain and simple routing based on
// urls
type Route map[string]*Response

// Options sets different response options to set in http testing responses
type Options struct {
	Status int
}

// Response is used for a testing response
type Response struct {
	Options  *Options
	Contents []byte
}

var notFound http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintf(writer, "not found")
}

func (route Route) match(path string) (string, bool) {
	//: a little bit `unoptimal` due each time we traverse over whole route (map)
	//: object to find matching route key. However, for test cases it's more than ok
	for key := range route {
		//: skip invalid regex
		regex, err := regexp.Compile(key)
		if err != nil {
			log.Printf("could not compile regex: `%s` please fix it: `%v`", key, err)
			continue
		}
		if regex.Match([]byte(path)) {
			return key, true
		}
	}
	return "", false
}

// NewHTTPTestServer creates a server with default 404 content response return
// logic. One can configure server with route objects and responses that should
// be returned on the certain urls (like ^/artifactory/.*$ will be matched with
// all urls that starts with /artifactory/)
func NewHTTPTestServer(route Route) (*httptest.Server, func()) {
	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		key, ok := route.match(request.URL.Path)
		if !ok {
			notFound(writer, request)
			return
		}

		response := route[key]
		if response.Options != nil {
			writer.WriteHeader(response.Options.Status)
		}
		_, _ = writer.Write(response.Contents)
	}

	server := httptest.NewServer(handler)
	return server, func() {
		server.Close()
	}
}

// HTTPBuffer is a http in-memory buffer used for testing
type HTTPBuffer struct {
	buffer *bytes.Buffer
	status int
}

// Write implements http.ResponseWriter interface.
func (h *HTTPBuffer) Write(p []byte) (n int, err error) {
	return h.buffer.Write(p)
}

// Header implements http.ResponseWriter interface.
func (h *HTTPBuffer) Header() http.Header {
	return http.Header{}
}

// WriteHeader implements http.ResponseWriter interface.
func (h *HTTPBuffer) WriteHeader(statusCode int) {
	h.status = statusCode
}

// String implements fmt.Stringer interface.
func (h *HTTPBuffer) String() string {
	return h.buffer.String()
}

// NewHTTPBuffer creates a new buffer object.
func NewHTTPBuffer() *HTTPBuffer {
	buffer := &HTTPBuffer{}
	buffer.buffer = bytes.NewBuffer([]byte{})
	return buffer
}
