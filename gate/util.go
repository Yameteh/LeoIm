package main

import (
	"crypto/md5"
	"encoding/json"
)

func GetMd5(in string) string {
 	h := md5.New()
	h.Write([]byte(in))
	return string(h.Sum(nil))
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