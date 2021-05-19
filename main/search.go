package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../lib"
	"../slack"
)

//SearchPageHandle show search page
func SearchPageHandle(w http.ResponseWriter, r *http.Request) {
	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
	}
	var SearchPage Search = Search{
		Header: Header{
			User:     user,
			Title:    "",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Footer: Footer{},
	}

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "search.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	w.Header().Add("Content-Type", "text/html")

	if e = t.ExecuteTemplate(w, "search", SearchPage); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	return
}

//ShowTagPage show page about tag
func ShowTagPage(w http.ResponseWriter, r *http.Request) {
	tag := lib.SplitTags(r.URL.Query().Get("Tag"))
	var getQuery GetQuery = make(GetQuery)
	getQuery["Page"] = "Search"

	var AllTags lib.Tag
	var err error
	if AllTags, err = AllTags.Get(); err != nil {
		lib.Logger(err)
		slack.SendError(err)
		getQuery["Error"] = "TagReadError"
		return
	}

	if len(tag) < 1 {
		return
	}

	videos, ok := AllTags.Tag[tag[0]]
	if !ok {
		lib.Logger(err)
		slack.SendError(err)
		getQuery["Error"] = "NotFound"
		return
	}

	Videos, err := lib.VideoIDsToVideos(videos)
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)
		getQuery["Error"] = "ReadVideosError"
		return
	}

	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
	}

	var SearchPage Search = Search{
		Header: Header{
			User:     user,
			Title:    "",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Footer: Footer{},
		Video:  Videos,
		Tag:    tag[0],
	}

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "search.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	w.Header().Add("Content-Type", "text/html")

	if e = t.ExecuteTemplate(w, "search", SearchPage); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}
	return

}

//SearchTags return tags which contains put words when there are multi choise and return videos when there is only one choise
func SearchTags(w http.ResponseWriter, r *http.Request) {
	var S func(w http.ResponseWriter, status int) = func(w http.ResponseWriter, status int) {
		res, _ := json.Marshal(Status{Status: status})
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		return
	}

	tags := lib.SplitTags(r.URL.Query().Get("Keywords"))
	mode := r.URL.Query().Get("Mode")

	if mode == "" {
		mode = "or"
	}

	var AllTags lib.Tag
	var err error
	if AllTags, err = AllTags.Get(); err != nil {
		lib.Logger(err)
		slack.SendError(err)
		S(w, 500)
		return
	}
	{
		if len(tags) == 1 {
			if videos, ok := AllTags.Tag[tags[0]]; ok {
				Videos, err := lib.VideoIDsToVideos(videos)
				if err != nil {
					lib.Logger(err)
					slack.SendError(err)
					S(w, 500)
					return
				}
				res, _ := json.Marshal(
					struct {
						Status int            `json:"status"`
						Type   string         `json:"type"`
						Video  []lib.ResVideo `json:"videos"`
					}{
						Status: 200,
						Type:   "videos",
						Video:  lib.ExcludeErrorAndPhaseVideos(lib.VideosToResVideos(Videos)),
					},
				)
				w.Header().Set("Content-Type", "application/json")
				w.Write(res)
				return
			}
		}
	}

	var res []byte
	var ResTags []string = make([]string, 0)
	switch mode {
	case "or":
		for key := range AllTags.Tag {
			for _, tag := range tags {
				if strings.Contains(key, tag) {
					ResTags = append(ResTags, key)
					break
				}
			}
		}

	case "and":
		for key := range AllTags.Tag {
			for i, tag := range tags {
				if !strings.Contains(key, tag) {
					break
				}
				if i == len(tags)-1 {
					ResTags = append(ResTags, key)
				}
			}
		}
	default:
		S(w, 400)
		return
	}

	res, _ = json.Marshal(
		struct {
			Status int      `json:"status"`
			Type   string   `json:"type"`
			Tags   []string `json:"tags"`
		}{
			Status: 200,
			Type:   "tags",
			Tags:   ResTags,
		},
	)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return

}

//SearchVideos search videos
func SearchVideos(w http.ResponseWriter, r *http.Request) {
	var S func(w http.ResponseWriter, status int) = func(w http.ResponseWriter, status int) {
		res, _ := json.Marshal(Status{Status: status})
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		return
	}

	keys := lib.SplitTags(r.URL.Query().Get("Keywords"))
	mode := r.URL.Query().Get("Mode")

	AllVideos, err := lib.ReadVideoData()
	if err != nil {
		lib.Logger(err)
		slack.SendError(err)

		S(w, 500)
	}

	var ResVideos []lib.Video
	switch mode {
	case "or":
		for _, v := range AllVideos {
			for _, word := range keys {
				if strings.Contains(v.Title, word) {
					ResVideos = append(ResVideos, v)
					break
				}
			}
		}

	case "and":
		for i, v := range AllVideos {
			for _, word := range keys {
				if !strings.Contains(v.Title, word) {
					break
				}
				if i == len(keys)-1 {
					ResVideos = append(ResVideos, v)
				}
			}
		}
	default:
		S(w, 400)
		return
	}

	res, _ := json.Marshal(
		struct {
			Status int            `json:"status"`
			Type   string         `json:"type"`
			Video  []lib.ResVideo `json:"videos"`
		}{
			Status: 200,
			Type:   "videos",
			Video:  lib.VideosToResVideos(ResVideos),
		},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return

}
