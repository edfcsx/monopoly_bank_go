package mjtp

import (
	"encoding/json"
	"errors"
	"strings"
)

type Message struct {
	Resource string
	Version  string
	Body     map[string]interface{}
}

func Make(resource string, body map[string]interface{}) *Message {
	return &Message{
		Resource: resource,
		Version:  "MJTP/1.0",
		Body:     body,
	}
}

func (m *Message) String() (string, error) {
	body := "{}"

	if m.Body != nil {
		j, err := json.Marshal(m.Body)

		if err != nil {
			return "", err
		}
		body = string(j)
	}

	return m.Resource + " " + m.Version + " " + body + "\r\n\r\n", nil
}

func Parse(data string) (*Message, error) {
	message := &Message{}
	message.Body = make(map[string]interface{})

	buffer := data
	bufferIdx := 0

	// Find the first space
	bufferIdx = strings.Index(buffer, " ")
	message.Resource = buffer[:bufferIdx]

	// Find the second space
	buffer = buffer[bufferIdx+1:]
	bufferIdx = strings.Index(buffer, " ")
	message.Version = buffer[:bufferIdx]

	// Find the first \r\n\r\n
	buffer = buffer[bufferIdx+1:]
	bufferIdx = strings.Index(buffer, "\r\n\r\n")

	if err := json.Unmarshal([]byte(buffer[:bufferIdx]), &message.Body); err != nil {

		return nil, err
	}
	if message.Version != "MJTP/1.0" {
		return nil, errors.New("invalid version")
	}

	return message, nil
}
