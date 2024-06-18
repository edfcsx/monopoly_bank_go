package main

import (
	"monopoly_bank_go/accounts"
	"monopoly_bank_go/commands"
	"monopoly_bank_go/connection"
	"monopoly_bank_go/http"
	"monopoly_bank_go/server"
	"monopoly_bank_go/static_files"
	"monopoly_bank_go/types"
	"monopoly_bank_go/webapi"
	"monopoly_bank_go/websocket"
	"time"
)

func main() {
	server.NewTCPServer(types.FILE, 80, static_files.Handler)
	server.NewTCPServer(types.WEBAPI, 7600, webapi.Handler)

	websocket.MessageHandler = commands.Handler

	// is created a go routine "done" to free function after call websocket.listen
	server.NewTCPServer(types.WEBSOCKET, 4444, func(c *connection.Connection) {
		done := make(chan bool)

		go http.HandlerRequest(c.Socket, func(r *http.Request, err error) {
			if err != nil {
				c.SendAndClose(http.MakeResponse(http.BadRequest, nil, ""))
				return
			}

			if acceptKey, ok := r.Headers["Sec-WebSocket-Key"]; ok {
				if r.Headers["Upgrade"] == "websocket" {
					// Authorization check
					if hash, ok := r.Query["player_hash"]; ok {
						account := accounts.ExistsByHash(hash)

						if account == nil {
							c.SendAndClose(http.MakeResponse(http.Unauthorized, nil, ""))
							return
						}
					}

					headers := map[string]string{
						"Upgrade":              "websocket",
						"Connection":           "Upgrade",
						"Sec-WebSocket-Accept": websocket.MakeHandshakeKey(acceptKey),
					}

					c.Send(http.MakeResponse(http.SwitchingProtocols, headers, ""))
					go websocket.Listen(c)

					// release goroutine
					done <- true
					return
				}
			}

			c.SendAndClose(http.MakeResponse(http.BadRequest, nil, ""))
		})
	})

	for {
		time.Sleep(10 * time.Second)
		go server.DumpClosedConnections()

		if !server.IsOpen {
			break
		}
	}
}
