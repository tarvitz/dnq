package config

import (
	"testing"

	"github.com/tarvitz/dnq/pkg/telegram"
)

const (
	testConfigFile    = "resources/config.yaml"
	testBadConfigFile = "resources/bad-config.yaml"
)

func TestReadConfig(t *testing.T) {
	originQuotes := telegram.Quotes
	defer func() { telegram.Quotes = originQuotes }()

	t.Run("ok", func(in *testing.T) {
		if _, err := ReadConfig(testConfigFile); err != nil {
			in.Errorf("got error: %v", err)
		}
	})

	t.Run("failure/config-file-not-found", func(in *testing.T) {
		if _, err := ReadConfig("non existent"); err == nil {
			in.Error("err expected, got nil instead.")
		}
	})

	t.Run("failure/bad-config", func(in *testing.T) {
		if _, err := ReadConfig(testBadConfigFile); err == nil {
			in.Error("err expected, got nil instead.")
		}
	})
}
