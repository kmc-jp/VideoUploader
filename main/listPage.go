package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"../lib"
)

// ListPageHandle  handle user list page
func ListPageHandle(w http.ResponseWriter, r *http.Request) {
	// ユーザデータの取得
	var user lib.User
	if e := user.Get(os.Getenv("REMOTE_USER")); e != nil {
		user = lib.User{
			Name:  os.Getenv("REMOTE_USER"),
			Video: []lib.Video{},
		}
	}
	var U UserListPage = UserListPage{
		Header: Header{
			User:     user,
			Title:    "RoomList",
			UserName: os.Getenv("REMOTE_USER"),
			Error:    ErrorHandle(r.URL.Query().Get("Error")),
			Success:  SuccessHandle(r.URL.Query().Get("Success")),
		},
		Footer: Footer{},
	}

	list, e := UserListMake()
	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	U.User = list

	t, e := template.New("").ParseFiles(
		filepath.Join("resources", "header.html"),
		filepath.Join("resources", "list.html"),
		filepath.Join("resources", "style.html"),
		filepath.Join("resources", "footer.html"),
		filepath.Join("resources", "script.html"),
	)

	if e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	w.Header().Add("Content-Type", "text/html")

	if e = t.ExecuteTemplate(w, "list", U); e != nil {
		fmt.Fprintf(w, "%s", e.Error())
		return
	}

	return
}

// UserListMake make user list, which is sorted by its user names
func UserListMake() ([]lib.User, error) {
	data, err := lib.ReadUserData()
	if err != nil {
		return []lib.User{}, err
	}

	var list []lib.User

	for _, u := range data {
		list = append(list, u)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })

	return list, nil
}
