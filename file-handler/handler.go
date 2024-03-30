package FileHandler

import (
	"fmt"
	"monopoly_bank_go/http"
	"monopoly_bank_go/server"
	"os"
	"strings"
)

var ROOT_DIR string

var contentType = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "text/javascript",
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"gif":  "image/gif",
	"ico":  "image/x-icon",
	"svg":  "image/svg+xml",
}

func init() {
	dir, err := os.Getwd()

	if err != nil {
		fmt.Println("error on getting working directory", err)
		os.Exit(1)
	}

	ROOT_DIR = dir + "/www"
}

func AcceptConnection(c *Server.Connection) {
	go HTTP.HandlerRequest(c.Socket, func(r *HTTP.Request, err error) {
		if err != nil {
			fmt.Println("error on handling request", err)
			c.SendAndClose(HTTP.MakeResponse(HTTP.BadRequest, nil, ""))
			return
		}

		sendFile(c, r)
	})
}

func sendFile(c *Server.Connection, r *HTTP.Request) {
	if r.Method != "GET" {
		c.SendAndClose(HTTP.MakeResponse(HTTP.BadRequest, nil, ""))
		return
	}

	if r.Path == "/" {
		r.Path = "/index.html"
	}

	if !strings.Contains(r.Path, ".") {
		r.Path = r.Path + ".html"
	}

	path := ROOT_DIR + r.Path

	if !fileExists(path) {
		c.SendAndClose(HTTP.MakeResponse(HTTP.NotFound, nil, ""))
		return
	}

	content, err := readFile(path)

	if err != nil {
		c.SendAndClose(HTTP.MakeResponse(HTTP.InternalServerError, nil, ""))
		return
	}

	headers := map[string]string{
		"Cache-Control": "no-cache, no-store, must-revalidate",
		"Pragma":        "no-cache",
		"Expires":       "0",
	}

	ext := strings.Split(r.Path, ".")[1]
	if extHeader, ok := contentType[ext]; ok {
		headers["Content-Type"] = extHeader
	}

	c.SendAndClose(HTTP.MakeResponse(HTTP.OK, headers, content))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

func readFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("error on reading file", err)
		return "", err
	}

	return string(bytes), nil
}
