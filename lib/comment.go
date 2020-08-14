package lib

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

//CommentPath return comment file path
func CommentPath(VideoID string) string {
	return filepath.Join("Comment", VideoID+".json")
}

//MakeCommentID generates random id
func MakeCommentID(comments []Comment) string {
again:
	var id = MakeRandStr(10)

	for _, c := range comments {
		if c.ID == id {
			goto again
		}
	}

	return id
}

//CommentFind find comment from slice
func CommentFind(comments []Comment, id string) (comment Comment, num int) {
	for i, c := range comments {
		if c.ID == id {
			return c, i
		}
	}
	return
}

// CommentDelete delete commment with the id
func CommentDelete(comments []Comment, id string) error {
	_, num := CommentFind(comments, id)

	if num == len(comments)-1 {
		comments = comments[:num]
	} else {
		comments = append(comments[:num], comments[num+1:]...)
	}

	return nil
}

// CommentSave save comments
func CommentSave(comments []Comment, VideoID string) error {
	var CommentPath = CommentPath(VideoID)

	bData, err := json.MarshalIndent(comments, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(CommentPath, bData, 0777)

}
