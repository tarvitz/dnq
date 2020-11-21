package telegram

// Quotes is a map for keywords and quotes
var Quotes map[string][]*Quote

// Quote keeps telegram voice message id (ID), caption and
// matching keywords
type Quote struct {
	ID      string   `json:"id"`
	Caption string   `json:"caption"`
	Matches []string `json:"matches"`
}

// SetQuotes processes Quote object slice (generally taken from
// a configuration file) and build Quotes map with it.
func SetQuotes(quotes []*Quote) {
	Quotes = make(map[string][]*Quote)
	for _, quote := range quotes {
		for _, match := range quote.Matches {
			Quotes[match] = append(Quotes[match], quote)
		}
	}
}
