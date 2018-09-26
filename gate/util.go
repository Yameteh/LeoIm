package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

func GetMd5(in string) string {
	h := md5.New()
	h.Write([]byte(in))
	return hex.EncodeToString(h.Sum(nil))
}

func GetBase64(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func CreateProtocolMsg(v uint8, t uint8, msg interface{}) *Protocol {
	p := new(Protocol)
	p.Version = v
	p.Type = t

	if j, err := json.Marshal(msg); err == nil {
		p.Body = j
		p.Length = uint32(len(j))
		return p
	} else {
		return nil
	}
}
