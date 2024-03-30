package WebsocketHandler

import (
	"fmt"
	Game "monopoly_bank_go/game"
	HTTP "monopoly_bank_go/http"
	Parse "monopoly_bank_go/parse"
	Server "monopoly_bank_go/server"
)

func AcceptConnection(c *Server.Connection) {
	done := make(chan bool)

	go HTTP.HandlerRequest(c.Socket, func(r *HTTP.Request, err error) {
		if err != nil {
			c.SendAndClose(HTTP.MakeResponse(HTTP.BadRequest, nil, ""))
			return
		}

		if acceptKey, ok := r.Headers["Sec-WebSocket-Key"]; ok {
			if r.Headers["Upgrade"] == "websocket" {
				headers := map[string]string{
					"Upgrade":              "websocket",
					"Connection":           "Upgrade",
					"Sec-WebSocket-Accept": makeHandshakeResponseKey(acceptKey),
				}

				c.Send(HTTP.MakeResponse(HTTP.SwitchingProtocols, headers, ""))
				go listen(c)

				// release goroutine
				done <- true
				return
			}
		}

		c.SendAndClose(HTTP.MakeResponse(HTTP.BadRequest, nil, ""))
	})

	<-done
}

func listen(c *Server.Connection) {
	for {
		frame, err := Parse.ReadFromWebsocket(c.Socket)

		if err != nil {
			fmt.Println("error on reading from socket", err.Message)
			wErr := Parse.WriteToWebsocket(c.Socket, "", Parse.CloseFrame)

			if wErr != nil {
				fmt.Println("error on writing to socket", wErr)
			}

			c.Close()
			break
		}

		if frame.Op == Parse.CloseFrame {
			c.Close()
			break
		} else if frame.Op == Parse.PingFrame {
			wErr := Parse.WriteToWebsocket(c.Socket, string(frame.Payload), Parse.PongFrame)

			if wErr != nil {
				fmt.Println("error on writing to socket", wErr)
			}

			continue
		}

		// handle command
		go Game.AcceptCommand(string(frame.Payload), c)
	}
}
