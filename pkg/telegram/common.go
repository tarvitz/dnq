package telegram

import (
	"bytes"
	"encoding/json"
)

// CommonService is served for common telegram methods
type CommonService struct {
	client *Client
}

// VoiceMessage processes telegram voice messages
type VoiceMessage struct {
	ChatID              int     `json:"chat_id"`
	Voice               string  `json:"voice"`
	Caption             *string `json:"caption,omitempty"`
	ParseMode           *string `json:"parse_mode,omitempty"`
	Duration            *int    `json:"duration,omitempty"`
	DisableNotification *int    `json:"disable_notification,omitempty"`
	ReplyToMessageID    *int    `json:"reply_to_message_id,omitempty"`
}

// SendVoice sends a voice message to the client
func (service *CommonService) SendVoice(message *VoiceMessage) (err error) {
	content, _ := json.Marshal(message)
	buffer := bytes.NewBuffer(content)

	request := &APIRequest{
		Method:    "POST",
		APIMethod: SendVoice,
		Body:      buffer,
	}
	err = service.client.Call(request, nil)
	return
}
