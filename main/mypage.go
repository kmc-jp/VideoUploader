package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../lib"
	"../slack"
)

//MyPageHandle handle my page
func MyPageHandle(w http.ResponseWriter, r *http.Request) {
	var user lib.User

	var V MyPage = MyPage{
		Header: Header{
			Title:    "MyPage",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Slack:  Settings.SlackWebhook != "",
		Footer: Footer{},
	}

	if err := user.Get(os.Getenv("REMOTE_USER")); err != nil {
		V.Video = []lib.Video{}
	} else {
		V.Video = user.Video
	}

	V.Header.User = user

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "mypage.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	if e = t.ExecuteTemplate(w, "mypage", V); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	return

}

//UpdateVideoInfo update user video infomation
func UpdateVideoInfo(w http.ResponseWriter, r *http.Request) {
	var getQuery GetQuery = make(GetQuery)
	getQuery["Page"] = "MyPage"

	// error handling
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(w, "Panic: %v\n", err)
		}
	}()

	r.ParseForm()

	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		getQuery["Error"] = "NotFound"
		lib.Logger(e)
		slack.SendError(e)
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	var video = lib.SearchVideo(user.Video, r.URL.Query().Get("Video"))
	if video.Video == "" {
		getQuery["Error"] = "NotFound"
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// update thumbnail
	func() {
		thumb, thumbH, err := r.FormFile("thumbnail")
		if err != nil || thumbH.Filename == "" {
			return
		}

		if err := lib.ImageConverter(thumb, filepath.Join("Videos", video.Video+".png")); err != nil {
			lib.Logger(err)
			slack.SendError(err)
			return
		}

	}()

	// update title
	func() {
		var title = r.FormValue("title")
		var newTitle string = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")

		if strings.ReplaceAll(strings.TrimSpace(newTitle), "\n", "") == "" {
			return
		}

		video.Title = newTitle
		return
	}()

	//update tags
	func() {
		var tags = lib.SplitTags(r.FormValue("tag"))
		if len(tags) == 0 {
			return
		}

		if err := lib.TagUpdate(video, tags); err != nil {
			lib.Logger(err)
			slack.SendError(err)

			return
		}

		video.Tags = tags

	}()

	var err = video.Update()
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)
		getQuery["Error"] = "UpdateError"
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	getQuery["Success"] = "update"

	w.Header().Set("Location", "index.up"+getQuery.Encode())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

//SendVideoStatus make json response of the video encoding status
func SendVideoStatus(w http.ResponseWriter, r *http.Request) {
	var S func(w http.ResponseWriter, status int) = func(w http.ResponseWriter, status int) {
		res, _ := json.Marshal(Status{Status: status})
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
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

	type Res struct {
		Status int            `json:"status"`
		Time   string         `json:"time"`
		Video  []lib.ResVideo `json:"video"`
	}

	var initTime time.Time
	if user.Time.String() == initTime.String() {
		user.Time = time.Now()
	}

	var vs []lib.ResVideo = lib.VideosToResVideos(user.Video)
	res, _ := json.Marshal(
		Res{
			Status: 200,
			Video:  vs,
			Time:   user.Time.Format("2006-01-02T15:04:05+09:00"),
		},
	)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}
