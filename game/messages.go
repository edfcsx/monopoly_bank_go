package game

const (
	AuthenticateCommand Command = "Authenticate"
	PingCommand         Command = "Ping"
	ProfileCommand      Command = "SendProfile"
	TransferCommand     Command = "Transfer"
)

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
