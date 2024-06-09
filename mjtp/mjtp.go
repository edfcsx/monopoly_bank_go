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
	jsonMap, err := json.Marshal(m.Body)

	if err != nil {
		return "", err
	}

	return m.Resource + " " + m.Version + " " + string(jsonMap) + "\r\n\r\n", nil
}

func Parse(data string) (*Message, error) {
	message := &Message{}
	message.Body = make(map[string]interface{})

	// Split the data into 3 parts
	parts := strings.Split(data, " ")

	if len(parts) < 3 {
		return nil, errors.New("invalid message format")
	}

	message.Resource = parts[0]
	message.Version = parts[1]

	if message.Version != "MJTP/1.0" {
		return nil, errors.New("invalid version")
	}

	// Parse the body
	bodyParts := strings.Split(parts[2], "\r\n\r\n")

	if len(bodyParts) < 2 {
		return nil, errors.New("invalid body format")
	}

	if err := json.Unmarshal([]byte(bodyParts[0]), &message.Body); err != nil {
		return nil, err
	}

	return message, nil
}
