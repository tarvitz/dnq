package telegram

import (
	"bytes"
	"encoding/json"
)

type CommonService struct {
	client *Client
}

type VoiceMessage struct {
	ChatID              int     `json:"chat_id"`
	Voice               string  `json:"voice"`
	Caption             *string `json:"caption,omitempty"`
	ParseMode           *string `json:"parse_mode,omitempty"`
	Duration            *int    `json:"duration,omitempty"`
	DisableNotification *int    `json:"disable_notification,omitempty"`
	ReplyToMessageID    *int    `json:"reply_to_message_id,omitempty"`
}

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
