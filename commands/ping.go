package commands

import "monopoly_bank_go/mjtp"

func PingHandler(request *Request) {
	response := mjtp.Make("/ping", nil)
	request.Respond(response, false)
	request.SendResponses()
}
