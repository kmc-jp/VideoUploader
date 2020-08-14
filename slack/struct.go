package slack

//Webhook put webhook message data
type Webhook struct {
	Attachments []Attachment `json:"attachments"`

	Channel  string `json:"channel"`
	UserName string `json:"username"`

	Text string `json:"text"`

	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
}

//Attachment put slack attachment data
type Attachment struct {
	Fallback  string `json:"fallback"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`

	Text    string `json:"text"`
	Pretext string `json:"pretext"`
	Color   string `json:"color"`

	Fields []Field `json:"fields"`

	AuthorName string `json:"author_name"`
	AuthorLink string `json:"author_link"`

	ImageURL string `json:"image_url"`
	ThumbURL string `json:"thumb_url"`

	Footer     string `json:"footer"`
	FooterIcon string `json:"footer_icon"`

	TS string `json:"ts"`
}

//Field put slack field data
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
