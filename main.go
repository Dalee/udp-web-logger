package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"udp-web-logger/pkg/server"
)

var usage = `
Usage: udp-web-logger [options].

Options:
--udp-read-buffer-size   size of buffer to read incoming UDP packet into. Default: 4096.
--max-messages           maximum amount of messages to keep. Default: 50.
--help                   prints this message.

Env:
UDP_LISTEN - address to listen UDP on. Default: 127.0.0.1:9010.
WEB_LISTEN - address to listen HTTP on. Default: 127.0.0.1:9000.

`

func main() {
	udpReadBufferSize := flag.Int("udp-read-buffer-size", 4096, "")
	maxMessages := flag.Int("max-messages", 50, "")
	udpListen := os.Getenv("UDP_LISTEN")
	webListen := os.Getenv("WEB_LISTEN")

	flag.Usage = func() {
		fmt.Print(usage)
	}

	flag.Parse()

	if udpListen == "" {
		udpListen = "127.0.0.1:9010"
	}

	if webListen == "" {
		webListen = "127.0.0.1:9000"
	}

	udp := server.NewUDPServer(udpListen, nil)
	web := server.NewHTTPServer(webListen, *maxMessages)

	defer udp.Shutdown()
	defer web.Shutdown()

	go web.Serve()

	go func() {
		for {
			message, err := udp.ReadMessage(*udpReadBufferSize)
			if err != nil {
				continue
			}

			web.AddMessage(message)
		}
	}()

	for {
		time.Sleep(time.Microsecond * 100)
	}
}
