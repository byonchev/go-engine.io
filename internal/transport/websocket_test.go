package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/byonchev/go-engine.io/internal/codec"
	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/byonchev/go-engine.io/internal/transport"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebsocketSend(t *testing.T) {
	codec, transport, server, client := setupWebsockets()
	defer server.Close()

	tests := []struct {
		packet      packet.Packet
		messageType int
	}{
		{
			packet.NewStringMessage("Hello"),
			websocket.TextMessage,
		},
		{
			packet.NewBinaryMessage([]byte("world")),
			websocket.BinaryMessage,
		},
	}

	for _, test := range tests {
		transport.Send(test.packet)

		messageType, data, _ := client.ReadMessage()

		payload, _ := codec.Decode(bytes.NewBuffer(data))

		assert.Equal(t, test.messageType, messageType, "wrong message type received by client")
		assert.Equal(t, packet.Payload{test.packet}, payload, "packet was not received by client")
	}
}

func TestWebsocketReceive(t *testing.T) {
	codec, transport, server, client := setupWebsockets()
	defer server.Close()

	tests := []struct {
		packet      packet.Packet
		messageType int
	}{
		{
			packet.NewStringMessage("Hello"),
			websocket.TextMessage,
		},
		{
			packet.NewBinaryMessage([]byte("world")),
			websocket.BinaryMessage,
		},
	}

	for _, test := range tests {
		var buffer bytes.Buffer

		codec.Encode(packet.Payload{test.packet}, &buffer)
		client.WriteMessage(test.messageType, buffer.Bytes())

		actual, _ := transport.Receive()

		assert.Equal(t, test.packet, actual, "packet was not received from client")
	}
}

func TestWebsocketSendAfterShutdown(t *testing.T) {
	_, transport, server, client := setupWebsockets()
	defer server.Close()

	transport.Shutdown()
	err := transport.Send(packet.NewNOOP())

	expected := []byte(nil)
	_, actual, _ := client.ReadMessage()

	assert.Error(t, err, "error was not returned after send on stopped transport")
	assert.Equal(t, expected, actual, "packet was sent to the client after shutdown")
}

func TestWebsocketReceiveAfterShutdown(t *testing.T) {
	codec, transport, server, client := setupWebsockets()
	defer server.Close()

	var buffer bytes.Buffer

	codec.Encode(packet.Payload{packet.NewNOOP()}, &buffer)
	client.WriteMessage(websocket.TextMessage, buffer.Bytes())

	transport.Shutdown()
	actual, err := transport.Receive()

	assert.Error(t, err, "error was not returned after receive on stopped transport")
	assert.Equal(t, packet.Packet{}, actual, "packet was received from client after shutdown")
}

func TestWebsocketUpgradeError(t *testing.T) {
	transport := createWebsocketTransport()

	request, _ := http.NewRequest("POST", "/", nil)
	writer := httptest.NewRecorder()

	transport.HandleRequest(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code, "upgrade failure doesn't return 400")
}

func createWebsocketTransport() *transport.Websocket {
	return transport.NewWebsocket(1024, 1024, false, func(*http.Request) bool { return true })
}

func createServer(transport *transport.Websocket) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(transport.HandleRequest))
}

func connectClient(server *httptest.Server) *websocket.Conn {
	url := "ws" + strings.TrimPrefix(server.URL, "http")

	client, _, _ := websocket.DefaultDialer.Dial(url, nil)

	return client
}

func setupWebsockets() (codec.Codec, *transport.Websocket, *httptest.Server, *websocket.Conn) {
	codec := codec.Websocket{}
	transport := createWebsocketTransport()
	server := createServer(transport)
	client := connectClient(server)

	return codec, transport, server, client
}
