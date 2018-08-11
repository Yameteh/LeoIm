package main


type MessageBody struct {
	MsgType  int
	From string
	To string
	Time int64
	MimeType string
	Content string
}

type SyncResponse struct {
	Server string
	Time  int64
}

type StreamSdp struct {
	InAddr string
	AudioCodec string
	AudioPort int
	AudioSampleRate int
	AudioBitrate int
	VideoCodec string
	VideoPort int
	VideoWidth int
	VideoHeight int
	VideoFrameRate int
	VideoBitrate int
}
