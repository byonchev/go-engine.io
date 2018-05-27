package main

import (
	"fmt"
	"net/http"

	"github.com/byonchev/go-engine.io"
)

func main() {
	engineIO := eio.NewServer()

	events := engineIO.Events()

	go func() {
		for event := range events {
			switch event := event.(type) {
			case eio.MessageEvent:
				fmt.Printf("Message received from %s: %s\n", event.SessionID, string(event.Data))
			case eio.ConnectEvent:
				fmt.Printf("Client %s connected\n", event.SessionID)

				engineIO.Send(event.SessionID, false, []byte("Hello"))
			case eio.DisconnectEvent:
				fmt.Printf("Client %s disconnected\n", event.SessionID)
			}
		}
	}()

	http.Handle("/engine.io/", engineIO)
	http.ListenAndServe(":8080", nil)
}
