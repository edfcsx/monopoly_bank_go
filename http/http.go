package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Response struct {
	Status  int
	Headers map[string]string
	Body    string
	Close   bool
}

type StatusCode int

var statusText = map[StatusCode]string{
	101: "Switching Protocols",
	200: "OK",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

const (
	SwitchingProtocols  StatusCode = 101
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	Unauthorized        StatusCode = 401
	NotFound            StatusCode = 404
	InternalServerError StatusCode = 500
)

var CorsHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "*",                                                                                  // Permite solicitações de qualquer origem
	"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",                                                    // Permite os métodos HTTP especificados
	"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", // Permite os cabeçalhos HTTP especificados
}

func MakeResponse(status StatusCode, headers map[string]string, body string) string {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, statusText[status])

	for key, value := range headers {
		response += key + ": " + value + "\r\n"
	}

	if _, ok := headers["Content-Length"]; !ok {
		response += "Content-Length: " + fmt.Sprint(len(body)) + "\r\n"
	}

	if _, ok := headers["Connection"]; !ok {
		response += "Connection: close\r\n"
	}

	for key, value := range CorsHeaders {
		response += key + ": " + value + "\r\n"
	}

	response += "\r\n" + body

	return response
}

type Request struct {
	socket  net.Conn
	Method  string
	Path    string
	Headers map[string]string
	Body    string
	Data    string
	Query   map[string]string
}

func HandlerRequest(socket net.Conn, cb func(r *Request, err error)) {
	request := &Request{
		socket:  socket,
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}

	err := request.parse()

	if err != nil {
		cb(nil, err)
		return
	}

	cb(request, nil)
}

// @TODO: refactor parseHeaders method in parseStatusLine and parseHeaders
func (r *Request) parseHeaders() {
	lines := strings.Split(r.Data, "\r\n")
	status := strings.Split(lines[0], " ")

	r.Method = status[0]
	r.Path = status[1]

	parseQuery(r)

	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			break
		}

		header := strings.Split(lines[i], ": ")
		if len(header) == 1 {
			header[1] = ""
		}

		r.Headers[header[0]] = header[1]
	}
}

func (r *Request) parse() error {
	reader := bufio.NewReader(r.socket)
	var buffer strings.Builder

	for {
		b, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		buffer.WriteString(b)

		if strings.HasSuffix(buffer.String(), "\r\n\r\n") || err == io.EOF {
			r.Data = buffer.String()
			break
		}
	}

	r.parseHeaders()

	if r.Headers["Content-Length"] != "" {
		length, _ := strconv.Atoi(r.Headers["Content-Length"])

		if length > 0 {
			body := make([]byte, length)
			_, err := io.ReadFull(reader, body)

			if err != nil {
				return err
			}

			r.Body = string(body)
		}
	}

	return nil
}

func parseQuery(r *Request) {
	idx := strings.Index(r.Path, "?")

	if idx != -1 {
		rawQuery := r.Path[idx+1:]
		r.Path = r.Path[:idx]
		query := strings.Split(rawQuery, "&")

		for _, params := range query {
			values := strings.Split(params, "=")

			if len(values) > 1 {
				r.Query[values[0]] = values[1]
			} else {
				r.Query[values[0]] = ""
			}
		}
	}
}
