package server

import (
	"fmt"
	"net/http"

	"github.com/jessevdk/go-flags"
	"golang.org/x/net/http2"

	"github.com/tarvitz/dnq/cmd/common"
	"github.com/tarvitz/dnq/pkg/config"
	"github.com/tarvitz/dnq/pkg/telegram"
)

// Command +
type Command struct {
	common.Auth `group:"auth" description:"authentication related options"`

	Port    int    `short:"p" long:"port" default:"8443" env:"DNQ_SERVER_PORT" description:"server port"`
	Address string `short:"H" long:"host" default:"0.0.0.0" env:"DNQ_SERVER_HOST" description:"server address"`

	Config flags.Filename `short:"c" long:"config" env:"DNQ_CONFIG" default:"config.yaml" description:"configuration file"`

	// TSL options
	Cert string `long:"cert" env:"DNQ_SERVER_CERT" default:"resources/server.pem" description:"server.{pem,crt} file"`
	Key  string `long:"key" env:"DNQ_SERVER_KEY" default:"resources/server.key" description:"server.key file"`

	config *config.Config
}

// Listen function that runs sever, could be overloaded by testing needs.
var Listen = func(srv *http.Server, cert, key string) (err error) {
	return srv.ListenAndServeTLS(cert, key)
}

func (command *Command) initServer() {
	httpServer = &http.Server{
		Addr: fmt.Sprintf("%s:%d", command.Address, command.Port),
	}
	_ = http2.ConfigureServer(httpServer, http2Server)
}

// Execute +
func (command *Command) Execute(_ []string) (err error) {
	// assign self to module variable
	cmd = command
	command.initServer()

	if cmd.config, err = config.ReadConfig(string(command.Config)); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Quotes loaded: %d\n", len(cmd.config.Quotes))
	telegram.SetQuotes(cmd.config.Quotes)
	return Listen(httpServer, command.Cert, command.Key)
}
