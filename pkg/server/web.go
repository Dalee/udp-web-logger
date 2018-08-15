package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"udp-web-logger/pkg/static"
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
		h.logger.Println(err)
	}
}

// HandleLog replies with messages taken from UDP.
func (h *HTTPServer) HandleLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	lastTimestampQuery := r.URL.Query().Get("ts")
	ts, _ := strconv.ParseInt(lastTimestampQuery, 10, 64)

	if ts == 0 {
		b, _ := json.Marshal(h.messages)
		w.Write(b)
		return
	}

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
func NewHTTPServer(addr string, maxMessages int, logger *log.Logger) *HTTPServer {
	if logger == nil {
		logger = log.New(os.Stdout, "http > ", log.Ldate|log.Ltime)
	}

	mux := http.NewServeMux()

	server := &HTTPServer{
		logger:      logger,
		messages:    make([]*Message, 0),
		server:      &http.Server{Addr: addr, Handler: mux},
		maxMessages: maxMessages,
	}

	mux.Handle("/", http.FileServer(static.HTTP))
	mux.HandleFunc("/api/log", server.HandleLog)

	return server
}
