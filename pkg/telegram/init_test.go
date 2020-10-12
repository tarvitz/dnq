package telegram

// pre-set data for testing
func init() {
	SetQuotes([]*Quote{
		{"AwA--1", "Good, bad, I'm a guy with a gun", []string{"good", "bad"}},
		{"AwA--2", "Groovy", []string{"groovy"}},
		{"AwA--3", "Come get some!", []string{"", "come", "get", "some"}},
	})
}
