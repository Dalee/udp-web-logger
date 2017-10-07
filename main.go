package main

import (
	"os"
	"time"
	"udp-web-logger/pkg"
)

func main() {
	udp := pkg.NewUDPServer(os.Getenv("UDP_LISTEN"))
	web := pkg.NewHTTPServer(os.Getenv("WEB_LISTEN"), 5)

	go web.Serve()

	defer udp.Shutdown()
	defer web.Shutdown()

	go func() {
		for {
			message, err := udp.ReadMessage(4096)
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
