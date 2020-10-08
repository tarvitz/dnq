package telegram

import "testing"

func TestReadQuotes(t *testing.T) {
	originQuotes := Quotes
	defer func() { Quotes = originQuotes }()

	t.Run("ok", func(in *testing.T) {
		if err := ReadQuotes(testConfigFile); err != nil {
			in.Errorf("got error: %v", err)
		}
	})

	t.Run("failure/config-file-not-found", func(in *testing.T) {
		if err := ReadQuotes("non existent"); err == nil {
			in.Error("err expected, got nil instead.")
		}
	})

	t.Run("failure/bad-config", func(in *testing.T) {
		if err := ReadQuotes(testBadConfigFile); err == nil {
			in.Error("err expected, got nil instead.")
		}
	})
}
