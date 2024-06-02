package types

type Protocol uint8

const (
	FILE      Protocol = 0x1
	WEBSOCKET Protocol = 0x2
	WEBAPI    Protocol = 0x3
)

func (p Protocol) String() string {
	switch p {
	case 0x1:
		return "FILE"
	case 0x2:
		return "WEBSOCKET"
	case 0x3:
		return "WEBAPI"
	default:
		return "UNKNOWN"
	}
}
