package server

import (
	"errors"
	"fmt"
	"monopoly_bank_go/connection"
	"monopoly_bank_go/types"
	"net"
	"os"
	"sync"
)

type Acceptor struct {
	port        uint
	ln          net.Listener
	acceptorCb  func(c *connection.Connection)
	stopChan    chan bool
	running     bool
	connections map[string]*connection.Connection
	mutex       sync.Mutex
}

type TCPServer struct {
	acceptors map[types.Protocol]*Acceptor
}

var server TCPServer
var IsOpen bool

func init() {
	server.acceptors = make(map[types.Protocol]*Acceptor)
	IsOpen = true
}

func NewTCPServer(p types.Protocol, port uint, cb func(c *connection.Connection)) {
	_, ok := server.acceptors[p]

	if ok {
		fmt.Println("[Server - " + string(p) + "] already started")
		return
	}

	host := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Println("[Server - "+string(p)+"] error on starting Server", err)
		os.Exit(1)
	}

	s := &Acceptor{
		port:        port,
		ln:          listener,
		acceptorCb:  cb,
		stopChan:    make(chan bool),
		running:     true,
		connections: make(map[string]*connection.Connection),
	}

	server.acceptors[p] = s
	go startServer(p)
}

func startServer(p types.Protocol) {
	s, ok := server.acceptors[p]

	if !ok {
		fmt.Println("[Server - " + p.String() + "] not found")
		return
	}

	fmt.Println("[Server - "+p.String()+"] started on port: ", s.port)

	for {
		select {
		case <-s.stopChan:
			fmt.Println("Stopping Server on port: ", s.port)
			s.running = false

			fmt.Println("[Server - " + p.String() + "] closed to accept new connections")
			fmt.Println("[Server - " + p.String() + "] closing all connections")

			s.mutex.Lock()

			for _, connections := range s.connections {
				err := connections.Socket.Close()

				if err != nil {
					fmt.Println("[Server - "+p.String()+"] error on closing connection", err)
				}
			}

			s.connections = make(map[string]*connection.Connection)
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

			c := connection.MakeConnection(p, conn)
			server.PushConnection(p, c)

			go s.acceptorCb(c)
		}
	}
}

func (s *TCPServer) PushConnection(p types.Protocol, c *connection.Connection) {
	attrs, ok := s.acceptors[p]

	if ok {
		attrs.mutex.Lock()
		attrs.connections[c.Id] = c
		attrs.mutex.Unlock()
	}
}

func StopTCPServer(p types.Protocol) {
	h, ok := server.acceptors[p]

	if !ok {
		fmt.Println("[Server - " + string(p) + "] not found")
		return
	}

	err := h.ln.Close()
	h.stopChan <- true

	if err != nil {
		fmt.Println("[Server - "+p.String()+"] error on closing", err)
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

func CloseConnection(c *connection.Connection) {
	h, ok := server.acceptors[c.Protocol]

	if !ok {
		fmt.Println("[Server - " + string(c.Protocol) + "] not found to close connection")
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if conn, ok := h.connections[c.Id]; ok {
		err := conn.Socket.Close()

		if err != nil {
			fmt.Println("[Server - "+c.Protocol.String()+"] error on closing connection", err)
		}

		delete(h.connections, c.Id)
		return
	}
}

func GetConnections(p types.Protocol) map[string]*connection.Connection {
	server.acceptors[p].mutex.Lock()
	defer server.acceptors[p].mutex.Unlock()

	return server.acceptors[p].connections
}

func DumpClosedConnections() {
	for _, acceptor := range server.acceptors {
		acceptor.mutex.Lock()

		for id, conn := range acceptor.connections {
			if conn.IsClosed {
				delete(acceptor.connections, id)
			}
		}

		acceptor.mutex.Unlock()
	}
}
