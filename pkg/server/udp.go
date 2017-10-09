package server

import (
	"bytes"
	"log"
	"net"
	"os"
	"time"
)

// UDPServer encapsulates UDP server API.
type UDPServer struct {
	conn   *net.UDPConn
	logger *log.Logger
	addr   *net.UDPAddr
}

// ReadMessage reads bytes from UDP and serializes them into
// Message struct.
func (u *UDPServer) ReadMessage(bufferSize int) (*Message, error) {
	buf := make([]byte, bufferSize)

	_, remote, err := u.conn.ReadFromUDP(buf)

	if err != nil {
		u.logger.Println(err)
		return nil, err
	}

	message := &Message{
		IP:      remote.IP.String(),
		Payload: string(bytes.Trim(buf[:], "\x00")),
		Time:    time.Now(),
	}

	return message, err
}

// Shutdown closes UDP listener.
func (u *UDPServer) Shutdown() {
	if err := u.conn.Close(); err != nil {
		panic(err)
	}
}

// NewUDPServer creates new UDP server.
func NewUDPServer(addr string, logger *log.Logger) *UDPServer {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		panic(err)
	}

	if logger == nil {
		logger = log.New(os.Stdout, "udp > ", log.Ldate|log.Ltime)
	}

	logger.Println("Listening " + a.String())

	return &UDPServer{
		logger: logger,
		conn:   conn,
	}
}
