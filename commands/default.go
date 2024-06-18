package commands

import "monopoly_bank_go/mjtp"

func makeErrorResponse(message string) *mjtp.Message {
	return mjtp.Make("/error", map[string]interface{}{
		"message": message,
	})
}

func makeForceLogoutResponse() *mjtp.Message {
	return mjtp.Make("/force_logout", nil)
}

func makeGlobalMessageResponse(message string) *mjtp.Message {
	return mjtp.Make("/global_message", map[string]interface{}{
		"message": message,
	})
}
