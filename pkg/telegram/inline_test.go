package telegram

import (
	"fmt"
	"net/http"
	"testing"
)

func TestInlineService_AnswerInlineQuery(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	mux.HandleFunc("/"+string(AnswerInlineQuery), func(writer http.ResponseWriter, request *http.Request) {
		testClientCallMethod(t, request, "POST")
		_, _ = fmt.Fprint(writer, `{"status": "ok"}`)
	})
	client := NewClient("").WithURL(server.URL)

	t.Run("ok", func(in *testing.T) {
		update := &Update{
			ID:          1337,
			InlineQuery: &InlineQuery{ID: "133733", Query: "good"},
		}
		err := client.Inlines.AnswerInlineQuery(update)
		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}
