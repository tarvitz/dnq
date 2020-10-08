package telegram

type MethodType string

const (
	BotAPIURL = "https://api.telegram.org/bot"

	// Methods
	SendVoice         MethodType = "sendVoice"
	AnswerInlineQuery MethodType = "answerInlineQuery"
)
