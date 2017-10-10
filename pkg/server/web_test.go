package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

var (
	httpServer *HTTPServer
)

func setupHTTP(maxMessages int) {
	logger := log.New(ioutil.Discard, "", 0)
	httpServer = NewHTTPServer("localhost:12051", maxMessages, logger)
}

func teardownHTTP() {
	httpServer.Shutdown()
}

func makeTestMessage() *Message {
	return &Message{IP: "127.0.0.1", Payload: "test message", Time: time.Now()}
}

func makeTestMessageWithTime(t time.Time) *Message {
	return &Message{IP: "127.0.0.1", Payload: "test message", Time: t}
}

func addTestMessages(n int) {
	for i := 0; i < n; i++ {
		message := makeTestMessage()
		httpServer.AddMessage(message)
	}
}

func requestTestMessages(ts int64, t *testing.T) []*jsonMessage {
	u := fmt.Sprintf("http://%s/api/log?ts=%d", httpServer.server.Addr, ts)
	resp, err := http.Get(u)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var messages []*jsonMessage
	if err := json.Unmarshal(b, &messages); err != nil {
		t.Fatal(err)
	}

	return messages
}

func TestHTTPServerAddMessagePrepend(t *testing.T) {
	setupHTTP(5)
	defer teardownHTTP()

	firstMessage := makeTestMessage()
	secondMessage := makeTestMessage()

	httpServer.AddMessage(firstMessage)
	httpServer.AddMessage(secondMessage)

	if httpServer.messages[0] != secondMessage {
		t.Fatal("Last added message is not the first in slice")
	}
}

func TestHTTPServerAddMessageTruncate(t *testing.T) {
	setupHTTP(1)
	defer teardownHTTP()

	addTestMessages(2)

	if len(httpServer.messages) != 1 {
		t.Fatal("Messages are not truncated")
	}
}

func TestHTTPServerHandleLogAllMessages(t *testing.T) {
	setupHTTP(5)
	defer teardownHTTP()

	addTestMessages(5)

	go httpServer.Serve()

	messages := requestTestMessages(0, t)

	if len(messages) != 5 {
		t.Errorf("/api/log?ts= is expected to return 5 message, actual: %d", len(messages))
	}
}

func TestHTTPServerHandleLogAllMessagesTsFilter(t *testing.T) {
	setupHTTP(5)
	defer teardownHTTP()

	go httpServer.Serve()

	firstMessage := makeTestMessageWithTime(
		time.Date(
			2009, 11, 17, 20,
			34, 58, 651387237, time.UTC,
		),
	)

	secondMessage := makeTestMessageWithTime(
		time.Date(
			2015, 11, 17, 20,
			34, 58, 651387237, time.UTC,
		),
	)

	httpServer.AddMessage(firstMessage)
	httpServer.AddMessage(secondMessage)

	messages := requestTestMessages(firstMessage.Time.Unix()+1, t)

	if len(messages) != 1 {
		t.Errorf("/api/log?ts > 0 is expected to return 1 message, actual: %d", len(messages))
	}
}
