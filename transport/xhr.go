package transport

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"../codec"
	"../packet"
)

const xhrFrameWindow = 50 * time.Millisecond

// XHR is the standard polling transport
type XHR struct {
	codec codec.Codec

	readLock  sync.Mutex
	writeLock sync.Mutex

	sendChannel    chan packet.Packet
	receiveChannel chan packet.Packet
}

// NewXHR creates new XHR transport
func NewXHR() *XHR {
	return &XHR{
		codec:          codec.XHR{},
		sendChannel:    make(chan packet.Packet),
		receiveChannel: make(chan packet.Packet),
	}
}

// HandleRequest handles HTTP polling requests
func (xhr *XHR) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		xhr.read(request.Body)
	} else {
		xhr.write(writer)
	}
}

// Send enqueues packet for sending
func (xhr *XHR) Send(packet packet.Packet) {
	go func() { xhr.sendChannel <- packet }()
}

// Receive blocks until a packet is received
func (xhr *XHR) Receive() packet.Packet {
	return <-xhr.receiveChannel
}

func (xhr *XHR) read(reader io.Reader) {
	xhr.readLock.Lock()
	defer xhr.readLock.Unlock()

	data, _ := ioutil.ReadAll(reader)

	payload, err := xhr.codec.Decode(data)

	if err != nil {
		return
	}

	for _, packet := range payload {
		xhr.receiveChannel <- packet
	}
}

func (xhr *XHR) write(writer io.Writer) {
	xhr.writeLock.Lock()
	defer xhr.writeLock.Unlock()

	payload := xhr.createPayload()

	data := xhr.codec.Encode(payload)

	writer.Write(data)
}

func (xhr *XHR) createPayload() packet.Payload {
	var payload packet.Payload

	timer := time.NewTimer(xhrFrameWindow)

	for {
		select {
		case packet := <-xhr.sendChannel:
			payload = append(payload, packet)
			continue
		case <-timer.C:
		}

		if len(payload) == 0 {
			timer.Reset(xhrFrameWindow)
			continue
		}

		return payload
	}
}
