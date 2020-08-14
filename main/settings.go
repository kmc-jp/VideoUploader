package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"../lib"
	"../slack"
)

//SettingsPageHandle hadle settings page
func SettingsPageHandle(w http.ResponseWriter, r *http.Request) {
	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		user = lib.User{
			Name:  os.Getenv("REMOTE_USER"),
			Video: []lib.Video{},
		}
	}
	var Page SettingsPage = SettingsPage{
		Header: Header{
			User:     user,
			Title:    "Settings",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Footer: Footer{},
	}

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "settings.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	if e = t.ExecuteTemplate(w, "settings", Page); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	return
}

//SetUserIcon Set user icon
func SetUserIcon(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(w, "Panic!:%v", err)
		}
	}()

	var getQuery GetQuery = make(GetQuery)
	getQuery["Page"] = "Settings"

	// postの内容をパース
	r.ParseForm()
	iconF, _, err := r.FormFile("icon")
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)
		getQuery["Error"] = "IconNotFound"
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	defer iconF.Close()
	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		user = lib.User{
			Name:  os.Getenv("REMOTE_USER"),
			Video: []lib.Video{},
		}
	}

	os.Mkdir(filepath.Join("Usericon"), 0777)

	{
		reader, err := lib.ImageSizer(IconSize, iconF)
		if err != nil {
			lib.Logger(err)
			slack.SendError(err)
			getQuery["Error"] = "ResizeError"
			w.Header().Set("Location", "index.up"+getQuery.Encode())
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		f, err := os.Create(filepath.Join("Usericon", user.Name+".png"))
		if err != nil {
			lib.Logger(err)
			slack.SendError(err)
			getQuery["Error"] = "ThumbWriteError"
			w.Header().Set("Location", "index.up"+getQuery.Encode())
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		defer f.Close()

		io.Copy(f, reader)
	}

	user.Icon = filepath.Join("Usericon", user.Name+".png")

	err = user.Update()
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)
		getQuery["Error"] = "UserDataUpdateError"
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	getQuery["Success"] = "SetIcon"
	w.Header().Set("Location", "index.up"+getQuery.Encode())
	w.WriteHeader(http.StatusTemporaryRedirect)

	return
}
