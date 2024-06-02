package game

import (
	"fmt"
)

func PingCommandHandler(cmd *CmdConnection) {
	req, err := GetRequest[Void](cmd)

	if err != nil {
		fmt.Println("error on parsing ping command", err)
		SendRawResponse(BadRequest, "", cmd.Connection)
	}

	req.AppendResponse(string(Pong), nil, false)
	SendResponse(req)
}
