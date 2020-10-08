package common

import "fmt"

type Auth struct {
	Token string `short:"t" long:"token" required:"true" env:"DNQ_API_TOKEN" description:"bot api token"`
}

// Method returns telegram API method name. Example
//   Method("getUpdates")
func (command *Auth) Method(name string) string {
	return fmt.Sprintf("%s%s/%s", TelegramBotAPIURL, command.Token, name)
}
