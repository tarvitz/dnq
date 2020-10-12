package common

import (
	"fmt"

	"github.com/tarvitz/dnq/pkg/telegram"
)

type Auth struct {
	Token string `short:"t" long:"token" required:"true" env:"DNQ_API_TOKEN" description:"bot api token"`

	client *telegram.Client
}

// Method returns telegram API method name. Example
//   Method("getUpdates")
func (command *Auth) Method(name string) string {
	return fmt.Sprintf("%s%s/%s", telegram.BotAPIURL, command.Token, name)
}

func (command *Auth) GetClient() *telegram.Client {
	if command.client == nil {
		command.client = telegram.NewClient(command.Token)
	}
	return command.client
}

func (command *Auth) SetClient(client *telegram.Client) {
	command.client = client
}
