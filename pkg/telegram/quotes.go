package telegram

import (
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

var Quotes map[string][]*Quote

type Config struct {
	Quotes []*Quote `json:"quotes"`
}

type Quote struct {
	ID      string   `json:"id"`
	Caption string   `json:"caption"`
	Matches []string `json:"matches"`
}

// ReadQuotes read Duke Nukem's quotes from configuration file.
func ReadQuotes(filename string) (err error) {
	var (
		fd     *os.File
		config *Config
	)

	fd, err = os.Open(filename)
	if err != nil {
		return
	}
	defer Close(fd)

	content, _ := ioutil.ReadAll(fd)
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return
	}

	Quotes = make(map[string][]*Quote)
	for _, quote := range config.Quotes {
		for _, match := range quote.Matches {
			Quotes[match] = append(Quotes[match], quote)
		}
	}
	return
}
