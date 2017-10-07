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
	u.conn.Close()
}

// NewUDPServer creates new UDP server.
func NewUDPServer(addr string) *UDPServer {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "udp > ", log.Ldate|log.Ltime)
	logger.Println("Listening " + a.String())

	return &UDPServer{
		logger: logger,
		conn:   conn,
	}
}
