package server

import (
	"bytes"
	"github.com/tarvitz/dnq/pkg/config"
	"github.com/tarvitz/dnq/pkg/telegram"
	"net/http"
	"testing"

	"github.com/tarvitz/dnq/pkg/tests"
)

const (
	testInlineRequestFile = "../../pkg/telegram/resources/inline-query.json"
	testInlineVoiceResult = "../../pkg/telegram/resources/voice.json"
)

func TestDefault(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost", nil)
	buffer := tests.NewHTTPBuffer()
	Default()(buffer, request)
	expected := "ok"

	if expected != buffer.String() {
		t.Errorf("expected: `%v`, got: `%v`", expected, buffer.String())
	}
}

func TestReload(t *testing.T) {
	// reload requires configure cmd
	originCmd := cmd
	defer func() { cmd = originCmd }()
	adminToken := "this-is-test"
	cmd = &Command{
		Config: testConfigFile,
		config: &config.Config{AdminToken: adminToken},
	}

	t.Run("ok", func(in *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost", nil)
		// secret is taken from the config file.
		request.Header.Set(adminTokenHeader, "thisIsASecret")

		buffer := tests.NewHTTPBuffer()
		Reload()(buffer, request)
		expected := `{"status": "ok"}`

		if expected != buffer.String() {
			in.Errorf("expected: `%v`, got: `%v`", expected, buffer.String())
		}
	})

	t.Run("ok/forbidden", func(in *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost", nil)

		buffer := tests.NewHTTPBuffer()
		Reload()(buffer, request)
		expected := `{"status": "forbidden"}`

		if expected != buffer.String() {
			in.Errorf("expected: `%v`, got: `%v`", expected, buffer.String())
		}
	})

	t.Run("failure/wrong-failure", func(in *testing.T) {
		cmd = &Command{Config: "non existent file"}
		request, _ := http.NewRequest("GET", "http://localhost", nil)
		buffer := tests.NewHTTPBuffer()
		Reload()(buffer, request)

		expected := `{"status": "oops, can't reload config file"}`

		if expected != buffer.String() {
			in.Errorf("expected: `%v`, got: `%v`", expected, buffer.String())
		}
	})
}

func TestInline(t *testing.T) {
	client := telegram.NewClient("")
	originCmd := cmd
	defer func() { cmd = originCmd }()

	cmd = &Command{}
	payload := bytes.NewBuffer(tests.MustReadFile(testInlineRequestFile))
	request, _ := http.NewRequest("POST", "http://localhost", payload)

	t.Run("ok", func(in *testing.T) {
		server, rollback := tests.NewHTTPTestServer(tests.Route{
			//: all queries to respond with ok
			"/": &tests.Response{Contents: tests.MustReadFile(testInlineVoiceResult)},
		})
		defer rollback()
		cmd.SetClient(client.WithURL(server.URL))

		buffer := tests.NewHTTPBuffer()
		Mast()(buffer, request)

		expected := `{"status": "ok"}`
		if expected != buffer.String() {
			in.Errorf("expected: `%v`, got: `%v`", expected, buffer.String())
		}
	})

	t.Run("failure/cant-read-update", func(in *testing.T) {
		server, rollback := tests.NewHTTPTestServer(tests.Route{
			"/": &tests.Response{Contents: []byte(`wrong json file`)},
		})
		defer rollback()
		cmd.SetClient(client.WithURL(server.URL))

		buffer := tests.NewHTTPBuffer()
		Mast()(buffer, request)

		expected := ``
		if expected != buffer.String() {
			in.Errorf("expected: `%v`, got: `%v`", expected, buffer.String())
		}

	})
}
