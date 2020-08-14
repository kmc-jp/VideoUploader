package slack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"../lib"
)

//SendVideoInfo sends video information to slack channel
func SendVideoInfo(video lib.Video, channel, URL string) error {
	link, err := func() (string, error) {
		b, err := ioutil.ReadFile(filepath.Join("Videos", video.Thumb))
		if err != nil {
			return "", err
		}

		return lib.GyazoFileUpload(b)
	}()

	if err != nil {
		return err
	}

	var Message Webhook = Webhook{
		Channel:  channel,
		UserName: "VideoUploader",

		Attachments: []Attachment{
			{
				Title:    video.Title,
				Text:     fmt.Sprintf("%sの神動画だよ！見てね！\nURI:%s", video.User, URL),
				ImageURL: link,

				AuthorName: "VideoUploder",
			},
		},
	}

	return Message.Send()

}

//SendComment Send sent comment
func SendComment(comment lib.Comment, video lib.Video, channel, URL string) error {
	link, err := func() (string, error) {
		b, err := ioutil.ReadFile(filepath.Join(comment.Icon))
		if err != nil {
			return "", err
		}

		return lib.GyazoFileUpload(b)
	}()
	if err != nil {
		return err
	}

	var Message Webhook = Webhook{
		Channel:  channel,
		UserName: "VideoUploader",
		Text:     fmt.Sprintf("あなたの動画にコメントが届いたよ！\n動画名:%s", video.Title),
		IconURL:  link,

		Attachments: []Attachment{
			{
				Text:       fmt.Sprintf(" %s", comment.Comment),
				AuthorName: comment.User,
			},
			{
				Text: fmt.Sprintf("動画URI:%s", URL),
			},
		},
	}

	return Message.Send()
}

// SendError send error to the Slack error channel
func SendError(e error) {
	var Message = Webhook{
		Text:     fmt.Sprintf("問題が発生しました。"),
		Channel:  lib.Settings.ErrorChannel,
		UserName: "VideoUploader -ErrorNotifier- ",
		Attachments: []Attachment{
			{
				Title: "問題の内容",
				Text:  fmt.Sprintf("%s", e.Error()),
			},
		},
	}

	Message.Send()
}
