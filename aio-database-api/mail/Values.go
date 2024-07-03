package mail

type Mail struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Text        string `json:"text"`
	Button      bool   `json:"button"`
	ButtonUrl   string `json:"button_url"`
	ButtonText  string `json:"button_text"`
	BelowButton string `json:"below_button"`
}
