package main

import (
	"net/http"
	"net/http/cgi"

	"../lib"
)

const (

	//ThumSizeMax define user thumbnail size(px)
	ThumSizeMax = 100

	//IconSize set icon size
	IconSize = 100
)

//Settings put setting data
var Settings lib.Setting
var (
	//Extension define save video's format
	Extension = lib.Extension
	// UserInfoFile put path to userinfo.json
	UserInfoFile string
	// AllVideosFile put path to allVideos.json
	AllVideosFile string
)

func init() {

	UserInfoFile = lib.UserInfoFile
	AllVideosFile = lib.AllVideosFile

	Extension = lib.Extension
	Settings = lib.Settings
}

func main() {
	cgi.Serve(
		http.HandlerFunc(Handler),
	)
}

//Handler handles http requests and excute a suitable function
func Handler(w http.ResponseWriter, r *http.Request) {
	//Handle Action first
	switch r.URL.Query().Get("Action") {
	case "Upload":
		UploadHandle(w, r)
		return
	case "Update":
		UpdateVideoInfo(w, r)
		return
	case "Delete":
		Delete(w, r)
		return
	case "Seticon":
		SetUserIcon(w, r)
		return
	case "SendSlack":
		SendSlack(w, r)
		return
	case "GetComment":
		ReadComment(w, r)
		return
	case "AddComment":
		AddComment(w, r)
		return
	case "Setchannel":
		SetSlackChannel(w, r)
		return
	case "VideoStatus":
		SendVideoStatus(w, r)
		return
	case "SearchTags":
	}

	//Handle showing page second
	switch r.URL.Query().Get("Page") {
	case "Upload":
		UploadPageHandle(w, r)
		return
	case "Play":
		PlayPageHandle(w, r)
		return
	case "MyPage":
		MyPageHandle(w, r)
		return
	case "Settings":
		SettingsPageHandle(w, r)
	case "Users":
		ListPageHandle(w, r)
	case "UserPage":
		UserPageHandle(w, r)
	default:
		IndexPageHandle(w, r)
		return
	}
}

//SuccessHandle returns success message from key
func SuccessHandle(suc string) string {
	switch suc {
	case "upload":
		return "アップロードに成功しました。エンコードの状態はMyRoomで確認できます。"
	case "delete":
		return "動画を削除しました。"
	case "update":
		return "保存しました。"
	case "SetIcon":
		return "設定しました。"
	case "SendSlack":
		return "Slackに送信しました。"
	case "SetChannel":
		return "送信先の設定をしました。"
	}
	return ""
}

//ErrorHandle returns success message from key
func ErrorHandle(err string) string {
	switch err {
	case "UpdateError":
		return "動画情報の更新に失敗しました。各種権限設定などに問題がある可能性があります。"
	case "EncodeError":
		return "動画のエンコードに失敗しました。形式が間違っている可能性があります。"
	case "NotFound":
		return "動画が見つかりませんでした。再度操作をやりなおしてください。"
	case "EncodeTimeError":
		return "動画の長さの解析に失敗しました。形式が間違っているか、極端に長さが短い可能性があります。"
	case "UserDataUpdateError":
		return "ユーザ情報の更新に失敗しました。何度も同じエラーが発生する場合は管理者に連絡してください。"
	case "IconNotFound":
		return "アイコンデータのアップロードに失敗しました。"
	case "ResizeError":
		return "画像の拡縮に失敗しました。"
	case "SendSlack":
		return "Slackへの送信に失敗しました。改善しない場合は管理者に連絡してください。"
	case "SetChannel":
		return "Slackチャンネル設定の書き込みに失敗しました。しばらくしても改善しない場合は管理者に連絡してください。"
	case "TagRemoveError":
		return "タグ情報の削除に失敗しました。"
	}

	return ""
}
