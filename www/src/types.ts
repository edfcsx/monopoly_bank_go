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
	GlobalMessage              = "GlobalMessage",
	OnlinePlayersData					 = "OnlinePlayersData"
}

export enum CommandsRequest {
	Authenticate               = "Authenticate",
	Ping                       = "Ping",
	SendProfile                = "SendProfile",
	Transfer                   = "Transfer",
	OnlinePlayers              = "OnlinePlayers"
}

export interface NetworkingMessage {
	command: CommandsRequest | CommandsResponse,
	args?: { [key:string]: any }
	[key: string]: any
}

export abstract class Commands {
	abstract execute (serverMessage: NetworkingMessage): void;
}
