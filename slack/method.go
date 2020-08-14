package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"../lib"
)

// Send message to slack
func (w Webhook) Send() error {
	bData, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		return err
	}

	_, err = http.Post(lib.Settings.SlackWebhook, "application/json", bytes.NewBuffer(bData))
	return err
}
