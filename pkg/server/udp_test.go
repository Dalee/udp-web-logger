package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"testing"
)

var (
	server *UDPServer
	conn   net.Conn
)

func setup() {
	var err error
	logger := log.New(ioutil.Discard, "", 0)
	server = NewUDPServer("localhost:9999", logger)
	raddr := server.conn.LocalAddr()

	if conn, err = net.Dial(raddr.Network(), raddr.String()); err != nil {
		panic(err)
	}
}

func teardown() {
	server.Shutdown()
}

func TestReadMessage(t *testing.T) {
	setup()
	defer teardown()

	bytesPayload := []byte("Some test message")
	conn.Write(bytesPayload)

	message, err := server.ReadMessage(512)

	if err != nil {
		t.Fatal(err)
	}

	payload := string(bytes.Trim(bytesPayload[:], "\x00"))

	if message.Payload != payload {
		t.Fatal("Payloads don't match")
	}
}
