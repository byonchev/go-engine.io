package transport

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
)

const xhrFrameWindow = 50 * time.Millisecond

// XHR is the standard polling transport
type XHR struct {
	codec codec.Codec

	readLock  sync.Mutex
	writeLock sync.Mutex

	sendChannel    <-chan packet.Packet
	receiveChannel chan<- packet.Packet
}

// NewXHR creates new XHR transport
func NewXHR(sendChannel <-chan packet.Packet, receiveChannel chan<- packet.Packet) *XHR {
	return &XHR{
		codec:          codec.XHR{},
		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
	}
}

// HandleRequest handles HTTP polling requests
func (xhr *XHR) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	method := request.Method

	switch method {
	case "GET":
		xhr.write(writer)
	case "POST":
		xhr.read(request.Body)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (xhr *XHR) read(reader io.Reader) {
	xhr.readLock.Lock()
	defer xhr.readLock.Unlock()

	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return
	}

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
