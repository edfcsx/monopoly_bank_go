package Server

import (
	"fmt"
	"net"
	"time"
)

type Connection struct {
	id       string
	protocol Protocol
	Socket   net.Conn
	Player   string
}

func (c *Connection) SendAndClose(data string) {
	err := c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		fmt.Println("error on setting write deadline", err)
		CloseConnection(c)
		return
	}

	_, err = c.Socket.Write([]byte(data))

	if err != nil {
		fmt.Println("error on writing response", err)
		CloseConnection(c)
		return
	}

	CloseConnection(c)
}

func (c *Connection) Close() {
	CloseConnection(c)
}

func (c *Connection) Send(data string) {
	err := c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		fmt.Println("error on setting write deadline", err)
		CloseConnection(c)
		return
	}

	_, err = c.Socket.Write([]byte(data))

	if err != nil {
		fmt.Println("error on writing response", err)
		CloseConnection(c)
		return
	}
}
