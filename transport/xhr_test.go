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

const waitTime = 100 * time.Millisecond

func TestXHRSendBufferedPayload(t *testing.T) {
	sendChannel := make(chan packet.Packet, 10)

	codec := codec.XHR{}
	transport := transport.NewXHR(sendChannel, nil)

	packets := []packet.Packet{
		packet.NewStringMessage("hello"),
		packet.NewStringMessage("world"),
	}

	for _, packet := range packets {
		sendChannel <- packet
	}

	close(sendChannel)

	// ensure all packets are buffered
	time.Sleep(waitTime)

	received := <-clientReceive(transport)

	expected := packet.Payload(packets)
	actual, _ := codec.Decode(received)

	assert.Equal(t, expected, actual, "packets were not delivered to client")
}

func TestXHRSendPayloadAfterRequest(t *testing.T) {
	sendChannel := make(chan packet.Packet, 10)

	codec := codec.XHR{}
	transport := transport.NewXHR(sendChannel, nil)

	sent := packet.NewClose()

	go func() {
		// ensure http request is sent
		time.Sleep(waitTime)
		sendChannel <- packet.NewClose()
	}()

	received := <-clientReceive(transport)

	expected := packet.Payload{sent}
	actual, _ := codec.Decode(received)

	assert.Equal(t, expected, actual, "packets were not delivered to client")
}

func TestXHRSendAndShutdown(t *testing.T) {
	sendChannel := make(chan packet.Packet, 10)

	transport := transport.NewXHR(sendChannel, nil)

	sendChannel <- packet.NewNOOP()
	transport.Shutdown()

	expected := []byte(nil)
	actual := <-clientReceive(transport)

	assert.Equal(t, expected, actual, "packets were sent to the client after shutdown")
}

func TestXHRReceivePayload(t *testing.T) {
	receiveChannel := make(chan packet.Packet, 10)

	codec := codec.XHR{}
	transport := transport.NewXHR(nil, receiveChannel)

	packets := []packet.Packet{
		packet.NewStringMessage("hello"),
		packet.NewStringMessage("world"),
	}

	sent := codec.Encode(packet.Payload(packets))

	clientSend(transport, sent)

	for _, expected := range packets {
		actual := <-receiveChannel

		assert.Equal(t, expected, actual, "packets were not received from client")
	}
}

func TestXHRReceiveAndShutdown(t *testing.T) {
	receiveChannel := make(chan packet.Packet)

	codec := codec.XHR{}
	transport := transport.NewXHR(nil, receiveChannel)

	sent := packet.NewNOOP()

	go func() {
		clientSend(transport, codec.Encode(packet.Payload{sent}))
		transport.Shutdown()
	}()

	// ensure packet is buffered and shutdown sequence is initiated
	time.Sleep(waitTime)

	expected := sent
	actual, _ := <-receiveChannel

	assert.Equal(t, expected, actual, "payload was not received due to transport shutdown")
}

func clientReceive(transport transport.Transport) <-chan []byte {
	result := make(chan []byte)

	request, _ := http.NewRequest("GET", "/", nil)
	writer := httptest.NewRecorder()

	go func() {
		transport.HandleRequest(writer, request)
		result <- writer.Body.Bytes()
		close(result)
	}()

	return result
}

func clientSend(transport transport.Transport, data []byte) {
	buffer := bytes.NewBuffer(data)

	request, _ := http.NewRequest("POST", "/", buffer)
	writer := httptest.NewRecorder()

	transport.HandleRequest(writer, request)
}
