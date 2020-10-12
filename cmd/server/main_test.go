package server

import (
	"net/http"
	"testing"
)

const (
	testConfigFile = "../../config.yaml"
)

func TestCommand_Execute(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		origListen := Listen
		defer func() { Listen = origListen }()
		Listen = func(srv *http.Server, cert, key string) (err error) {
			return
		}
		cmd := &Command{
			Config: testConfigFile,
		}
		err := cmd.Execute([]string{})
		if err != nil {
			in.Errorf("got error")
		}
	})

	t.Run("failure/cant-read-config", func(in *testing.T) {
		cmd := &Command{Config: "non existent file"}
		err := cmd.Execute([]string{})
		if err == nil {
			in.Errorf("expected error, got nil instead")
		}
	})
}

func TestCommand_initServer(t *testing.T) {
	//: just register a call
	origHttpServer := httpServer
	defer func() { httpServer = origHttpServer }()

	cmd := &Command{
		Port:    1337,
		Address: "0.0.0.0",
	}
	cmd.initServer()
	expected := "0.0.0.0:1337"
	if httpServer.Addr != expected {
		t.Errorf("expected: %v, got: %v", expected, httpServer.Addr)
	}
}
