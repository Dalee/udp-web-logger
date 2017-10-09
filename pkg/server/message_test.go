package server

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageMarshal(t *testing.T) {
	now := time.Date(
		2009, 11, 17, 20,
		34, 58, 651387237, time.UTC,
	)

	message := &Message{IP: "127.0.0.1", Payload: "some test message", Time: now}
	expectedMessage := &jsonMessage{
		Time: "Nov 17 20:34:58",
		Unix: 1258490098,
		Alias: &Alias{
			IP:      "127.0.0.1",
			Payload: "some test message",
		},
	}

	encodedMessage, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}

	var decodedMessage jsonMessage
	err = json.Unmarshal(encodedMessage, &decodedMessage)

	if err != nil {
		t.Fatal(err)
	}

	if expectedMessage.Unix != decodedMessage.Unix {
		t.Fatalf("Unix times don't match: %d and %d.", expectedMessage.Unix, decodedMessage.Unix)
	}

	if expectedMessage.Time != decodedMessage.Time {
		t.Fatalf("Times don't match: %s and %s.", expectedMessage.Time, decodedMessage.Time)
	}

	if expectedMessage.IP != decodedMessage.IP {
		t.Fatalf("IPs don't match: %s and %s.", expectedMessage.IP, decodedMessage.IP)
	}

	if expectedMessage.Payload != decodedMessage.Payload {
		t.Fatalf("Payloads don't match: %s and %s.", expectedMessage.Payload, decodedMessage.Payload)
	}
}
