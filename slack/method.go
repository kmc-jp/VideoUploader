package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"../lib"
)

// Send message to slack
func (w Webhook) Send() error {
	if lib.Settings.SlackWebhook == "" {
		return fmt.Errorf("Webhook not setuped")
	}
	bData, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		return err
	}

	_, err = http.Post(lib.Settings.SlackWebhook, "application/json", bytes.NewBuffer(bData))
	return err
}
