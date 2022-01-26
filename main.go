package main

import (
	"github.com/crazy-me/apps-produce/server"
)

func main() {
	ws := server.New("0.0.0.0", 8888)
	ws.Start()
}
