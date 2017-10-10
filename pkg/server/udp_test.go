package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"testing"
)

var (
	udpServer *UDPServer
	conn      net.Conn
)

func setupUDP() {
	var err error
	logger := log.New(ioutil.Discard, "", 0)
	udpServer = NewUDPServer("localhost:9999", logger)
	raddr := udpServer.conn.LocalAddr()

	if conn, err = net.Dial(raddr.Network(), raddr.String()); err != nil {
		panic(err)
	}
}

func teardownUDP() {
	udpServer.Shutdown()
}

func TestUDPServerReadMessage(t *testing.T) {
	setupUDP()
	defer teardownUDP()

	bytesPayload := []byte("Some test message")
	conn.Write(bytesPayload)

	message, err := udpServer.ReadMessage(512)

	if err != nil {
		t.Fatal(err)
	}

	payload := string(bytes.Trim(bytesPayload[:], "\x00"))

	if message.Payload != payload {
		t.Fatal("Payloads don't match")
	}
}
