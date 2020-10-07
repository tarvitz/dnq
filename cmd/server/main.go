package server

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/tarvitz/dnq/cmd/common"
	"github.com/tarvitz/dnq/pkg/telegram"
	"net/http"

	"golang.org/x/net/http2"
)

// Command +
type Command struct {
	common.Auth `group:"auth" description:"authentication related options"`

	Port    int    `short:"p" long:"port" default:"8443" env:"DNQ_SERVER_PORT" description:"server port"`
	Address string `short:"H" long:"host" default:"" env:"DNQ_SERVER_HOST" description:"server address"`

	Config flags.Filename `short:"c" long:"config" env:"DNQ_CONFIG" default:"config.yaml" description:"configuration file"`

	// TSL options
	Cert string `long:"cert" env:"DNQ_SERVER_CERT" default:"resources/server.pem" description:"server.{pem,crt} file"`
	Key  string `long:"key" env:"DNQ_SERVER_KEY" default:"resources/server.key" description:"server.key file"`
}

// Listen function that runs sever, could be overloaded by testing needs.
var Listen = func(srv *http.Server, cert, key string) (err error) {
	return srv.ListenAndServeTLS(cert, key)
}

func (command *Command) initServer() {
	httpServer = &http.Server{Addr: fmt.Sprintf(":%d", command.Port)}
	_ = http2.ConfigureServer(httpServer, http2Server)
}

// Execute +
func (command *Command) Execute(_ []string) (err error) {
	// assign self to module variable
	cmd = command
	// init
	command.initServer()
	if err = telegram.ReadQuotes(string(command.Config)); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	return Listen(httpServer, command.Cert, command.Key)
}
