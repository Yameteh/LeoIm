package main

type AuthRequest struct {
	User string
	Response string
}

type AuthResponse struct {
	Code int
	Nonce string
	Method string
}