package connection

import (
	"fmt"
	"github.com/google/uuid"
	"monopoly_bank_go/types"
	"net"
	"time"
)

type Connection struct {
	Id       string
	Protocol types.Protocol
	Socket   net.Conn
	IsClosed bool
	PlayerId string
}

func MakeConnection(protocol types.Protocol, socket net.Conn) *Connection {
	return &Connection{
		Id:       uuid.New().String(),
		Protocol: protocol,
		Socket:   socket,
		IsClosed: false,
		PlayerId: "",
	}
}

func (c *Connection) Close() {
	if c.IsClosed {
		return
	}

	c.IsClosed = true
	err := c.Socket.Close()

	if err != nil {
		fmt.Println("error on closing connection", err)
	}
}

func (c *Connection) SendAndClose(data string) {
	err := c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		fmt.Println("error on setting write deadline", err)
		c.Close()
		return
	}

	_, err = c.Socket.Write([]byte(data))

	if err != nil {
		fmt.Println("error on writing response", err)
		c.Close()
		return
	}

	c.Close()
}

func (c *Connection) Send(data string) {
	err := c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		fmt.Println("error on setting write deadline", err)
		c.Close()
		return
	}

	_, err = c.Socket.Write([]byte(data))

	if err != nil {
		fmt.Println("error on writing response", err)
		c.Close()
		return
	}
}
