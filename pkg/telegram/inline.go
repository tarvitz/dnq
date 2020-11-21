package telegram

import (
	"bytes"
	"encoding/json"
)

// InlineService is for service requests for inline queries
// See also: https://core.telegram.org/bots/api#inline-mode
type InlineService struct {
	client *Client
}

// AnswerInlineQuery sends to the requester a response for a inline request.
func (service *InlineService) AnswerInlineQuery(update *Update) (err error) {
	answer := NewAnswerInline(update)
	content, _ := json.Marshal(answer)
	buffer := bytes.NewBuffer(content)

	request := &APIRequest{
		Method:    "POST",
		APIMethod: AnswerInlineQuery,
		Body:      buffer,
	}
	err = service.client.Call(request, nil)
	return
}
