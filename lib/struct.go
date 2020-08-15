package lib

import "time"

//Setting put setting data
type Setting struct {
	FFmpeg       string `json:"ffmpeg"`
	FFprobe      string `json:"ffprobe"`
	SlackWebhook string `json:"slack_webhook"`
	ErrorChannel string `json:"error_channel"`
	GyazoToken   string `json:"gyazo_token"`
}

//Video put video info
type Video struct {
	Video  string    `json:"video"`
	User   string    `json:"user"`
	Title  string    `json:"title"`
	Time   time.Time `json:"time"`
	Thumb  string    `json:"thumb"`
	Status Status    `json:"status"`
	Tags   []string  `json:"tags"`
}

// Comment put comment data
type Comment struct {
	Comment string    `json:"comment"`
	User    string    `json:"user_name"`
	Icon    string    `json:"user_icon"`
	Time    time.Time `json:"time"`
	ID      string    `json:"comment_id"`
}

//Status put uploading phase and errors
type Status struct {
	Phase string
	Error string
}

//User put user data
type User struct {
	Name  string  `json:"username"`
	Video []Video `json:"files"`
	Icon  string  `json:"icon"`
	Slack string
}

//Tag put tag data
type Tag struct {
	Tag map[string][]string
}
