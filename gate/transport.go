package main

type Transport interface {
	Listen(domain string, port int)

}
