package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"../lib"
	"../slack"
)

func main() {

	for {
		if err := WatchTmpDir(); err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Millisecond * 500)
	}

}

// WatchTmpDir watch temporary directory and encode videos when video found
func WatchTmpDir() (err error) {
	files, err := ioutil.ReadDir("tmp")
	if err != nil {
		return
	}

	AllVideos, err := lib.ReadVideoData()
	if err != nil {
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		_, extension := lib.FindExtension(f.Name())
		if "."+extension != lib.Extension {
			continue
		}

		video := lib.SearchVideo(AllVideos, f.Name())
		if video.Video == "" {
			continue
		}

		if video.Status.Phase != "calling encode process..." {
			continue
		}

		Encode(video)

	}
	return

}

// Encode encodes put videos
func Encode(newData lib.Video) (err error) {

	lib.Progress(newData, lib.Status{Phase: "search for the saved video"})

	var videoName string

	{
		files, err := ioutil.ReadDir(filepath.Join("tmp", newData.Video))
		if err != nil {
			lib.Logger(fmt.Errorf("Encode:%s", err.Error()))
			slack.SendError(fmt.Errorf("Encode:%s", err.Error()))
			lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})
			return nil
		}
		for _, f := range files {
			name, _ := lib.FindExtension(f.Name())
			if name == "video" {
				videoName = f.Name()
				break
			}
		}
		if videoName == "" {
			return fmt.Errorf("Video Not Found")
		}

	}

	lib.Progress(newData, lib.Status{Phase: "encoding"})

	out, err := exec.Command(
		lib.Settings.FFmpeg,
		"-i", filepath.Join("tmp", newData.Video, videoName),
		filepath.Join("Videos", newData.Video),
	).CombinedOutput()

	if err != nil {
		lib.Logger(fmt.Errorf("Encode:%s\nFFmpegState:%s", err.Error(), out))
		slack.SendError(fmt.Errorf("Encode:%s\nFFmpegState:%s", err.Error(), out))
		lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})
		lib.TmpClear(newData.Video)
		return nil
	}

	lib.Progress(newData, lib.Status{Phase: "generating thumbnail"})

	if !lib.FileExistance(filepath.Join("Videos", newData.Thumb)) {
		out, err = exec.Command(lib.Settings.FFprobe, filepath.Join("Videos", newData.Video)).CombinedOutput()
		if err != nil {
			lib.Logger(fmt.Errorf("Encode:%s\nFFproveState:%s", err.Error(), out))
			slack.SendError(fmt.Errorf("Encode:%s\nFFproveState:%s", err.Error(), out))
			lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})

			os.Remove(filepath.Join("Videos", newData.Video))
			lib.TmpClear(newData.Video)

			return nil
		}
		t, err := getffTimeStr(out)
		if err != nil {
			lib.Logger(fmt.Errorf("Encode:%s\nFFproveState:%s", err.Error(), out))
			slack.SendError(fmt.Errorf("Encode:%s\nFFproveState:%s", err.Error(), out))
			lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})
			os.Remove(filepath.Join("Videos", newData.Video))
			lib.TmpClear(newData.Video)
			return nil
		}
		if t[2] == "00" && t[1] == "00" {
			lib.Logger(fmt.Errorf("Encode:%s\nFFproveState:%s", err.Error(), out))
			slack.SendError(fmt.Errorf("Encode:%s\nFFproveState:%s\ntime%#v", err.Error(), out, t))
			lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})
			lib.TmpClear(newData.Video)

			return nil
		}

		out, err = exec.Command(
			lib.Settings.FFmpeg, "-i", filepath.Join("tmp", newData.Video, videoName),
			"-ss", t[2],
			"-vframes", "1",
			"-filter:v",
			"scale=360:-1",
			filepath.Join("Videos", newData.Video+".png"),
		).CombinedOutput()

		if err != nil {
			lib.Logger(fmt.Errorf("Encode:%s\nFFmpegState:%s", err.Error(), out))
			slack.SendError(fmt.Errorf("Encode:%s\nFFmpegState:%s", err.Error(), out))
			lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})
			os.Remove(filepath.Join("Videos", newData.Video))
			lib.TmpClear(newData.Video)
			return nil
		}

	}

	lib.TmpClear(newData.Video)
	lib.Progress(newData, lib.Status{})
	return nil
}

func getffTimeStr(out []byte) ([]string, error) {
	var t = strings.Split(string(out), "Duration: ")
	if len(t) == 1 {
		return []string{}, fmt.Errorf("TypeError")
	}
	return strings.Split(strings.Split(t[1], ".")[0], ":"), nil
}
