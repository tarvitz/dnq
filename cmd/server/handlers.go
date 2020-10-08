package server

import (
	"encoding/json"
	"fmt"
	"github.com/tarvitz/dnq/pkg/telegram"
	"net/http"
)

func Default() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(writer, "ok")
	}
}

func Echo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "ok")
	}
}

func Inline() http.HandlerFunc {
	return inline
}

func inline(writer http.ResponseWriter, request *http.Request) {
	var (
		update *telegram.Update
		err    error
	)
	writer.WriteHeader(http.StatusOK)

	update, err = telegram.ReadUpdate(request)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	client := telegram.NewClient(cmd.Token)
	err = client.Inlines.AnswerInlineQuery(update)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func Reload() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := telegram.ReadQuotes(string(cmd.Config))
		if err != nil {
			writer.WriteHeader(400)
			content, _ := json.Marshal(map[string]string{
				"ok":  "false",
				"err": err.Error(),
			})
			_, _ = writer.Write(content)
		}
		_, _ = fmt.Fprintf(writer, `{"status": "ok"}`)
	}
}
