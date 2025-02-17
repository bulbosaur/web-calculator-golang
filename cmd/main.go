package main

import (
	config "github.com/bulbosaur/web-calculator-golang/config"
	server "github.com/bulbosaur/web-calculator-golang/internal/http"
)

func main() {
	config.Init()
	server.RunServer()
}
