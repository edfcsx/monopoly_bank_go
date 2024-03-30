package Game

import (
	"encoding/json"
	"fmt"
	Parse "monopoly_bank_go/parse"
	Server "monopoly_bank_go/server"
	"strings"
)

type Command string

const (
	AuthenticateCommand Command = "Authenticate"
	PingCommand         Command = "Ping"
	ProfileCommand      Command = "Profile"
	TransferCommand     Command = "Transfer"
)

type CommandResponse string

const (
	AuthenticateFailed        CommandResponse = "AuthenticateFailed"
	AuthenticateSuccess       CommandResponse = "AuthenticateSuccess"
	Pong                      CommandResponse = "Pong"
	ProfileData               CommandResponse = "ProfileData"
	TransferSuccess           CommandResponse = "TransferSuccess"
	TransferFailed            CommandResponse = "TransferFailed"
	TransferInsufficientFunds CommandResponse = "TransferInsufficientFunds"
	TransferReceived          CommandResponse = "TransferReceived"
	BadRequest                CommandResponse = "BadRequest"
	GlobalMessage             CommandResponse = "GlobalMessage"
)

var commandsHandler = map[Command]func(conn *CmdConnection){
	AuthenticateCommand: AuthenticateCommandHandler,
	PingCommand:         PingCommandHandler,
	ProfileCommand:      ProfileCommandHandler,
	TransferCommand:     TransferCommandHandler,
}

type Response struct {
	Cmd       string
	Data      map[string]interface{}
	SendToAll bool
}

type CmdConnection struct {
	Raw        string
	RawJson    map[string]interface{}
	Cmd        Command
	Connection Server.Connection
}

type Request[T any] struct {
	Raw        string
	RawJson    map[string]interface{}
	Cmd        Command
	Connection Server.Connection
	Data       *T
	Res        Response
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

func GetCmdRequest[T any](c *CmdConnection) (*Request[T], error) {
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
		Res:        Response{},
	}

	req.Res.Data = make(map[string]interface{})
	return req, nil
}

func AcceptCommand(commandRaw string, c *Server.Connection) {
	conn := &CmdConnection{
		Raw:        commandRaw,
		Connection: *c,
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

type AuthenticateReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthenticateCommandHandler(cmd *CmdConnection) {
	req, err := GetCmdRequest[AuthenticateReq](cmd)

	if err != nil {
		fmt.Println("error on parsing authenticate command", err)
		SendRawResponse(BadRequest, "", &cmd.Connection)
		return
	}

	player := PlayerExists(req.Data.Username)

	if player == nil {
		CreatePlayer(req.Data.Username, req.Data.Password)
		req.Res.Cmd = string(AuthenticateSuccess)
		req.Res.Data["message"] = fmt.Sprintf("O jogador %s se juntou ao jogo!", req.Data.Username)
		req.Res.SendToAll = true
	} else if player.Password != req.Data.Password {
		req.Res.Cmd = string(AuthenticateFailed)
	} else {
		req.Connection.Player = req.Data.Username
	}

	SendResponse(req)
}

func PingCommandHandler(cmd *CmdConnection) {
	req, err := GetCmdRequest[map[string]interface{}](cmd)

	if err != nil {
		fmt.Println("error on parsing ping command", err)
		SendRawResponse(BadRequest, "", &cmd.Connection)
	}

	req.Res.Cmd = string(Pong)
	SendResponse(req)
}

func ProfileCommandHandler(cmd *CmdConnection) {
	fmt.Println("ProfileCommandHandler", cmd)
}

func TransferCommandHandler(cmd *CmdConnection) {
	fmt.Println("TransferCommandHandler", cmd)
}

func SendResponse[T any](req *Request[T]) {
	if arg, ok := req.RawJson["args_id"]; ok {
		req.Res.Data["args_id"] = arg
	}

	dataBytes, err := json.Marshal(req.Res.Data)

	if err != nil {
		fmt.Println("error on parsing response data", err)
		return
	}

	response := fmt.Sprint(req.Res.Cmd, "|", string(dataBytes))

	var errW *Parse.FrameError

	if req.Res.SendToAll {
		for _, v := range Server.GetConnections(Server.WEBSOCKET) {
			errW = Parse.WriteToWebsocket(v.Socket, response, Parse.TextFrame)

			if errW != nil {
				fmt.Println("error on writing response", errW)
				v.Close()
			}
		}
	} else {
		errW = Parse.WriteToWebsocket(req.Connection.Socket, response, Parse.TextFrame)
	}

	if errW != nil {
		fmt.Println("error on writing response", errW)
	}
}

func SendRawResponse(cmd CommandResponse, data string, c *Server.Connection) {
	response := fmt.Sprint(cmd, "|", data)
	errW := Parse.WriteToWebsocket(c.Socket, response, Parse.TextFrame)

	if errW != nil {
		fmt.Println("error on writing response", errW)
	}
}
