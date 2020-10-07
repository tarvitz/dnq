package telegram

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	u "net/url"
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
	client := NewClient("https://localhost/telegram/api/bot1337token")
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

func TestClient_request(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	endpoint, url := url("/getUpdates", server)
	contentExpected := []byte(`ok`)
	mux.HandleFunc(endpoint, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(writer, string(contentExpected))
	})

	client := NewClient("testToken")
	request := client.request("GET", url, nil)
	response, err := client.httpClient.Do(request)

	if err != nil {
		t.Errorf("got error: %v", err)
	}
	defer func() { Close(response.Body) }()

	result, _ := ioutil.ReadAll(response.Body)
	if !bytes.Equal(result, contentExpected) {
		t.Errorf("expected: `%s`, got: `%s`", contentExpected, result)
	}
}

func TestClient_Retrieve(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	mux.HandleFunc("/getUpdates", func(writer http.ResponseWriter, request *http.Request) {
		testClientCallMethod(t, request, "GET")
		_, _ = fmt.Fprint(writer, `[{"id":1},{"id":2}]`)
	})

	client := NewClient("").WithURL(server.URL)

	values := u.Values{"ShowOk": {"True"}}
	reader := bytes.NewBuffer([]byte(values.Encode()))

	var testObject []*TestObject
	expected := []*TestObject{{ID: 1}, {ID: 2}}

	t.Run("ok", func(in *testing.T) {
		err := client.Retrieve("getUpdates", reader, &testObject)

		if err != nil {
			in.Errorf("got error: %v", err)
		}

		if !reflect.DeepEqual(testObject, expected) {
			in.Errorf("expected: `%+v`, got: `%+v`", expected, testObject)
		}
	})

	t.Run("failure/404", func(t *testing.T) {
		err := client.Retrieve("404", nil, nil)
		if err == nil {
			t.Errorf("expected error got nil instead")
		}
	})

	t.Run("failure/wrong-domain", func(in *testing.T) {
		client.httpClient = &http.Client{Timeout: 10 * time.Millisecond}
		client.apiURL = "http://api.fake-domain"
		err := client.Retrieve("1", nil, nil)

		if err == nil {
			t.Errorf("expected error got nil instead")
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
			"voice": strings.NewReader("fake voice content"),
			"chat_id": strings.NewReader("1337"),
			"caption": strings.NewReader("I am Duke!"),
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


// Testing 404s of Retrieve method
func TestClient_404s(t *testing.T) {
	_, server := setup()
	defer teardown(server)
	client := NewClient(server.URL)

	for _, entry := range []struct {
		name   string
		method interface{}
	}{
		{"Retrieve", client.Retrieve},
	} {
		t.Run(entry.name, func(in *testing.T) {
			var err error
			switch entry.name {
			case "Retrieve":
				err = entry.method.(func(string, io.Reader, interface{}) error)("/404", nil, nil)
			}

			if err == nil {
				in.Errorf("expected error, got nil instead.")
			}
		})
	}
}

// Testing timeouts issues of Retrieve  method
func TestClient_Timeouts(t *testing.T) {
	_, server := setup()
	defer teardown(server)
	client := NewClient(server.URL)
	client.httpClient = &http.Client{Timeout: 10 * time.Millisecond}

	for _, entry := range []struct {
		name   string
		method interface{}
	}{
		{"Retrieve", client.Retrieve},
	} {
		t.Run(entry.name, func(in *testing.T) {
			var err error
			switch entry.name {
			case "Retrieve":
				err = entry.method.(func(string, io.Reader, interface{}) error)("1", nil, nil)
			}

			if err == nil {
				t.Errorf("expected error got nil instead")
			}
		})
	}
}

func TestClient_updateRequest(t *testing.T) {
	t.Run("query", func(in *testing.T) {
		client := NewClient("")
		request := client.request("GET", "http://localhost:8000/", nil)
		query := u.Values{}
		query.Set("test", "1337")
		client.updateRequest(&APIRequest{Query: query.Encode()}, request)

		result := request.URL.Query()
		if !reflect.DeepEqual(query, result) {
			in.Errorf("query does not match: `%v` != `%v`", query, result)
		}
	})
}

func Test_makeMultipartFormDataRequest(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		fd := mustOpen(testVoiceFile)
		_, _,  err := makeMultipartFormDataPayload(map[string]io.Reader{
			"voice": fd,
			"chat_id": strings.NewReader("1337"),
			"caption": strings.NewReader("I am Duke!"),
			"duration": strings.NewReader("2"),
		})

		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}
