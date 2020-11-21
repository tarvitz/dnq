package config

import (
	"fmt"
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

func TestConfig_RandomQuote(t *testing.T) {
	t.Run("blank", func(in *testing.T) {
		config := &Config{}
		quote := config.RandomQuote()
		if quote.ID != wtfQuoteID {
			in.Errorf("expected: %v, got: %v", wtfQuoteID, quote.ID)
		}
	})

	t.Run("random", func(in *testing.T) {
		var quotes []*telegram.Quote
		for i := 0; i < 100; i++ {
			quotes = append(quotes, &telegram.Quote{ID: fmt.Sprintf("%d", i)})
		}
		config := &Config{Quotes: quotes}
		quote := config.RandomQuote()
		if quote.ID == wtfQuoteID {
			in.Errorf("expected not `%s` but got it", wtfQuoteID)
		}
	})
}
