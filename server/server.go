package Server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"os"
	"sync"
)

type Protocol string

const (
	FILE      Protocol = "FILE"
	WEBSOCKET Protocol = "WEBSOCKET"
)

type Attrs struct {
	port        string
	ln          net.Listener
	acceptorCb  func(c *Connection)
	stopChan    chan bool
	running     bool
	connections map[string]*Connection
	mutex       sync.Mutex
}

type TCPServer struct {
	acceptors map[Protocol]*Attrs
}

var server TCPServer
var IsOpen bool

func init() {
	server.acceptors = make(map[Protocol]*Attrs)
	IsOpen = true
}

func NewTCPServer(p Protocol, port string, cb func(c *Connection)) {
	_, ok := server.acceptors[p]

	if ok {
		fmt.Println("[Server - " + string(p) + "] already started")
		return
	}

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println("[Server - "+string(p)+"] error on starting Server", err)
		os.Exit(1)
	}

	s := &Attrs{
		port:        port,
		ln:          listener,
		acceptorCb:  cb,
		stopChan:    make(chan bool),
		running:     true,
		connections: make(map[string]*Connection),
	}

	server.acceptors[p] = s
	go startServer(p)
}

func startServer(p Protocol) {
	s, ok := server.acceptors[p]

	if !ok {
		fmt.Println("[Server - " + p + "] not found")
		return
	}

	fmt.Println("[Server - "+p+"] started on port: ", s.port)

	for {
		select {
		case <-s.stopChan:
			fmt.Println("Stopping Server on port: ", s.port)
			s.running = false

			fmt.Println("[Server - " + p + "] closed to accept new connections")
			fmt.Println("[Server - " + p + "] closing all connections")

			s.mutex.Lock()

			for _, connection := range s.connections {
				err := connection.Socket.Close()

				if err != nil {
					fmt.Println("[Server - "+p+"] error on closing connection", err)
				}
			}

			s.connections = make(map[string]*Connection)
			s.mutex.Unlock()
			return
		default:
			conn, err := s.ln.Accept()

			if err != nil {
				var opErr *net.OpError
				if errors.As(err, &opErr) && opErr.Op == "accept" && !s.running {
					return
				}

				fmt.Println("error on accepting connection: ", err)
				continue
			}

			s.mutex.Lock()

			c := Connection{
				id:       uuid.New().String(),
				protocol: p,
				Socket:   conn,
			}

			s.connections[c.id] = &c
			s.mutex.Unlock()
			go s.acceptorCb(&c)
		}
	}
}

func StopTCPServer(p Protocol) {
	h, ok := server.acceptors[p]

	if !ok {
		fmt.Println("[Server - " + string(p) + "] not found")
		return
	}

	err := h.ln.Close()
	h.stopChan <- true

	if err != nil {
		fmt.Println("[Server - "+p+"] error on closing", err)
		panic(err)
	}

	var count int

	for _, v := range server.acceptors {
		if v.running {
			count++
		}
	}

	if count == 0 {
		IsOpen = false
	}
}

func CloseConnection(c *Connection) {
	h, ok := server.acceptors[c.protocol]

	if !ok {
		fmt.Println("[Server - " + string(c.protocol) + "] not found to close connection")
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if conn, ok := h.connections[c.id]; ok {
		err := conn.Socket.Close()

		if err != nil {
			fmt.Println("[Server - "+c.protocol+"] error on closing connection", err)
		}

		delete(h.connections, c.id)
		return
	}
}

func GetConnections(p Protocol) map[string]*Connection {
	server.acceptors[p].mutex.Lock()
	defer server.acceptors[p].mutex.Unlock()

	return server.acceptors[p].connections
}
