package lib

import (
	"sort"
	"strings"
)

//SplitTags divide str into tags
func SplitTags(str string) []string {
	var strs = strings.Split(str, ",")
	var res []string = make([]string, 0)
	for _, s := range strs {
		res = append(res, TrimSpaces(strings.Split(s, " "))...)
	}

	return res
}

// TrimSpaces trim spaces in []string
func TrimSpaces(strs []string) []string {
	var tags = make(map[string]bool)
	var res []string = make([]string, 0)
	for i := 0; i < len(strs); i++ {
		strs[i] = strings.TrimSpace(strings.ReplaceAll(strs[i], "\n", ""))
		if _, ok := tags[strs[i]]; ok || strs[i] == "" {
			continue
		}

		res = append(res, strs[i])
		tags[strs[i]] = true
	}
	return res
}

//StoM converts []tags to map[tag]bool
func StoM(tags []string) map[string]bool {
	var Res = make(map[string]bool)

	for _, t := range tags {
		Res[t] = true
	}
	return Res
}

//MtoS converts map[tag]bool to []tag
func MtoS(tags map[string]bool) []string {
	var Res []string = make([]string, 0)

	for t, b := range tags {
		if b {
			Res = append(Res, t)
		}
	}

	sort.Slice(Res, func(i, j int) bool { return Res[i] < Res[j] })

	return Res
}

//TagUpdate update tag info
func TagUpdate(oldVideo Video, newTags []string) error {
	var AllTag Tag
	AllTag, err := AllTag.Get()
	if err != nil {
		return err
	}

	AllTag, err = AllTag.Remove(oldVideo)
	if err != nil {
		return err
	}

	oldVideo.Tags = newTags

	AllTag = AllTag.Append(oldVideo)

	return AllTag.Save()

}

//TagRemove remove put video from tag list
func TagRemove(oldVideo Video) error {
	var AllTag Tag
	AllTag, err := AllTag.Get()
	if err != nil {
		return err
	}

	AllTag, err = AllTag.Remove(oldVideo)
	if err != nil {
		return err
	}

	return AllTag.Save()
}

// VideoIDsToVideos return Video from VideoID
func VideoIDsToVideos(VideoIDs []string) (ResVideo []Video, err error) {
	AllVideos, err := ReadVideoData()
	if err != nil {
		return
	}

	var AllVideoMap = func(vs []Video) map[string]Video {
		var Res = make(map[string]Video)
		for _, v := range vs {
			Res[v.Video] = v
		}

		return Res
	}(AllVideos)

	for _, video := range VideoIDs {
		v, ok := AllVideoMap[video]
		if !ok {
			continue
		}
		ResVideo = append(ResVideo, v)
	}

	return ResVideo, nil
}
