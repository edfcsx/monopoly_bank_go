package commands

import (
	"monopoly_bank_go/accounts"
	"monopoly_bank_go/mjtp"
)

func StatusHandler(request *Request) {
	account := accounts.ExistsByHash(request.Connection.PlayerId)

	if account == nil {
		request.Respond(makeErrorResponse("Seus dados de conta não foram encontrados."), false)
		request.Respond(makeForceLogoutResponse(), false)
		request.SendResponses()
		return
	}

	response := mjtp.Make("/status", map[string]interface{}{
		"balance":        account.Balance,
		"players_online": accounts.All(),
	})

	request.Respond(response, false)
	request.SendResponses()
}
