package game

import "monopoly_bank_go/connection"

type Response struct {
	Cmd       string
	Data      map[string]interface{}
	SendToAll bool
}

type Request[T any] struct {
	Raw        string
	RawJson    map[string]interface{}
	Cmd        Command
	Connection *connection.Connection
	Data       *T
	res        []Response
}

func (r *Request[T]) AppendResponse(cmd string, data map[string]interface{}, sendToAll bool) {
	res := Response{
		Cmd:       cmd,
		Data:      data,
		SendToAll: sendToAll,
	}

	r.res = append(r.res, res)
}

func (r *Request[T]) AddResponse(res Response) {
	r.res = append(r.res, res)
}
