package pkg

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

// MarshalJSON implements Marshaler interface to reformat
// time to human readable format.
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		Time string `json:"time"`
		Unix int64  `json:"unix"`
		*Alias
	}{
		Time:  m.Time.Format(time.Stamp),
		Unix:  m.Time.Unix(),
		Alias: (*Alias)(m),
	})
}
