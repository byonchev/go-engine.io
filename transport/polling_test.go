package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/transport"
	"github.com/stretchr/testify/assert"
)

func TestPollingSendBufferedPayload(t *testing.T) {
	codec := codec.XHR{}
	transport := createPollingTransport()

	packets := []packet.Packet{
		packet.NewStringMessage("hello"),
		packet.NewStringMessage("world"),
	}

	for _, packet := range packets {
		transport.Send(packet)
	}

	received := <-clientReceive(transport)

	expected := packet.Payload(packets)
	actual, _ := codec.Decode(received)

	assert.Equal(t, expected, actual, "packets were not delivered to client")
}

func TestPollingSendPayloadAfterRequest(t *testing.T) {
	codec := codec.XHR{}
	transport := createPollingTransport()

	sent := packet.NewClose()

	transport.Send(sent)

	received := <-clientReceive(transport)

	expected := packet.Payload{sent}
	actual, _ := codec.Decode(received)

	assert.Equal(t, expected, actual, "packets were not delivered to client")
}

func TestPollingSendAndShutdown(t *testing.T) {
	transport := createPollingTransport()

	transport.Send(packet.NewNOOP())
	transport.Shutdown()

	expected := []byte(nil)
	actual := (<-clientReceive(transport)).Bytes()

	assert.Equal(t, expected, actual, "packets were sent to the client after shutdown")
}

func TestPollingReceivePayload(t *testing.T) {
	codec := codec.XHR{}
	transport := transport.NewPolling(0, 10, nil)

	payload := packet.Payload{
		packet.NewStringMessage("hello"),
		packet.NewStringMessage("world"),
	}

	var buffer bytes.Buffer

	codec.Encode(payload, &buffer)

	clientSend(transport, &buffer)

	for _, expected := range payload {
		actual, err := transport.Receive()

		assert.NoError(t, err, "error while receiving sent packets")
		assert.Equal(t, expected, actual, "packets sents from client were not received")
	}
}

func TestPollingReceiveAndShutdown(t *testing.T) {
	codec := codec.XHR{}
	transport := transport.NewPolling(0, 10, nil)

	sent := packet.NewNOOP()

	var buffer bytes.Buffer

	codec.Encode(packet.Payload{sent}, &buffer)

	clientSend(transport, &buffer)
	transport.Shutdown()

	expected := sent
	actual, err := transport.Receive()

	assert.NoError(t, err, "error while receiving sent packets")
	assert.Equal(t, expected, actual, "payload was not received due to transport shutdown")
}

func TestPollingInvalidHTTPMethod(t *testing.T) {
	transport := createPollingTransport()

	request, _ := http.NewRequest("DELETE", "/", nil)
	writer := httptest.NewRecorder()

	transport.HandleRequest(writer, request)

	assert.Equal(t, http.StatusMethodNotAllowed, writer.Code, "http handler responded to invalid method")
}

func TestPollingReceiveInvalidPayload(t *testing.T) {
	transport := createPollingTransport()

	buffer := bytes.NewBuffer([]byte("INVALID:INVALID"))

	clientSend(transport, buffer)

	go func() {
		transport.Receive()
		t.Error("invalid received packet was processed")
	}()

	time.Sleep(100 * time.Millisecond)
}

func createPollingTransport() *transport.Polling {
	return transport.NewPolling(0, 0, nil)
}

func clientReceive(transport transport.Transport) <-chan *bytes.Buffer {
	result := make(chan *bytes.Buffer)

	request, _ := http.NewRequest("GET", "/", nil)
	writer := httptest.NewRecorder()

	go func() {
		transport.HandleRequest(writer, request)
		result <- writer.Body
		close(result)
	}()

	return result
}

func clientSend(transport transport.Transport, buffer *bytes.Buffer) {
	request, _ := http.NewRequest("POST", "/", buffer)
	writer := httptest.NewRecorder()

	transport.HandleRequest(writer, request)

	time.Sleep(100 * time.Millisecond)
}
