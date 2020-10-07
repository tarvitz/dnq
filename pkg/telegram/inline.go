package telegram

import (
	"bytes"
	"encoding/json"
)

type InlineService struct {
	client *Client
}

func (service *InlineService) AnswerInlineQuery(update *Update) (err error) {
	answer := NewAnswerInline(update)
	content, _ := json.Marshal(answer)
	buffer := bytes.NewBuffer(content)

	request := &APIRequest{
		Method:    "POST",
		APIMethod: "answerInlineQuery",
		Body:      buffer,
	}
	err = service.client.Call(request, nil)
	return
}
