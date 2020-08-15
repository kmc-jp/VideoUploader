package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"../lib"
	"../slack"
)

// AddComment append comment
func AddComment(w http.ResponseWriter, r *http.Request) {
	var S func(w http.ResponseWriter, status int) = func(w http.ResponseWriter, status int) {
		res, _ := json.Marshal(Status{Status: status})
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		return
	}

	r.ParseForm()

	var Comment string = r.FormValue("Comment")
	var VideoID string = r.FormValue("VideoID")

	if Comment == "" {
		S(w, 400)
		return
	}

	var user lib.User
	err := user.Get(os.Getenv("REMOTE_USER"))
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)
		S(w, 500)
		return
	}

	if !lib.VideoExistance(VideoID) {
		S(w, 400)
		return
	}

	comments, err := lib.ReadComments(VideoID)
	if err != nil {
		S(w, 500)
		lib.Logger(err)
		slack.SendError(err)
		return
	}

	if user.Icon == "" {
		os.Mkdir(filepath.Join("Usericon"), 0777)
		{
			f, _ := os.Open(filepath.Join("static", "icon", "user.png"))
			ioutil.WriteFile(filepath.Join("Usericon", user.Name+".png"), bufio.NewScanner(f).Bytes(), 0777)
		}
		user.Icon = filepath.Join("Usericon", user.Name+".png")

		user.Update()
	}

	var NewComment lib.Comment = lib.Comment{
		Comment: Comment,
		User:    user.Name,
		Icon:    user.Icon,
		Time:    time.Now(),
		ID:      lib.MakeCommentID(comments),
	}

	err = NewComment.Save(VideoID)
	if err != nil {
		S(w, 500)
		lib.Logger(err)
		slack.SendError(err)
		return
	}

	videos, _ := lib.ReadVideoData()
	v := lib.SearchVideo(videos, VideoID)

	var upUser lib.User
	if err := upUser.Get(v.User); err != nil {
		lib.Logger(err)
		slack.SendError(err)
		S(w, 200)
		return
	}

	if upUser.Slack == "" {
		S(w, 200)
		return
	}

	var get GetQuery = make(GetQuery)

	get["Video"] = v.Video
	get["Page"] = "Play"
	get["User"] = user.Name

	err = slack.SendComment(NewComment, v, upUser.Slack, r.URL.Scheme+"://"+r.URL.Host+r.URL.Path+get.Encode())
	if err != nil {
		S(w, 200)
		lib.Logger(err)
		slack.SendError(err)
		return
	}
	S(w, 200)
	return

}

//ReadComment return json comment data
func ReadComment(w http.ResponseWriter, r *http.Request) {
	var S func(w http.ResponseWriter, status int) = func(w http.ResponseWriter, status int) {
		res, _ := json.Marshal(Status{Status: status})
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		return
	}

	var VideoID string = r.URL.Query().Get("Video")

	if !lib.VideoExistance(VideoID) {
		S(w, 400)
		return
	}

	cs, err := lib.ReadComments(VideoID)
	if err != nil {
		S(w, 500)
		return
	}

	res, _ := json.Marshal(
		struct {
			Status   int           `json:"status"`
			Comments []lib.Comment `json:"comments"`
		}{200, cs},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return

}
