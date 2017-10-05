package pkg

import (
	"time"
)

type Message struct {
	IP      string    `json:"ip"`
	Payload string    `json:"payload"`
	Time    time.Time `json:"time"`
}
