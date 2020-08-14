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
    "ffprobe": "ffproveの相対パス"
}
```

あとはビルドする。