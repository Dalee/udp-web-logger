package pkg

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type HTTPServer struct {
	Server   *http.Server
	Logger   *log.Logger
	Messages []*Message
}

func (h *HTTPServer) Serve() {
	h.Logger.Println("Serving http://" + h.Server.Addr)
	if err := h.Server.ListenAndServe(); err != nil {
		h.Logger.Println(err)
	}
}

func (h *HTTPServer) HandleLog(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	b, _ := json.Marshal(h.Messages)
	w.Write(b)
}

func (h *HTTPServer) Shutdown() {
	if err := h.Server.Shutdown(nil); err != nil {
		panic(err)
	}
}

func (h *HTTPServer) AddMessage(message *Message) {
	h.Messages = append(h.Messages, message)
}

func NewHTTPServer(addr string) *HTTPServer {
	server := &HTTPServer{
		Logger:   log.New(os.Stdout, "http > ", log.Ldate|log.Ltime),
		Messages: make([]*Message, 0),
		Server:   &http.Server{Addr: addr},
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	http.HandleFunc("/api/log", server.HandleLog)

	return server
}
