package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tarvitz/dnq/pkg/config"
	"github.com/tarvitz/dnq/pkg/telegram"
)

const (
	//: help voice message (already preloaded)
	helpVoiceID = "AwACAgIAAxkDAAM1X4TWF69s2eFdU3INm1VZdlGQHcwAAv0KAAJl8ShIiwo4VB6W5hobBA"
	statusOk    = `{"status": "ok"}`
)

// Default is a default http response handler responding with HTTP 200 OK
func Default() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(writer, "ok")
	}
}

// Mast is like a telegraph mast that receives messages from telegram
// and responds back.
func Mast() http.HandlerFunc {
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

		switch update.Type() {
		case telegram.UpdateTypeInline:
			inline(writer, update)
		case telegram.UpdateTypeMessage:
			message(writer, update)
		default:
			fmt.Printf("unknown type of message, skipping \n")
		}
	}
}

func inline(writer http.ResponseWriter, update *telegram.Update) {
	var err error
	client := cmd.GetClient()

	if err = client.Inlines.AnswerInlineQuery(update); err == nil {
		// well, telegram does not require a payload in answer,
		// thus far it just returns a simple json object
		_, _ = fmt.Fprintf(writer, statusOk)
	}
}

func message(writer http.ResponseWriter, update *telegram.Update) {
	var err error
	client := cmd.GetClient()
	text := strings.ToLower(update.Message.Text)

	quote := cmd.config.RandomQuote()
	voice := &telegram.VoiceMessage{
		ChatID:  update.Message.From.ID,
		Voice:   quote.ID,
		Caption: &quote.Caption,
	}
	if strings.HasPrefix(text, "/help") {

		voice.Voice = helpVoiceID
		caption := "Help me! Help me, ehehehe"
		voice.Caption = &caption
	}
	err = client.Commons.SendVoice(voice)

	if err == nil {
		_, _ = fmt.Fprintf(writer, statusOk)
	} else {
		fmt.Printf("err: %v\n", err)
	}
}

// Reload reloads quotes configuration, this is useful to reload quotes without
// restarting a server
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
