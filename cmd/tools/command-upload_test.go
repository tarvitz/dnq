package tools

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/tarvitz/dnq/cmd/common"
	"github.com/tarvitz/dnq/pkg/telegram"
	"github.com/tarvitz/dnq/pkg/tests"
)

const (
	testUploadResponse = "resources/message.json"
	testOggOpusFile    = "../../pkg/ogg/resources/opus-headers-only.ogg"
	testNotOggFile     = "../../go.mod"
)

var (
	testUploadCommand = UploadCommand{
		Auth: common.Auth{},
		UploadPositional: UploadPositional{
			Filename: testOggOpusFile,
		},
		ChatID:  "1337",
		Caption: "this is a test caption.",
		Matches: []string{"test", "me"},
		Output:  "config",
	}

	testMessage = &telegram.Message{
		ID: 1337,
		From: &telegram.From{
			ID:        1337,
			IsBot:     true,
			FirstName: "Duke Nukem Quotes",
			Username:  "duqe_bot",
		},
		Chat: &telegram.From{
			ID:        133733,
			IsBot:     false,
			FirstName: "Nickolas",
			LastName:  "F.",
			Username:  "nickolasfox",
			Type:      "private",
		},
		Date: 1602355724,
		Voice: &telegram.Voice{
			Duration:     2,
			MimeType:     "audio/ogg",
			FileID:       "AwACAgIAAxkDAAMzX4ICDLk-BdYO1ce_sz0Fdy9o-ngAAngKAAKmRBFI0AdUqZbbDLsbBA",
			FileUniqueID: "AgADeAoAAqZEEUg",
			FileSize:     24169,
		},
	}
)

func TestUploadCommand_Execute(t *testing.T) {
	client := telegram.NewClient("")
	cmd := testUploadCommand

	t.Run("ok", func(in *testing.T) {
		contents, _ := ioutil.ReadFile(testUploadResponse)
		server, rollback := tests.NewHTTPTestServer(tests.Route{
			//: i.e. every traffic
			"/": &tests.Response{Contents: contents},
		})
		defer rollback()

		cmd.SetClient(client.WithURL(server.URL))
		err := cmd.Execute([]string{})
		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})

	t.Run("failure/not-authenticated", func(in *testing.T) {
		server, rollback := tests.NewHTTPTestServer(tests.Route{
			"/": &tests.Response{
				Options:  &tests.Options{Status: http.StatusUnauthorized},
				Contents: []byte(`{"status: false}`),
			},
		})
		defer rollback()
		cmd.SetClient(client.WithURL(server.URL))

		err := cmd.Execute([]string{})
		if err == nil {
			in.Errorf("expected error but got unexpected nil")
		}
	})

	t.Run("failure/cant-open-ogg-file", func(in *testing.T) {
		cmd := testUploadCommand
		cmd.Filename = "not existent file"

		err := cmd.Execute([]string{})
		if err == nil {
			in.Errorf("expected error but got unexpected nil")
		}
	})

	t.Run("failure/not-an-ogg-file", func(in *testing.T) {
		cmd := testUploadCommand
		cmd.Filename = testNotOggFile

		err := cmd.Execute([]string{})
		if err == nil {
			in.Errorf("expected error but got unexpected nil")
		}
	})
}

//: just register calls as far as no much to check/verify.
func TestUploadCommand_print(t *testing.T) {
	cmd := testUploadCommand
	for _, entry := range []struct {
		name string
	}{
		{"config"}, {"json"}, {"yaml"}, {"unknown"},
	} {
		t.Run(entry.name, func(in *testing.T) {
			cmd.Output = entry.name
			cmd.print(testMessage)
		})
	}
}
