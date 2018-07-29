package main

import (
	"encoding/json"
)

const (
	AUTHACL_CODE_OK     = 200
	AUTHACL_CODE_ERROR  = 402
	AUTHACL_CODE_REAUTH = 401
)

type AuthRequest struct {
	User     string
	Response string
}

type AuthResponse struct {
	Code   int
	Nonce  string
	Method string
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
