package models

type TextLine struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (t *TextLine) CreateLineMessageText(text string) {
	t.Type = "text"
	t.Text = text
}

type CardLine struct {
	Type     string       `json:"type"`
	AltText  string       `json:"altText"`
	Template TemplateLine `json:"template"`
}

type TemplateLine struct {
	Type              string       `json:"type"`
	Actions           []ActionLine `json:"actions"`
	ThumbnailImageUrl string       `json:"thumbnailImageUrl"`
	Title             string       `json:"title"`
	Text              string       `json:"text"`
}

type ActionLine struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Uri   string `json:"uri"`
}

func (c *CardLine) CreateCardLine(Uri, Title, Text string) {
	action := ActionLine{Type: "uri", Label: "More Detail", Uri: Uri}
	actions := []ActionLine{action}
	template := TemplateLine{Type: "buttons", Actions: actions, ThumbnailImageUrl: "https://sv1.picz.in.th/images/2020/03/01/xXntke.jpg", Title: Title, Text: Text}
	c.Type = "template"
	c.AltText = Title
	c.Template = template
}

// {
// 	"type": "template",
// 	"altText": "this is a buttons template",
// 	"template": {
// 	  "type": "buttons",
// 	  "actions": [
// 		{
// 		  "type": "uri",
// 		  "label": "News Detail",
// 		  "uri": "https://www.google.com"
// 		}
// 	  ],
// 	  "thumbnailImageUrl": "SPECIFY_YOUR_IMAGE_URL",
// 	  "title": "sdasdasd",
// 	  "text": "scasc"
// 	}
//   }
