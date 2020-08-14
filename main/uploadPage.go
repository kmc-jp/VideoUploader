package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"../lib"
)

//UploadPageHandle handle uploaded video file
func UploadPageHandle(w http.ResponseWriter, r *http.Request) {
	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		user = lib.User{
			Name:  os.Getenv("REMOTE_USER"),
			Video: []lib.Video{},
		}
	}
	var T Upload = Upload{
		Header: Header{
			User:     user,
			Title:    "Upload",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    "",
			Success:  "",
		},
		Footer: Footer{},
	}

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "upload.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)
	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	if e = t.ExecuteTemplate(w, "upload", T); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
}
