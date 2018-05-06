package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/transport"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketSend(t *testing.T) {
	codec, transport, server, client := setupWebSockets()
	defer server.Close()

	expected := packet.NewClose()

	transport.Send(expected)

	_, data, _ := client.ReadMessage()

	payload, _ := codec.Decode(bytes.NewBuffer(data))

	assert.Equal(t, payload[0], expected, "packet was not received by client")
}

func TestWebSocketReceive(t *testing.T) {
	codec, transport, server, client := setupWebSockets()
	defer server.Close()

	expected := packet.NewStringMessage("Test")

	var buffer bytes.Buffer

	codec.Encode(packet.Payload{expected}, &buffer)
	client.WriteMessage(websocket.TextMessage, buffer.Bytes())

	actual, _ := transport.Receive()

	assert.Equal(t, expected, actual, "packet was not received from client")
}

func TestWebSocketShutdown(t *testing.T) {
	_, transport, server, client := setupWebSockets()
	defer server.Close()

	transport.Shutdown()
	transport.Send(packet.NewNOOP())

	expected := []byte(nil)
	_, actual, _ := client.ReadMessage()

	assert.Equal(t, expected, actual, "packets were sent to the client after shutdown")
}

func createServer(transport transport.Transport) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(transport.HandleRequest))
}

func connectClient(server *httptest.Server) *websocket.Conn {
	url := "ws" + strings.TrimPrefix(server.URL, "http")

	client, _, _ := websocket.DefaultDialer.Dial(url, nil)

	return client
}

func setupWebSockets() (codec.Codec, *transport.WebSocket, *httptest.Server, *websocket.Conn) {
	codec := codec.WebSocket{}
	transport := transport.NewWebSocket()
	server := createServer(transport)
	client := connectClient(server)

	return codec, transport, server, client
}
