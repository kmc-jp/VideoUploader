package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"../lib"
)

//PlayPageHandle show player page
func PlayPageHandle(w http.ResponseWriter, r *http.Request) {
	var Play Player = Player{}
	var user lib.User
	var videoUser lib.User

	user.Get(os.Getenv("REMOTE_USER"))
	videoUser.Get(r.URL.Query().Get("User"))

	Play.Video = lib.SearchVideo(videoUser.Video, r.URL.Query().Get("Video"))
	Play.User = videoUser
	Play.Slack = Settings.SlackWebhook != ""

	Play.Header = Header{
		Title:    Play.Video.Title,
		UserName: os.Getenv("REMOTE_USER"),
		User:     user,
		Error:    ErrorHandle(r.URL.Query().Get("Error")),
		Success:  SuccessHandle(r.URL.Query().Get("Success")),
	}

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "play.html"),
		filepath.Join("resources", "script.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	if e = t.ExecuteTemplate(w, "play", Play); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	return
}
