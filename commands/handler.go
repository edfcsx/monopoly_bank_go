package commands

import (
	"fmt"
	"monopoly_bank_go/connection"
	"monopoly_bank_go/mjtp"
	"monopoly_bank_go/server"
	"monopoly_bank_go/types"
	"monopoly_bank_go/websocket"
)

type Request struct {
	Data       *mjtp.Message
	Connection *connection.Connection
	response   []Response
}

type Response struct {
	Msg       *mjtp.Message
	SendToAll bool
}

func (r *Request) Respond(msg *mjtp.Message, sendToAll bool) {
	r.response = append(r.response, Response{msg, sendToAll})
}

func (r *Request) SendResponses() {
	for _, res := range r.response {
		if arg, ok := r.Data.Body["args_id"]; ok {
			if !res.SendToAll {
				res.Msg.Body["args_id"] = arg
			}
		}

		responseText, err := res.Msg.String()

		if err != nil {
			fmt.Println("error on parsing response data", err)
			return
		}

		var errWrite *websocket.FrameError

		if res.SendToAll {
			for _, conn := range server.GetConnections(types.WEBSOCKET) {
				errWrite = websocket.Write(conn.Socket, responseText, websocket.TextFrame)

				if errWrite != nil {
					fmt.Println("error on writing response", errWrite)
					conn.Close()
				}
			}
		} else {
			errWrite = websocket.Write(r.Connection.Socket, responseText, websocket.TextFrame)

			if errWrite != nil {
				fmt.Println("error on writing response", errWrite)
				r.Connection.Close()
			}
		}
	}
}

func sendRawResponse(data string, c *connection.Connection) {
	errWrite := websocket.Write(c.Socket, data, websocket.TextFrame)

	if errWrite != nil {
		fmt.Println("error on writing response", errWrite)
	}
}

var Resources = map[string]func(req *Request){
	"/ping": PingHandler,
}

func Handler(msgRaw string, c *connection.Connection) {
	msg, err := mjtp.Parse(msgRaw)

	if err != nil {
		errBody := make(map[string]interface{})
		errBody["message"] = "invalid message format"

		errMsg, errMakeMsg := mjtp.Make("/invalid", errBody).String()

		if errMakeMsg != nil {
			fmt.Println("error on make error message")
			return
		}

		sendRawResponse(errMsg, c)
		return
	}

	if handler, ok := Resources[msg.Resource]; ok {
		req := &Request{
			Data:       msg,
			Connection: c,
			response:   []Response{},
		}

		handler(req)
	} else {
		fmt.Printf("resource %s not found", msg.Resource)
	}
}
