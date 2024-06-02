package game

import "fmt"

type ProfileReq struct {
	PlayerHash string `json:"player_hash"`
}

func ProfileCommandHandler(cmd *CmdConnection) {
	req, err := GetRequest[ProfileReq](cmd)

	if err != nil {
		fmt.Println("error on parsing profile command", err)
		SendRawResponse(BadRequest, "", cmd.Connection)
		return
	}

	if _, ok := req.RawJson["player_hash"]; !ok {
		fmt.Println("profile command: player_hash not found")
		SendRawResponse(BadRequest, "", cmd.Connection)
		return
	}

	player := PlayerExistsByHash(req.Data.PlayerHash)

	if player == nil {
		fmt.Println("profile command: player not found")
		SendRawResponse(BadRequest, "", cmd.Connection)
		return
	}

	req.AppendResponse(
		string(ProfileData),
		map[string]interface{}{
			"balance":        player.Balance,
			"players_online": GetOnlinePlayersNames(),
		},
		false,
	)

	SendResponse(req)
}
