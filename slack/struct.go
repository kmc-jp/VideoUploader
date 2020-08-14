package slack

//Webhook put webhook message data
type Webhook struct {
	Attachments []Attachment `json:"attachments,omitempty"`

	Channel  string `json:"channel,omitempty"`
	UserName string `json:"username,omitempty"`

	Text string `json:"text,omitempty"`

	IconEmoji string `json:"icon_emoji,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
}

//Attachment put slack attachment data
type Attachment struct {
	Fallback  string `json:"fallback,omitempty"`
	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`

	Text    string `json:"text,omitempty"`
	Pretext string `json:"pretext,omitempty"`
	Color   string `json:"color,omitempty"`

	Fields []Field `json:"fields,omitempty"`

	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`

	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`

	Footer     string `json:"footer,omitempty"`
	FooterIcon string `json:"footer_icon,omitempty"`

	TS string `json:"ts,omitempty"`
}

//Field put slack field data
type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}
