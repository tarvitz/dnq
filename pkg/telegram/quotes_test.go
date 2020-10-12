package telegram

import "testing"

func TestSetQuotes(t *testing.T) {
	quotes := []*Quote{
		{"1", "caption 1", []string{"1", "2", "3"}},
		{"2", "caption 2", []string{"2", "3", "4"}},
		{"3", "caption 3", []string{"1337", "test", ""}},
	}
	SetQuotes(quotes)
}
