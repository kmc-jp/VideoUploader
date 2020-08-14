package lib

import (
	"bytes"
	"fmt"

	"github.com/tomohiro/go-gyazo/gyazo"
)

//GyazoFileUpload Upload image to Gyazo and return PermalinkURL
func GyazoFileUpload(dataBytes []byte) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			Logger(fmt.Errorf("Gyazo:%v", err))
			return
		}
	}()

	gyazo, err := gyazo.NewClient(Settings.GyazoToken)
	if err != nil {
		return "", err
	}

	image, err := gyazo.Upload(bytes.NewReader(dataBytes))

	return image.URL, err
}
