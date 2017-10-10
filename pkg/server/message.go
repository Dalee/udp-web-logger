package server

import (
	"encoding/json"
	"time"
)

// Message is a struct that is created from UDP and
// served through HTTP.
type Message struct {
	IP      string    `json:"ip"`
	Payload string    `json:"payload"`
	Time    time.Time `json:"time"`
}

type alias Message

type jsonMessage struct {
	Time string `json:"time"`
	Unix int64  `json:"unix"`
	*alias
}

// MarshalJSON implements Marshaler interface to reformat
// time to human readable format.
func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonMessage{
		Time:  m.Time.Format(time.Stamp),
		Unix:  m.Time.Unix(),
		alias: (*alias)(m),
	})
}
