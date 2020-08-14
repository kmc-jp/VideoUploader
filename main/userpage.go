package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"../lib"
)

//UserPageHandle handle user page
func UserPageHandle(w http.ResponseWriter, r *http.Request) {
	var SelectUser lib.User
	if e := SelectUser.Get(r.URL.Query().Get("User")); e != nil {
		SelectUser = lib.User{
			Name:  r.URL.Query().Get("User"),
			Video: []lib.Video{},
		}
	}

	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		user = lib.User{
			Name:  os.Getenv("REMOTE_USER"),
			Video: []lib.Video{},
		}
	}

	var U UserPage = UserPage{
		Header: Header{
			User:     user,
			Title:    SelectUser.Name + "のお部屋",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Slack:  Settings.SlackWebhook != "",
		User:   SelectUser,
		Footer: Footer{},
	}

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "userpage.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	if e = t.ExecuteTemplate(w, "user", U); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	return
}
