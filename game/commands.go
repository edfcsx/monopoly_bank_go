package game

import (
	"encoding/json"
	"fmt"
	"monopoly_bank_go/connection"
	"monopoly_bank_go/server"
	"monopoly_bank_go/types"
	"monopoly_bank_go/websocket"
	"strings"
)

type Command string
type CommandResponse string

var commandsHandler = map[Command]func(conn *CmdConnection){
	AuthenticateCommand: AuthenticateCommandHandler,
	PingCommand:         PingCommandHandler,
	ProfileCommand:      ProfileCommandHandler,
	TransferCommand:     TransferCommandHandler,
}

type CmdConnection struct {
	Raw        string
	RawJson    map[string]interface{}
	Cmd        Command
	Connection *connection.Connection
}

func (c *CmdConnection) parse() error {
	com := strings.Split(c.Raw, "|")

	if len(com) < 1 {
		return fmt.Errorf("command expected but not found")
	}

	if len(com[1]) > 0 {
		err := json.Unmarshal([]byte(com[1]), &c.RawJson)

		if err != nil {
			return err
		}
	}

	c.Cmd = Command(com[0])
	return nil
}

func GetRequest[T any](c *CmdConnection) (*Request[T], error) {
	var data T

	cmd := strings.Split(c.Raw, "|")

	if len(cmd) == 2 && len(cmd[1]) > 0 {
		err := json.Unmarshal([]byte(cmd[1]), &data)

		if err != nil {
			return nil, err
		}
	}

	req := &Request[T]{
		Raw:        c.Raw,
		RawJson:    c.RawJson,
		Cmd:        c.Cmd,
		Connection: c.Connection,
		Data:       &data,
		res:        []Response{},
	}

	return req, nil
}

func AcceptCommand(commandRaw string, c *connection.Connection) {
	fmt.Println("AcceptCommand", commandRaw)

	conn := &CmdConnection{
		Raw:        commandRaw,
		Connection: c,
		RawJson:    map[string]interface{}{},
	}

	err := conn.parse()

	if err != nil {
		fmt.Println("error on parsing command", err)
		SendRawResponse(BadRequest, "", c)
		return
	}

	if commandFunc, ok := commandsHandler[conn.Cmd]; ok {
		commandFunc(conn)
	}
}

func SendResponse[T any](req *Request[T]) {
	for _, res := range req.res {

		if arg, ok := req.RawJson["args_id"]; ok {
			if !res.SendToAll {
				res.Data["args_id"] = arg
			}
		}

		dataBytes, err := json.Marshal(res.Data)

		if err != nil {
			fmt.Println("error on parsing response data", err)
			return
		}

		response := fmt.Sprint(res.Cmd, "|", string(dataBytes))
		fmt.Println("response", response)

		var errW *websocket.FrameError

		if res.SendToAll {
			for _, v := range server.GetConnections(types.WEBSOCKET) {
				errW = websocket.Write(v.Socket, response, websocket.TextFrame)

				if errW != nil {
					fmt.Println("error on writing response", errW)
					v.Close()
				}
			}
		} else {
			errW = websocket.Write(req.Connection.Socket, response, websocket.TextFrame)
		}

		if errW != nil {
			fmt.Println("error on writing response", errW)
		}
	}
}

func SendRawResponse(cmd CommandResponse, data string, c *connection.Connection) {
	response := fmt.Sprint(cmd, "|", data)
	errW := websocket.Write(c.Socket, response, websocket.TextFrame)

	if errW != nil {
		fmt.Println("error on writing response", errW)
	}
}
