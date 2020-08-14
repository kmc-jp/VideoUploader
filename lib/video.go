package lib

import "fmt"

//UpdateVideo replace the video data in the list with put video data
func UpdateVideo(Videos []Video, newVideo Video) error {
	for i, v := range Videos {
		if v.Video == newVideo.Video {
			Videos[i] = newVideo
			return nil
		}
	}
	return fmt.Errorf("Not Found")
}

//DeleteVideo delete video from list
func DeleteVideo(Videos []Video, newVideo Video) ([]Video, error) {
	for i, v := range Videos {
		if v.Video == newVideo.Video {
			if len(Videos) < i+1 {
				return Videos[:i], nil
			}
			return append(Videos[:i], Videos[i+1:]...), nil
		}
	}
	return []Video{}, fmt.Errorf("Not Found")
}

//SearchVideo returns put video from []Video
func SearchVideo(Videos []Video, name string) Video {
	for _, v := range Videos {
		if v.Video == name {
			return v
		}
	}
	return Video{}
}
