package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/tarvitz/dnq/pkg/ogg"
	"sigs.k8s.io/yaml"

	"github.com/tarvitz/dnq/cmd/common"
	"github.com/tarvitz/dnq/pkg/telegram"
)

type UploadPositional struct {
	Filename flags.Filename `positional-arg-name:"file" env:"DNQ_UPLOAD_FILE" description:"upload file name."`
}

type UploadCommand struct {
	common.Auth      `group:"auth options"`
	UploadPositional `positional-args:"yes" required:"yes"`

	ChatID  string   `short:"C" long:"chat-id" env:"DNQ_CHAT_ID" required:"true" description:"chat id: unique int id or @username."`
	Caption string   `short:"c" long:"caption" description:"set caption for the uploaded file."`
	Matches []string `short:"m" long:"matches" description:"a keyword matching list"`
	Output  string   `short:"o" long:"output" choice:"config" choice:"json" choice:"yaml" default:"config" env:"DNQ_OUTPUT" description:"output time once file is uploaded."`
}

func (command *UploadCommand) print(message *telegram.Message) {
	var content []byte
	switch command.Output {
	case "json":
		content, _ = json.MarshalIndent(message, "", "    ")
	case "yaml":
		content, _ = yaml.Marshal(message)
	case "config":
		quote := &telegram.Quote{
			ID:      message.Voice.FileID,
			Caption: command.Caption,
			Matches: command.Matches,
		}
		content, _ = yaml.Marshal([]*telegram.Quote{quote})
	default:
		fmt.Printf("%v\n", message)
	}
	fmt.Printf("%s\n", content)

}

func (command *UploadCommand) upload(file *os.File) (err error) {
	var message *telegram.Message
	client := command.GetClient()
	message, err = client.Upload(telegram.SendVoice, map[string]io.Reader{
		"voice":   file,
		"chat_id": strings.NewReader(command.ChatID),
		"caption": strings.NewReader(command.Caption),
	})
	if err != nil {
		return
	}
	command.print(message)
	return
}

func (command *UploadCommand) Execute(_ []string) (err error) {
	var fd *os.File
	if fd, err = os.Open(string(command.Filename)); err != nil {
		return
	}
	defer func() { _ = fd.Close() }()

	if ogg.IsOggOpusFile(fd) {
		_, _ = fd.Seek(0, 0)
		err = command.upload(fd)
	} else {
		err = fmt.Errorf("you can't downloaded this file")
		fmt.Printf("%v\n", err)
	}
	return
}
