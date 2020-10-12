package server

import (
	"fmt"
	"net/http"

	"github.com/tarvitz/dnq/pkg/config"
	"github.com/tarvitz/dnq/pkg/telegram"
)

func Default() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(writer, "ok")
	}
}

func Inline() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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

		client := cmd.GetClient()
		err = client.Inlines.AnswerInlineQuery(update)
		// well, telegram does not require a payload in answer,
		// thus far it just returns a simple json object
		if err == nil {
			_, _ = fmt.Fprintf(writer, `{"status": "ok"}`)
		}
	}
}

func Reload() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error
		if cmd.config, err = config.ReadConfig(string(cmd.Config)); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(writer, `{"status": "oops, can't reload config file"}`)
			return
		}

		if hasAdminToken(request, cmd.config.AdminToken) {
			fmt.Printf("Quotes reloaded: %d\n", len(cmd.config.Quotes))
			telegram.SetQuotes(cmd.config.Quotes)
			_, _ = fmt.Fprintf(writer, `{"status": "ok"}`)
		} else {
			writer.WriteHeader(http.StatusForbidden)
			_, _ = fmt.Fprintf(writer, `{"status": "forbidden"}`)
		}
	}
}
