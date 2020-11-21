package config

import (
	"io/ioutil"
	"math/rand"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/tarvitz/dnq/pkg/telegram"
)

const (
	wtfQuoteID = "AwACAgIAAxkDAAMWX3okXL1AZ-aOQTpL2t-7tExt2YIAArMIAAJgG9BLmWJEVGtI5hwbBA"
)

// Config stores DNQ configuration assets, quotes, tokens and so forth
type Config struct {
	// Administrator's token to perform s2s queries from administrator privileges
	AdminToken string `json:"admin-token"`
	// Quotes configuration
	Quotes []*telegram.Quote `json:"quotes"`
}

// ReadConfig read Duke Nukem's quotes from configuration file.
func ReadConfig(filename string) (config *Config, err error) {
	var fd *os.File

	fd, err = os.Open(filename)
	if err != nil {
		return
	}
	defer func() { _ = fd.Close() }()

	content, _ := ioutil.ReadAll(fd)
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return
	}
	return
}

// RandomQuote returns a random quote from the quotes slice, otherwise
// returns pre-set one.
func (config *Config) RandomQuote() (quote *telegram.Quote) {
	if len(config.Quotes) > 0 {
		idx := rand.Int() % len(config.Quotes)
		return config.Quotes[idx]
	}
	//: fallback quote
	return &telegram.Quote{
		ID:      wtfQuoteID,
		Caption: "What the ..?",
		Matches: nil,
	}
}
