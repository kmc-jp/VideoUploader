package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	//Extension define save video's format
	Extension = ".mp4"
	// UserInfoFile put path to userinfo.json
	UserInfoFile string
	// AllVideosFile put path to allVideos.json
	AllVideosFile string
	//TagFile put path to tags.json
	TagFile string

	//Settings put setting data
	Settings Setting
)

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	bytes, err := ioutil.ReadFile(filepath.Join("Settings.json"))
	if err != nil {
		panic("Panic:Settings.json not found!")
	}
	json.Unmarshal(bytes, &Settings)

	UserInfoFile = filepath.Join("User", "userinfo.json")
	AllVideosFile = filepath.Join("User", "allVideos.json")
	TagFile = filepath.Join("tags.json")
}

//MakeRandStr returns ramdom string (length n)
func MakeRandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

//FindExtension returns filename and extensions
func FindExtension(filename string) (name, extension string) {
	if !strings.Contains(filename, ".") {
		name = filename
		return name, extension
	}

	var sep []string = strings.Split(filename, ".")
	name = strings.Join(sep[:len(sep)-1], ".")
	extension = sep[len(sep)-1]

	extension = strings.ReplaceAll(strings.ReplaceAll(extension, "/", "_"), "\\", "_")

	return
}

//FileExistance returns whether file exists or not
func FileExistance(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return true
}

//MakeSuitName make new file name
func MakeSuitName() string {
	var name = MakeRandStr(20)

	if FileExistance(filepath.Join("Videos", name+Extension)) {
		return MakeSuitName()
	}
	return name
}

//TmpClear delete the temporary directory
func TmpClear(videoName string) error {
	err := os.RemoveAll(filepath.Join("tmp", videoName))
	if err != nil {
		return err
	}
	return nil
}

//Logger exports put err to log.txt
func Logger(err error) error {
	if !FileExistance("log.txt") {
		ioutil.WriteFile("log.txt", []byte{}, 0777)
	}
	bData, e := ioutil.ReadFile("log.txt")
	if e != nil {
		return err
	}

	bData = bytes.Join([][]byte{[]byte(fmt.Sprintf("%s: %s\n", time.Now().Format(time.Stamp), err.Error())), bData}, []byte(""))

	e = ioutil.WriteFile("log.txt", bData, 0777)

	return e
}

//Progress export uploading status
func Progress(video Video, phase Status) {
	video.Status = phase

	video.Update()

	return
}

// VideoExistance Check video existance
func VideoExistance(VideoID string) bool {
	if VideoID == "" {
		return false
	}
	videos, _ := ReadVideoData()
	video := SearchVideo(videos, VideoID)
	return video.Video == VideoID
}
