package pkg

import (
	"bytes"
	"log"
	"net"
	"os"
	"time"
)

type UDPServer struct {
	conn   *net.UDPConn
	Logger *log.Logger
}

func (u *UDPServer) ReadMessage(bufferSize int) (*Message, error) {
	buf := make([]byte, bufferSize)

	_, remote, err := u.conn.ReadFromUDP(buf)

	if err != nil {
		u.Logger.Println(err)
		return nil, err
	}

	message := &Message{
		IP:      remote.IP.String(),
		Payload: string(bytes.Trim(buf[:], "\x00")),
		Time:    time.Now(),
	}

	return message, err
}

func (u *UDPServer) Shutdown() {
	u.conn.Close()
}

func NewUDPServer(addr string) *UDPServer {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		panic(err)
	}

	return &UDPServer{
		Logger: log.New(os.Stdout, "udp > ", log.Ldate|log.Ltime),
		conn:   conn,
	}
}
