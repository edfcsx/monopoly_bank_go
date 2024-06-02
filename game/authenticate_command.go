package game

import "fmt"

type AuthenticateReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthenticateCommandHandler(cmd *CmdConnection) {
	req, err := GetRequest[AuthenticateReq](cmd)

	if err != nil {
		fmt.Println("error on parsing authenticate command", err)
		SendRawResponse(BadRequest, "", cmd.Connection)
		return
	}

	player := PlayerExistsByName(req.Data.Username)
	res := Response{}

	if player == nil {
		playerHash := CreatePlayer(req.Data.Username, req.Data.Password)

		res.Cmd = string(AuthenticateSuccess)
		res.Data["player_hash"] = playerHash

		connAlert := Response{}
		connAlert.Cmd = string(GlobalMessage)
		connAlert.SendToAll = true
		connAlert.Data["message"] = fmt.Sprintf("O jogador %s se juntou ao jogo!", req.Data.Username)
		req.AddResponse(connAlert)
	} else {
		if player.Password != req.Data.Password {
			res.Cmd = string(AuthenticateFailed)
		} else {
			res.Cmd = string(AuthenticateSuccess)
			res.Data["player_hash"] = player.AuthHash
		}
	}

	req.AddResponse(res)
	SendResponse(req)
}
