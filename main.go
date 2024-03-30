package main

import (
	"monopoly_bank_go/file-handler"
	"monopoly_bank_go/server"
	WebsocketHandler "monopoly_bank_go/websocket-handler"
	"time"
)

func main() {
	Server.NewTCPServer(Server.FILE, "80", FileHandler.AcceptConnection)
	Server.NewTCPServer(Server.WEBSOCKET, "4444", WebsocketHandler.AcceptConnection)

	for {
		time.Sleep(10 * time.Second)

		if !Server.IsOpen {
			break
		}
	}
}
