# VideoUploader
## 概要
CGIの動画投稿サービス作成中。

## 用意

ルートディレクトリに
`Settings.json`
を用意し、次のようにする。

```json
{
    "ffmpeg": "ffmpegの相対パス",
    "ffprobe": "ffproveの相対パス",
    "gyazo_token": "GyazoのToken",
    "slack_webhook": "SlackのIncommingWebhookURL"
}
```

`main/`をビルドし、`index.up`という名前に変更する。

他、
`encoder/`もビルドし、常態化させておく。

以上で準備は終わり。