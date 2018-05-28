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

	engineIO := eio.NewServer()
	engineIO.SetLogger(logger)

	engineIO.PingInterval = 5 * time.Second
	engineIO.PingTimeout = 10 * time.Second
	engineIO.AllowUpgrades = false
	engineIO.Transports = []string{"polling"}
	engineIO.PollingBufferFlushLimit = 100
	engineIO.PollingBufferReceiveLimit = 50

	events := engineIO.Events()

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

	http.Handle("/engine.io/", engineIO)
	http.ListenAndServe(":8080", nil)
}
