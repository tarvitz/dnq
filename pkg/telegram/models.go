package telegram

import (
	"strings"
)

type APIResponse struct {
	Status  bool     `json:"status"`
	Message *Message `json:"result"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot,omitempty"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username"`
	Type         string `json:"type,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Message struct {
	ID   int   `json:"message_id"`
	From *From `json:"from"`
	Chat *From `json:"chat"`
	// todo: replace with serializable date (some other day :)
	Date int    `json:"date"`
	Text string `json:"text,omitempty"`

	// Optional fields
	Voice *Voice `json:"voice,omitempty"`
}

type UpdateType int

const (
	UpdateTypeUnknown UpdateType = iota
	UpdateTypeInline
	UpdateTypeMessage
)

type Update struct {
	ID          int          `json:"update_id"`
	Message     *Message     `json:"message,omitempty"`
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
}

func (update *Update) Type() UpdateType {
	if update.Message != nil {
		return UpdateTypeMessage
	}
	if update.InlineQuery != nil {
		return UpdateTypeInline
	}
	return UpdateTypeUnknown
}

type InlineQuery struct {
	ID     string `json:"id"`
	From   *From  `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

type Voice struct {
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
}

type AnswerInline struct {
	ID      string                `json:"inline_query_id"`
	Results []*AnswerInlineResult `json:"results"`
}

type AnswerInlineType string

const (
	AnswerInlineTypeVoice AnswerInlineType = "voice"
)

type AnswerInlineResult struct {
	Type        AnswerInlineType `json:"type"`
	ID          string           `json:"id"`
	VoiceFileId *string          `json:"voice_file_id,omitempty"`
	Title       string           `json:"title"`
	Caption     string           `json:"caption"`
}

func NewAnswerInline(update *Update) (result *AnswerInline) {
	result = new(AnswerInline)
	result.ID = update.InlineQuery.ID
	for _, quote := range NewInlineQueryResultCachedVoice(update) {
		result.Results = append(result.Results, quote)
	}
	return result
}

func NewInlineQueryResultCachedVoice(update *Update) (results []*AnswerInlineResult) {
	in := strings.ToLower(update.InlineQuery.Query)
	quotes, found := Quotes[in]
	if !found {
		//: set default
		quotes = Quotes[""]
	}
	for _, quote := range quotes {
		result := &AnswerInlineResult{
			ID:          uuid(),
			Type:        AnswerInlineTypeVoice,
			Caption:     quote.Caption,
			Title:       quote.Caption,
			VoiceFileId: &quote.ID,
		}
		results = append(results, result)
	}
	return
}
