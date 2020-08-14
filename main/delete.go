package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"../lib"
)

//Delete delete put video
func Delete(w http.ResponseWriter, r *http.Request) {
	var getQuery GetQuery = make(GetQuery)
	getQuery["Page"] = "MyPage"

	// error handling
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(w, "Panic: %v\n", err)
		}
	}()

	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		getQuery["Error"] = "NotFound"
		lib.Logger(e)
		w.Header().Set("Location", "index.up"+getQuery.Encode())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	var videoData = lib.SearchVideo(user.Video, r.URL.Query().Get("Video"))

	var D func(V []lib.Video) []lib.Video = func(V []lib.Video) []lib.Video {
		var video = lib.SearchVideo(V, r.URL.Query().Get("Video"))
		if video.Video == "" {
			getQuery["Error"] = "NotFound"
			w.Header().Set("Location", "index.up"+getQuery.Encode())
			w.WriteHeader(http.StatusTemporaryRedirect)
			return V
		}

		v, err := lib.DeleteVideo(V, video)
		if err != nil {
			getQuery["Error"] = "DeleteError"
			w.Header().Set("Location", "index.up"+getQuery.Encode())
			w.WriteHeader(http.StatusTemporaryRedirect)
		}

		return v
	}

	user.Video = D(user.Video)

	var data map[string]lib.User

	bData, err := ioutil.ReadFile(UserInfoFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(bData, &data)
	if err != nil {
		return
	}

	data[user.Name] = user

	bData, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return
	}

	ioutil.WriteFile(UserInfoFile, bData, 0777)

	bData, err = ioutil.ReadFile(AllVideosFile)
	if err != nil {
		return
	}

	var V []lib.Video
	err = json.Unmarshal(bData, &V)
	if err != nil {
		return
	}

	V = D(V)

	bData, err = json.MarshalIndent(V, "", "    ")
	if err != nil {
		return
	}

	ioutil.WriteFile(AllVideosFile, bData, 0777)

	os.Remove(filepath.Join("Videos", videoData.Video+".png"))
	os.Remove(filepath.Join("Videos", videoData.Video))

	getQuery["Success"] = "delete"
	w.Header().Set("Location", "index.up"+getQuery.Encode())
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}
