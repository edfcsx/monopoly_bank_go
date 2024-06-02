package static_files

import (
	"fmt"
	"monopoly_bank_go/connection"
	"monopoly_bank_go/http"
	"os"
	"strings"
)

var rootPath string

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

	rootPath = dir + "/www"
}

func Handler(c *connection.Connection) {
	go http.HandlerRequest(c.Socket, func(r *http.Request, err error) {
		if err != nil {
			c.SendAndClose(http.MakeResponse(http.BadRequest, nil, ""))
			return
		}

		send(c, r)
	})
}

func send(c *connection.Connection, r *http.Request) {
	if r.Method != "GET" {
		c.SendAndClose(http.MakeResponse(http.BadRequest, nil, ""))
		return
	}

	if r.Path == "/" {
		r.Path = "/index.html"
	}

	if !strings.Contains(r.Path, ".") {
		r.Path = r.Path + ".html"
	}

	path := rootPath + r.Path

	if !fileExists(path) {
		c.SendAndClose(http.MakeResponse(http.NotFound, nil, ""))
		return
	}

	content, err := readFile(path)

	if err != nil {
		c.SendAndClose(http.MakeResponse(http.InternalServerError, nil, ""))
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

	c.SendAndClose(http.MakeResponse(http.OK, headers, content))
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
