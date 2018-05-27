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

// Polling is the standard polling transport
type Polling struct {
	originCheck func(*http.Request) bool

	running bool

	buffer *packet.Buffer

	receiving sync.WaitGroup

	received chan packet.Packet
}

// NewPolling creates new polling transport
func NewPolling(bufferFlushLimit int, receiveBufferSize int, originCheck func(*http.Request) bool) *Polling {
	transport := &Polling{
		originCheck: originCheck,
		buffer:      packet.NewBuffer(bufferFlushLimit),
		received:    make(chan packet.Packet, receiveBufferSize),
		running:     true,
	}

	return transport
}

// HandleRequest handles HTTP polling requests
func (transport *Polling) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	if !transport.running {
		return
	}

	if !transport.originCheck(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	method := request.Method
	codec := transport.createCodec(request)

	switch method {
	case "GET":
		transport.write(writer, codec)
	case "POST":
		transport.read(request.Body, codec)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Send buffers packets for sending on next poll cycle
func (transport *Polling) Send(packet packet.Packet) error {
	transport.buffer.Add(packet)

	return nil
}

// Receive returns the last received packet or blocks until a packet is present
func (transport *Polling) Receive() (packet.Packet, error) {
	received, success := <-transport.received

	if !success {
		return packet.Packet{}, errors.New("transport is stopped")
	}

	return received, nil
}

// Shutdown stops the transport from receiving or sending packets
func (transport *Polling) Shutdown() {
	if !transport.running {
		return
	}

	transport.running = false

	transport.Send(packet.NewNOOP())

	transport.receiving.Wait()
	transport.buffer.Close()

	close(transport.received)
}

// Running returns true if the transport is active
func (transport *Polling) Running() bool {
	return transport.running
}

// Type returns the transport identifier
func (transport *Polling) Type() string {
	return PollingType
}

// Upgrades returns the possible transport upgrades
func (transport *Polling) Upgrades() []string {
	return []string{WebsocketType}
}

func (transport *Polling) read(reader io.Reader, codec codec.Codec) {
	payload, err := codec.Decode(reader)

	if err != nil {
		logger.Error("Error decoding messages:", err)
		return
	}

	transport.receiving.Add(1)

	for _, packet := range payload {
		transport.received <- packet
	}

	transport.receiving.Done()
}

func (transport *Polling) write(writer io.Writer, codec codec.Codec) {
	payload := transport.buffer.Flush()

	err := codec.Encode(payload, writer)

	if err != nil {
		logger.Error("Error encoding messages:", err)
		return
	}
}

func (transport *Polling) createCodec(request *http.Request) codec.Codec {
	query := request.URL.Query()

	b64 := query.Get("b64")
	j := query.Get("j")

	if j != "" {
		return codec.JSONP{Index: j}
	}

	return codec.XHR{ForceBase64: b64 != ""}
}
