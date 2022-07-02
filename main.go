package main

import (
	"fmt"
	"log"

	"github.com/davyj0nes/generic-web-handler/client"
	"github.com/davyj0nes/generic-web-handler/server"
)

func main() {
	port := "8080"
	s := server.NewServer(port)
	go func() {
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	addr := fmt.Sprintf("http://localhost:%s/sum", port)
	client.NewClient(addr, nil).Run()
}
