package main

import (
	"fmt"
	"net/http"
	"os"

	"../lib"
	"../slack"
)

// SendSlack send video info for slack
func SendSlack(w http.ResponseWriter, r *http.Request) {
	var user lib.User
	user.Get(r.URL.Query().Get("User"))

	var video = lib.SearchVideo(user.Video, r.URL.Query().Get("Video"))

	var channel string = r.URL.Query().Get("Channel")
	var URL string = r.URL.Scheme + "://" + r.URL.Host + r.URL.Path
	var get GetQuery = make(GetQuery)

	get["Video"] = video.Video
	get["Page"] = "Play"
	get["User"] = user.Name

	var getQuery GetQuery = make(GetQuery)

	getQuery["Page"] = r.URL.Query().Get("Page")
	switch getQuery["Page"] {
	case "Play":
		getQuery["Video"] = video.Video
		getQuery["User"] = user.Name
	case "UserPage":
		getQuery["User"] = user.Name
	}

	err := slack.SendVideoInfo(video, channel, URL+get.Encode())
	if err != nil {
		getQuery["Error"] = "SendSlack"
		lib.Logger(err)
		slack.SendError(err)
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	getQuery["Success"] = "SendSlack"

	w.Header().Set("Location", "index.up"+getQuery.Encode())
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}

// SetSlackChannel set user slack channel
func SetSlackChannel(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(w, "Panic!:%v", err)
		}
	}()
	var getQuery GetQuery = make(GetQuery)

	getQuery["Video"] = r.URL.Query().Get("Video")
	getQuery["Page"] = r.URL.Query().Get("Page")
	getQuery["User"] = r.URL.Query().Get("User")

	r.ParseForm()
	var channel = r.FormValue("channel")

	var user lib.User
	user.Get(os.Getenv("REMOTE_USER"))

	user.Slack = channel

	err := user.Update()
	if err != nil {
		getQuery["Error"] = "SetChannel"
		lib.Logger(err)
		slack.SendError(err)
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	getQuery["Success"] = "SetChannel"
	w.Header().Set("Location", "index.up"+getQuery.Encode())
	w.WriteHeader(http.StatusTemporaryRedirect)

	return
}
