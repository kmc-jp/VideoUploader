package main

import (
	"../lib"
)

//GetQuery handle query data
type GetQuery map[string]string

//Header put header objects
type Header struct {
	Title    string
	UserName string
	Error    string
	Success  string
	User     lib.User
}

//Footer put header objects
type Footer struct {
}

//Index put index objects
type Index struct {
	Header Header
	Footer Footer
	Video  []lib.Video
}

//Upload put upload page data
type Upload struct {
	Header Header
	Footer Footer
}

//Player put player page data
type Player struct {
	Header  Header
	Footer  Footer
	Comment []lib.Comment
	Video   lib.Video
	User    lib.User
	Slack   bool
}

//MyPage put mypage data
type MyPage struct {
	Header Header
	Footer Footer
	Video  []lib.Video
	Slack  bool
}

//UserListPage put user list page data
type UserListPage struct {
	Header Header
	Footer Footer
	User   []lib.User
}

//UserPage put user page data
type UserPage struct {
	Header Header
	Footer Footer
	User   lib.User
	Slack  bool
}

//SettingsPage put settings page data
type SettingsPage struct {
	Header Header
	Footer Footer
}

// Status http status
type Status struct {
	Status int `json:"status"`
}

//Script put index objects
type Script struct {
}
