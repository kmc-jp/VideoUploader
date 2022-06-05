package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/floostack/transcoder/ffmpeg"
	"github.com/kmc-jp/VideoUploader/encoder/lib"
	"github.com/kmc-jp/VideoUploader/encoder/slack"
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
		"-codec:a", "libfdk_aac",
		"-b:a", "320k",
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
		err = makeVideoThumb(newData.Video)
		if err != nil {
			lib.Logger(fmt.Errorf("Encode:%s\n", err.Error()))
			slack.SendError(fmt.Errorf("Encode:%s\n", err.Error()))
			lib.Progress(newData, lib.Status{Error: fmt.Sprintf("Encode:\n%s\n", err.Error())})
			os.Remove(filepath.Join("Videos", newData.Video))
			lib.TmpClear(newData.Video)
			return nil
		}

		os.Chmod(filepath.Join("Videos", newData.Video+".png"), 0777)
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

// makeVideoThumb Make video thumbnail
func makeVideoThumb(Video string) error {
	const ThumbSizeMax = 360
	var ff = ffmpeg.
		New(
			&ffmpeg.Config{
				FfmpegBinPath:   lib.Settings.FFmpeg,
				FfprobeBinPath:  lib.Settings.FFprobe,
				ProgressEnabled: false,
			},
		)

	meta, err := ff.Input(filepath.Join("Videos", Video)).GetMetadata()
	if err != nil {
		return err
	}

	du, err := strconv.ParseFloat(meta.GetFormat().GetDuration(), 64)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	var width, height int

	var scale string
	// var format = "jpg"
	var overwrite = true
	var ss = strconv.Itoa(rand.Intn(int(du)))
	var vframes = 1

	for _, st := range meta.GetStreams() {
		if st.GetWidth() == 0 || st.GetHeight() == 0 {
			continue
		}
		width = st.GetWidth()
		height = st.GetHeight()
	}

	switch {
	case width == 0 || height == 0:
		return fmt.Errorf("VideoSizeError")
	case width > height:
		var nHeight = int(float32(height*ThumbSizeMax) / float32(width))
		scale = fmt.Sprintf("%dx%d", ThumbSizeMax, nHeight)
	case width <= height:
		var nWidth = int(float32(width*ThumbSizeMax) / float32(height))
		scale = fmt.Sprintf("%dx%d", nWidth, ThumbSizeMax)
	}

	var opts = ffmpeg.Options{
		// OutputFormat: &format,
		Overwrite:  &overwrite,
		SeekTime:   &ss,
		Vframes:    &vframes,
		Resolution: &scale,
	}

	_, err = ff.
		Input(filepath.Join("Videos", Video)).
		Output(filepath.Join("Videos", Video+".png")).
		WithOptions(opts).
		Start(opts)

	return err
}
