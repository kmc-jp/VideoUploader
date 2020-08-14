package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"../lib"
	"../slack"
)

//UploadHandle handle uploaded video file
func UploadHandle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(w, "Panic!:%v", err)
		}
	}()

	var getQuery GetQuery = make(GetQuery)

	// postの内容をパース
	r.ParseForm()

	// ユーザデータの取得
	var newData lib.Video
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		user = lib.User{
			Name:  os.Getenv("REMOTE_USER"),
			Video: []lib.Video{},
		}
	}

	// 動画データを変数に格納
	video, videoH, err := r.FormFile("video")
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)
		lib.Progress(newData, lib.Status{Error: fmt.Sprintf("%s", err.Error())})
		lib.TmpClear(newData.Video)
		getQuery["Error"] = "UploadError"
		return
	}

	// サムネイルデータを変数に格納
	var thumbB bool
	thumb, thumbH, err := r.FormFile("thumbnail")
	if err == nil && thumbH.Filename != "" {
		thumbB = true
	}

	// タグデータを変数に格納
	var tags = lib.SplitTags(r.FormValue("tag"))
	tags = lib.TrimSpaces(tags)

	newData.Tags = tags
	newData.Video = lib.MakeSuitName() + Extension
	newData.Time = time.Now()
	newData.User = user.Name
	newData.Title = r.FormValue("title")

	if newData.Title == "" {
		newData.Title = videoH.Filename
	}

	// make index
	newData.Thumb = newData.Video + ".png"
	if cap(user.Video) == 0 {
		user.Video = append(user.Video, newData)
	} else {
		user.Video = append([]lib.Video{newData}, user.Video[0:]...)
	}

	// export index
	if err = user.WriteAll(); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		lib.Logger(err)
		slack.SendError(err)
		return
	}

	//export taglist
	if err = lib.TagVideoAppend(newData); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		lib.Logger(err)
		slack.SendError(err)
		return
	}

	// 一時領域にディレクトリ作成
	os.Mkdir(filepath.Join("tmp"), 0777)
	os.Mkdir(filepath.Join("tmp", newData.Video), 0777)

	// 動画データを一時領域に格納
	lib.Progress(newData, lib.Status{Phase: "exporting"})

	_, extension := lib.FindExtension(videoH.Filename)
	file, err := os.Create(filepath.Join("tmp", newData.Video, "video"+"."+extension))
	if err != nil {
		lib.Progress(newData, lib.Status{Error: fmt.Sprintf("%s", err.Error())})
		return
	}

	if _, err = io.Copy(file, video); err != nil {
		lib.Logger(err)
		slack.SendError(err)
		lib.Progress(newData, lib.Status{Error: fmt.Sprintf("%s", err.Error())})
		return
	}

	file.Close()

	// サムネイルを指定されていた場合はpng形式にして保存
	lib.Progress(newData, lib.Status{Phase: "saving thumbnail"})

	if thumbB {
		if err := lib.ImageConverter(thumb, filepath.Join("Videos", newData.Video, newData.Thumb)); err != nil {
			thumbB = false
			goto next
		}
	}

next:
	// エンコード指示
	lib.Progress(newData, lib.Status{Phase: "calling encode process..."})

	getQuery["Success"] = "upload"
	w.Header().Set("Location", "index.up"+getQuery.Encode())
	w.WriteHeader(http.StatusTemporaryRedirect)

}
