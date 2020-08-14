package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

//ReadVideoData returns the data of allVideos.json
func ReadVideoData() ([]Video, error) {
	var Videos []Video

	bytes, err := ioutil.ReadFile(filepath.Join("User", "allVideos.json"))
	if err != nil {
		return Videos, err
	}

	err = json.Unmarshal(bytes, &Videos)
	if err != nil {
		return Videos, err
	}
	return Videos, nil
}

//ReadUserData read userinfo.json
func ReadUserData() (map[string]User, error) {
	var data map[string]User

	bData, err := ioutil.ReadFile(UserInfoFile)
	if err != nil {
		return data, fmt.Errorf("NotFound")
	}

	json.Unmarshal(bData, &data)

	return data, nil
}

// ReadComments read comment data
func ReadComments(VideoID string) (c []Comment, err error) {
	var CommentPath string = CommentPath(VideoID)

	if !FileExistance(CommentPath) {
		ioutil.WriteFile(CommentPath, []byte("[]"), 0777)
	}

	bData, err := ioutil.ReadFile(CommentPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(bData, &c)
	if err != nil {
		return
	}

	return
}
