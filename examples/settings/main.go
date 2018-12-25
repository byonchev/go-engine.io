package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/byonchev/go-engine.io"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	server := eio.NewServer()
	server.Configure(
		eio.Logger(logger),
		eio.PingInterval(5*time.Second),
		eio.PingTimeout(10*time.Second),
		eio.DisableUpgrades(),
		eio.Transports("polling"),
		eio.PollingBufferFlushLimit(100),
		eio.PollingBufferReceiveLimit(50),
	)

	events := server.Events()

	go func() {
		for event := range events {
			switch event := event.(type) {
			case eio.MessageEvent:
				fmt.Printf("Message received from %s: %s\n", event.SessionID, string(event.Data))
			case eio.ConnectEvent:
				fmt.Printf("Client %s connected\n", event.SessionID)
			case eio.DisconnectEvent:
				fmt.Printf("Client %s disconnected. Reason: %s\n", event.SessionID, event.Reason)
			}
		}
	}()

	http.Handle("/engine.io/", server)
	http.ListenAndServe(":8080", nil)
}
