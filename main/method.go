package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"../lib"
)

//Get get video data from allVideos.json
func (i *Index) Get() error {
	var V []lib.Video

	bData, err := ioutil.ReadFile(AllVideosFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bData, &V)

	i.Video = V

	return nil

}

//Encode returns get Query string from put map
func (query GetQuery) Encode() string {
	var result string = "?"
	var i int = 0
	for key, q := range query {
		if i == 0 {
			result += url.QueryEscape(key) + "=" + url.QueryEscape(q)
			i++
			continue
		}
		result += "&" + url.QueryEscape(key) + "=" + url.QueryEscape(q)
	}
	return result
}
