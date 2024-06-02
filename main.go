package main

import (
	"monopoly_bank_go/game"
	"monopoly_bank_go/server"
	"monopoly_bank_go/static_files"
	"monopoly_bank_go/types"
	"monopoly_bank_go/webapi"
	"monopoly_bank_go/websocket"
	"time"
)

func main() {
	server.NewTCPServer(types.FILE, 80, static_files.Handler)

	websocket.MessageHandler = game.AcceptCommand
	server.NewTCPServer(types.WEBSOCKET, 4444, websocket.Handler)

	server.NewTCPServer(types.WEBAPI, 7600, webapi.Handler)

	for {
		time.Sleep(10 * time.Second)
		go server.DumpClosedConnections()

		if !server.IsOpen {
			break
		}
	}
}
