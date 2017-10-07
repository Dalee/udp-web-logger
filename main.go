package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"udp-web-logger/pkg"
)

var usage = `
Usage: udp-web-logger [options].

Options:
--udp-read-buffer-size   size of buffer to read incoming UDP packet into. Default: 4096.
--max-messages           maximum amount of messages to keep. Default: 50.
--help                   prints this message.

Env:
UDP_LISTEN - address to listen UDP on.
WEB_LISTEN - address to listen HTTP on.

`

func main() {
	udpReadBufferSize := flag.Int("udp-read-buffer-size", 4096, "Buffer size to read into")
	maxMessages := flag.Int("max-messages", 50, "Maximum message to keep")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Print(usage)
		return
	}

	udp := pkg.NewUDPServer(os.Getenv("UDP_LISTEN"))
	web := pkg.NewHTTPServer(os.Getenv("WEB_LISTEN"), *maxMessages)

	go web.Serve()

	defer udp.Shutdown()
	defer web.Shutdown()

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
