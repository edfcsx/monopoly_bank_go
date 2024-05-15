export type json = string;

export enum CommandsResponse {
	AuthenticateFailed         = "AuthenticateFailed",
	AuthenticateSuccess        = "AuthenticateSuccess",
	Pong                       = "Pong",
	ProfileData                = "ProfileData",
	TransferSuccess            = "TransferSuccess",
	TransferFailed             = "TransferFailed",
	TransferInsufficientFunds  = "TransferInsufficientFunds",
	TransferReceived           = "TransferReceived",
	BadRequest                 = "BadRequest",
	GlobalMessage              = "GlobalMessage"
}

export enum CommandsRequest {
	Authenticate               = "Authenticate",
	Ping                       = "Ping",
	SendProfile                = "SendProfile",
	Transfer                   = "Transfer"
}
