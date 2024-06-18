package commands

import "monopoly_bank_go/accounts"

func AuthenticateHandler(request *Request) {
	if id, ok := request.Data.Body["id"]; ok {
		request.Connection.PlayerId = id.(string)

		account := accounts.ExistsByHash(request.Connection.PlayerId)

		if account != nil {
			request.Respond(makeGlobalMessageResponse(account.Name+", se juntou ao jogo!"), true)
		}

		request.SendResponses()
	} else {
		request.Respond(makeErrorResponse("Conexão não autorizada"), false)
		request.Respond(makeForceLogoutResponse(), false)
		request.SendResponses()
		return
	}
}
