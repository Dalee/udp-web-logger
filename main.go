package main

import (
	"time"
	"udp-web-logger/pkg"
)

func main() {
	udp := pkg.NewUDPServer("127.0.0.1:9010")
	web := pkg.NewHTTPServer("127.0.0.1:9000")

	go web.Serve()

	defer udp.Shutdown()
	defer web.Shutdown()

	go func() {
		for {
			message, err := udp.ReadMessage(2048)
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
