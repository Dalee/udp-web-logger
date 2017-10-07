package pkg

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

// HTTPServer encapsulates HTTP server API.
type HTTPServer struct {
	server      *http.Server
	messages    []*Message
	logger      *log.Logger
	maxMessages int
}

// Serve starts serving requests.
func (h *HTTPServer) Serve() {
	h.logger.Println("Serving http://" + h.server.Addr)

	if err := h.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// HandleLog replies with messages taken from UDP.
func (h *HTTPServer) HandleLog(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	lastTimestamp := r.URL.Query().Get("ts")

	if lastTimestamp == "" {
		b, _ := json.Marshal(h.messages)
		w.Write(b)
		return
	}

	ts, _ := strconv.ParseInt(lastTimestamp, 10, 64)
	lastMessages := make([]*Message, 0)

	for _, message := range h.messages {
		if message.Time.Unix() > ts {
			lastMessages = append([]*Message{message}, lastMessages...)
		}
	}

	b, _ := json.Marshal(lastMessages)
	w.Write(b)
}

// Shutdown stops http listener.
func (h *HTTPServer) Shutdown() {
	if err := h.server.Shutdown(nil); err != nil {
		panic(err)
	}
}

// AddMessage prepends message to the slice.
func (h *HTTPServer) AddMessage(message *Message) {
	if len(h.messages) >= h.maxMessages {
		h.messages = []*Message{message}
	} else {
		h.messages = append([]*Message{message}, h.messages...)
	}
}

// NewHTTPServer creates new HTTP server.
func NewHTTPServer(addr string, maxMessages int) *HTTPServer {
	server := &HTTPServer{
		logger:      log.New(os.Stdout, "http > ", log.Ldate|log.Ltime),
		messages:    make([]*Message, 0),
		server:      &http.Server{Addr: addr},
		maxMessages: maxMessages,
	}

	fs := http.FileServer(http.Dir("public"))

	http.Handle("/", fs)
	http.HandleFunc("/api/log", server.HandleLog)

	return server
}
