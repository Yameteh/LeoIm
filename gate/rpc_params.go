package main

type User struct {
	Account  string
	Password string
}

type Message struct {
	Version uint8
	Type    uint8
	Body    string
}

