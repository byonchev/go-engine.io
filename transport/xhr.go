package transport

import (
	"io"
	"net/http"
	"sync"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
)

// XHR is the standard polling transport
type XHR struct {
	running bool

	codec  codec.Codec
	buffer packet.Buffer

	receiving sync.WaitGroup

	sendChannel    <-chan packet.Packet
	receiveChannel chan<- packet.Packet
}

// NewXHR creates new XHR transport
func NewXHR(sendChannel <-chan packet.Packet, receiveChannel chan<- packet.Packet) *XHR {
	transport := &XHR{
		codec:          codec.XHR{},
		buffer:         packet.NewBuffer(10),
		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
		running:        true,
	}

	go transport.bufferPackets()

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

// Shutdown stops the transport from receiving or sending packets
func (transport *XHR) Shutdown() {
	transport.running = false

	transport.receiving.Wait()
	transport.buffer.Close()
}

func (transport *XHR) read(reader io.Reader) {
	payload, err := transport.codec.Decode(reader)

	if err != nil {
		return
	}

	transport.receiving.Add(1)

	for _, packet := range payload {
		transport.receiveChannel <- packet
	}

	transport.receiving.Done()
}

func (transport *XHR) write(writer io.Writer) {
	payload := transport.buffer.Flush()

	err := transport.codec.Encode(payload, writer)

	if err != nil {
		return
	}
}

func (transport *XHR) bufferPackets() {
	for packet := range transport.sendChannel {
		transport.buffer.Add(packet)
	}
}
