package telegram

var Quotes map[string][]*Quote

type Quote struct {
	ID      string   `json:"id"`
	Caption string   `json:"caption"`
	Matches []string `json:"matches"`
}

func SetQuotes(quotes []*Quote) {
	Quotes = make(map[string][]*Quote)
	for _, quote := range quotes {
		for _, match := range quote.Matches {
			Quotes[match] = append(Quotes[match], quote)
		}
	}
}
