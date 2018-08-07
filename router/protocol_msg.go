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

