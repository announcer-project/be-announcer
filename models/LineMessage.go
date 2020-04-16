package models

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (t *Text) CreateLineMessageText(text string) {
	t.Type = "text"
	t.Text = text
}
