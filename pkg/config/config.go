package config

import (
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/tarvitz/dnq/pkg/telegram"
)

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
