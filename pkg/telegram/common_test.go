package telegram

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCommonService_SendVoice(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	mux.HandleFunc("/"+string(SendVoice), func(writer http.ResponseWriter, request *http.Request) {
		testClientCallMethod(t, request, "POST")
		_, _ = fmt.Fprint(writer, `{"status": "ok"}`)
	})
	client := NewClient("").WithURL(server.URL)

	t.Run("ok", func(in *testing.T) {
		update := &VoiceMessage{
			ChatID: 1337,
			Voice:  "voice-id",
		}
		err := client.Commons.SendVoice(update)
		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}
