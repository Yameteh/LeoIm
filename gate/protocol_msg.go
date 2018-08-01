package main



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


