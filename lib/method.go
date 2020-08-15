package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Get user data
func (u *User) Get(name string) error {
	var data, err = ReadUserData()
	if err != nil {
		return err
	}

	_, ok := data[name]
	if !ok {
		data[name] = User{
			Name:  name,
			Video: []Video{},
		}
	}

	*u = data[name]

	return nil
}

//WriteAll exports user data into allVideos.json and userinfo.json
func (u *User) WriteAll() error {
	var V []Video
	bData, err := ioutil.ReadFile(AllVideosFile)
	err = json.Unmarshal(bData, &V)
	if err != nil {
		return err
	}

	if cap(V) == 0 {
		V = append(V, u.Video[len(u.Video)-1])
	} else {
		V = append([]Video{u.Video[0]}, V[0:]...)
	}

	bData, err = json.MarshalIndent(V, "", "    ")
	err = ioutil.WriteFile(AllVideosFile, bData, 0777)
	if err != nil {
		return err
	}

	var data map[string]User

	bData, err = ioutil.ReadFile(UserInfoFile)
	if err != nil {
		return fmt.Errorf("NotFound")
	}

	json.Unmarshal(bData, &data)

	data[u.Name] = *u

	bData, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(UserInfoFile, bData, 0777)
}

//Write exports user data into userinfo.json
func (u *User) Write() error {
	var data map[string]User

	bData, err := ioutil.ReadFile(UserInfoFile)
	if err != nil {
		return fmt.Errorf("NotFound")
	}

	json.Unmarshal(bData, &data)

	data[u.Name] = *u

	bData, err = json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(UserInfoFile, bData, 0777)
}

//Update updates index data with put video data
func (v Video) Update() (err error) {
	var U func(Videos []Video) error = func(Videos []Video) error {
		var video = SearchVideo(Videos, v.Video)
		if video.Video != v.Video {
			return fmt.Errorf("Not Found")
		}

		switch {
		case v.Title != "":
			video.Title = v.Title
			fallthrough
		case len(v.Tags) > 0:
			video.Tags = v.Tags
			fallthrough
		case v.Status.Phase != "":
			video.Status.Phase = v.Status.Phase
		case v.Status.Error != "":
			video.Status.Error = v.Status.Error
		}

		UpdateVideo(Videos, video)
		return nil
	}

	var data map[string]User

	bData, err := ioutil.ReadFile(UserInfoFile)
	if err != nil {
		return fmt.Errorf("NotFound")
	}

	err = json.Unmarshal(bData, &data)
	if err != nil {
		return err
	}

	if err := U(data[v.User].Video); err != nil {
		return err
	}

	bData, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(UserInfoFile, bData, 0777)

	bData, err = ioutil.ReadFile(AllVideosFile)
	if err != nil {
		return fmt.Errorf("NotFound")
	}

	var V []Video

	err = json.Unmarshal(bData, &V)
	if err != nil {
		return err
	}

	err = U(V)
	if err != nil {
		return err
	}

	bData, err = json.MarshalIndent(V, "", "    ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(AllVideosFile, bData, 0777)

	return nil
}

//TagAppend append video to tag list
func (v Video) TagAppend() (err error) {
	var AllTag Tag
	err = AllTag.Get()
	if err != nil {
		return err
	}

	AllTag.Append(v)

	err = AllTag.Save()
	if err != nil {
		return err
	}
	return nil
}

//Update update userinfo.json
func (u User) Update() error {

	var data map[string]User

	bData, err := ioutil.ReadFile(UserInfoFile)
	if err != nil {
		return fmt.Errorf("NotFound")
	}

	err = json.Unmarshal(bData, &data)
	if err != nil {
		return err
	}

	data[u.Name] = u

	bData, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(UserInfoFile, bData, 0777)

	return nil
}

//Save save comment data
func (c Comment) Save(VideoID string) error {
	var CommentPath string = CommentPath(VideoID)
	var comments []Comment

	if !FileExistance(CommentPath) {
		ioutil.WriteFile(CommentPath, []byte("[]"), 0777)
	}

	bData, err := ioutil.ReadFile(CommentPath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bData, &comments); err != nil {
		return err
	}

	comments = append([]Comment{c}, comments...)

	return CommentSave(comments, VideoID)

}

//Get read tag.json
func (tag *Tag) Get() (err error) {
	var data Tag
	if !FileExistance(TagFile) {
		ioutil.WriteFile(TagFile, []byte("{}"), 0777)
	}

	bData, err := ioutil.ReadFile(TagFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(bData, &data)

	tag = &data
	if tag.Tag == nil {
		tag.Tag = make(map[string][]string)
	}
	return
}

//Save save tag data
func (tag Tag) Save() (err error) {
	bData, err := json.MarshalIndent(tag, "", "    ")
	if err != nil {
		return
	}

	return ioutil.WriteFile(TagFile, bData, 0777)
}

//Append append tag
func (tag *Tag) Append(video Video) {
	if tag.Tag == nil {
		tag.Tag = make(map[string][]string)
	}
	for _, key := range video.Tags {
		tag.Tag[key] = append(tag.Tag[key], video.Video)
	}
}

//Remove remove one video from tag list
func (tag *Tag) Remove(video Video) (err error) {
	if tag.Tag == nil {
		return
	}
	for _, t := range video.Tags {
		var vs []string
		for _, v := range tag.Tag[t] {
			if v != video.Video {
				vs = append(vs, video.Video)
			}
		}

		tag.Tag[t] = vs
	}

	return
}
