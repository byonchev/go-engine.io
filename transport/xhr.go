package transport

import (
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/logger"
	"github.com/byonchev/go-engine.io/packet"
)

// XHR is the standard polling transport
type XHR struct {
	running bool

	codec  codec.Codec
	buffer *packet.Buffer

	receiving sync.WaitGroup

	received chan packet.Packet
}

// NewXHR creates new XHR transport
func NewXHR() *XHR {
	transport := &XHR{
		codec:    codec.XHR{},
		buffer:   packet.NewBuffer(),
		received: make(chan packet.Packet),
		running:  true,
	}

	return transport
}

// HandleRequest handles HTTP polling requests
func (transport *XHR) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	method := request.Method

	if !transport.running {
		return
	}

	switch method {
	case "GET":
		transport.write(writer)
	case "POST":
		transport.read(request.Body)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Send buffers packets for sending on next poll cycle
func (transport *XHR) Send(packet packet.Packet) error {
	transport.buffer.Add(packet)

	return nil
}

// Receive returns the last received packet or blocks until a packet is present
func (transport *XHR) Receive() (packet.Packet, error) {
	received, success := <-transport.received

	if !success {
		return packet.Packet{}, errors.New("transport is stopped")
	}

	return received, nil
}

// Shutdown stops the transport from receiving or sending packets
func (transport *XHR) Shutdown() {
	if !transport.running {
		return
	}

	transport.running = false

	transport.receiving.Wait()
	transport.buffer.Close()

	close(transport.received)
}

func (transport *XHR) read(reader io.Reader) {
	payload, err := transport.codec.Decode(reader)

	if err != nil {
		logger.Error("error while receiving messages:", err)
		return
	}

	transport.receiving.Add(1)

	for _, packet := range payload {
		transport.received <- packet
	}

	transport.receiving.Done()
}

func (transport *XHR) write(writer io.Writer) {
	payload := transport.buffer.Flush()

	err := transport.codec.Encode(payload, writer)

	if err != nil {
		logger.Error("error while sending messages:", err)
		return
	}
}
