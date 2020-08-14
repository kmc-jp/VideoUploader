package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"../lib"
)

//IndexPageHandle show index page
func IndexPageHandle(w http.ResponseWriter, r *http.Request) {
	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
	}
	var IndexPage Index = Index{
		Header: Header{
			User:     user,
			Title:    "",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Footer: Footer{},
	}

	var Videos, err = lib.ReadVideoData()
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	if err = IndexPage.Get(); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	var IndexVideos []lib.Video

	func() {
		for _, v := range Videos {
			if v.Status.Error != "" || v.Status.Phase != "" {
				continue
			}

			IndexVideos = append(IndexVideos, v)
		}
	}()

	IndexPage.Video = IndexVideos

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "index.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	if e = t.ExecuteTemplate(w, "index", IndexPage); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	return
}
