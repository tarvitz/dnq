package telegram

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TestObject struct {
	ID int  `json:"id"`
	Ok bool `json:"ok"`
}

func TestNewAPI(t *testing.T) {
	client := NewClient("1337token")
	if reflect.TypeOf(client) != reflect.TypeOf(&Client{}) {
		t.Error("types does not match")
	}
}

func TestNewClient(t *testing.T) {
	expected := "1337token"
	client := NewClient(expected)
	if client.token != expected {
		t.Errorf("\nexp: %v\ngot: %v", expected, client.token)
	}
}

func TestClient_URL(t *testing.T) {
	client := Client{
		apiURL: "http://api-example/v4",
	}
	url := client.URL("/getUpdates")
	expected := "http://api-example/v4/getUpdates"
	if url != expected {
		t.Errorf("expected: %v, got: %v", expected, url)
	}
}

func TestClient_WithURL(t *testing.T) {
	expected := "https://localhost/api/"
	client := NewClient("").WithURL(expected)
	if client.apiURL != expected {
		t.Errorf("expected: %v, got: %v", expected, client.apiURL)
	}
}

func TestClient_Call(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	mux.HandleFunc("/getUpdates", func(writer http.ResponseWriter, request *http.Request) {
		testClientCallMethod(t, request, "GET")
		_, _ = fmt.Fprint(writer, `[{"id":1},{"id":2}]`)
	})
	client := NewClient("").WithURL(server.URL)

	t.Run("failure/404", func(t *testing.T) {
		err := client.Call(&APIRequest{Method: "404"}, nil)
		if err == nil {
			t.Errorf("expected error got nil instead")
		}
	})

	t.Run("failure/wrong-domain", func(in *testing.T) {
		client.httpClient = &http.Client{Timeout: 10 * time.Millisecond}
		client.apiURL = "http://api.fake-domain"
		err := client.Call(&APIRequest{Method: "1"}, nil)

		if err == nil {
			t.Errorf("expected error got nil instead")
		}
	})
}

func TestClient_Retrieve(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	mux.HandleFunc("/getUpdates", func(writer http.ResponseWriter, request *http.Request) {
		testClientCallMethod(t, request, "GET")
		_, _ = fmt.Fprint(writer, `[{"id":1},{"id":2}]`)
	})

	client := NewClient("").WithURL(server.URL)

	var testObject []*TestObject
	expected := []*TestObject{{ID: 1}, {ID: 2}}

	t.Run("ok", func(in *testing.T) {
		err := client.Retrieve("getUpdates", &testObject)

		if err != nil {
			in.Errorf("got error: %v", err)
		}

		if !reflect.DeepEqual(testObject, expected) {
			in.Errorf("expected: `%+v`, got: `%+v`", expected, testObject)
		}
	})
}

func TestClient_Upload(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	mux.HandleFunc("/sendVoice", func(writer http.ResponseWriter, request *http.Request) {
		testClientCallMethod(t, request, "POST")
		voiceResponse := mustReadFile(testVoiceResponseFile)
		_, _ = fmt.Fprint(writer, string(voiceResponse))
	})

	t.Run("ok", func(in *testing.T) {
		client := NewClient("").WithURL(server.URL)
		message, err := client.Upload(SendVoice, map[string]io.Reader{
			"voice":    strings.NewReader("fake voice content"),
			"chat_id":  strings.NewReader("1337"),
			"caption":  strings.NewReader("I am Duke!"),
			"duration": strings.NewReader("2"),
		})

		if err != nil {
			in.Errorf("got error: %v", err)
		}

		expected := &Voice{
			Duration:     2,
			MimeType:     "audio/ogg",
			FileID:       "AwACAgIAAxkDAAMVX3m--xYKLFLzyvr54__SS5DDJGAAAi4IAAJgG9BLwR3pDPW5EM4bBA",
			FileUniqueID: "AgADLggAAmAb0Es",
			FileSize:     27275,
		}
		if !reflect.DeepEqual(message.Voice, expected) {
			in.Errorf("\nexp: %v\ngot: `%v`", expected, message.Voice)
		}
	})
}

func TestReadUpdate(t *testing.T) {
	content := mustReadFile(testInlineQueryFile)
	body := bytes.NewBuffer(content)
	request, _ := http.NewRequest("GET", "http://localhost/", body)
	update, err := ReadUpdate(request)

	expected := &Update{
		ID: 292124505,
		InlineQuery: &InlineQuery{
			ID: "600597106931592670",
			From: &From{
				ID:           1337,
				IsBot:        false,
				FirstName:    "Nickolas",
				LastName:     "F.",
				Username:     "nickolasfox",
				LanguageCode: "ru",
			},
			Query:  "groovy",
			Offset: "",
		},
	}

	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if !reflect.DeepEqual(update, expected) {
		t.Errorf("\nexp: `%v`\ngot: %v", expected, update)
	}
}

func Test_makeMultipartFormDataRequest(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		fd := mustOpen(testVoiceFile)
		_, contentType := makeMultipartFormDataPayload(map[string]io.Reader{
			"voice":    fd,
			"chat_id":  strings.NewReader("1337"),
			"caption":  strings.NewReader("I am Duke!"),
			"duration": strings.NewReader("2"),
		})

		expected := "multipart/form-data; boundary="
		if strings.Contains(expected, contentType) {
			in.Errorf("expected: `%v` hasn't been found in: `%v`", expected, contentType)
		}
	})
}
