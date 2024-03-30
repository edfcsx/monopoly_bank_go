package WebsocketHandler

import (
	"crypto/sha1"
	"encoding/base64"
)

const MagicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func makeHandshakeResponseKey(key string) string {
	combined := key + MagicGUID

	hasher := sha1.New()
	hasher.Write([]byte(combined))
	hash := hasher.Sum(nil)

	output := base64.StdEncoding.EncodeToString(hash)
	return output
}
