package telegram

// MethodType +
type MethodType string

// Exportable constants
const (
	BotAPIURL = "https://api.telegram.org/bot"

	// Methods
	SendVoice         MethodType = "sendVoice"
	AnswerInlineQuery MethodType = "answerInlineQuery"
)
