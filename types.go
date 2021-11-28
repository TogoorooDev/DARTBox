package main

type confServer struct {
	Port string
}

type config_format struct {
	Server confServer
}

type confirmFormat struct {
	Ip       string
	Filename string
}
