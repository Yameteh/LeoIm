package main

import (
	"crypto/md5"
	"encoding/json"
	"encoding/hex"
	"encoding/base64"
)

func GetMd5(in string) string {
 	h := md5.New()
	h.Write([]byte(in))
	return hex.EncodeToString(h.Sum(nil))
}

func GetBase64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func CreateProtocolMsg(v uint8, t uint8, msg interface{}) *Protocol {
	p := new(Protocol)
	p.Version = v
	p.Type = t

	if j, err := json.Marshal(msg); err == nil {
		p.Body = string(j)
		p.Length = uint32(len(p.Body))
		return p
	} else {
		return nil
	}
}