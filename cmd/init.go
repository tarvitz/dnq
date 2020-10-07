package main

import (
	"github.com/tarvitz/dnq/cmd/server"
	"github.com/tarvitz/dnq/cmd/tools"
)

func init() {
	for _, command := range []struct {
		name             string
		shortDescription string
		longDescription  string
		cmd              interface{}
	}{
		{
			name:             "server",
			shortDescription: "bot http webhook",
			longDescription:  "",
			cmd:              &server.Command{},
		},
		{
			name:             "upload",
			shortDescription: "upload an ogg opus encoded file like",
			longDescription:  "",
			cmd:              &tools.UploadCommand{},
		},
	} {
		_, _ = parser.AddCommand(
			command.name,
			command.shortDescription,
			command.longDescription,
			command.cmd)
	}
}
