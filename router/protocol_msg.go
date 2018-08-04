package main


type MessageBody struct {
	MsgType  int
	From string
	To string
	Time int64
	MimeType string
	Content string
}

