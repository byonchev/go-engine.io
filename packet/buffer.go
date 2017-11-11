package packet

import (
	"sync"
)

// Payload is a collection of packets
type Payload []Packet

// Buffer is a synchronized buffer of packets
type Buffer struct {
	sync.Mutex

	closed  bool
	packets chan Packet
}

// NewBuffer returns new packet buffer with fixed size
func NewBuffer(size int) Buffer {
	return Buffer{
		packets: make(chan Packet, size),
		closed:  false,
	}
}

// Add adds new packet to the payload buffer
func (buffer *Buffer) Add(packet Packet) {
	buffer.Lock()
	defer buffer.Unlock()

	if buffer.closed {
		return
	}

	buffer.packets <- packet
}

// Pop blocks untils a packet is buffered and then returns it
func (buffer *Buffer) Pop() Packet {
	return <-buffer.packets
}

// Close stops the buffering of packets
func (buffer *Buffer) Close() {
	buffer.Lock()
	defer buffer.Unlock()

	buffer.closed = true

	close(buffer.packets)
}

// Flush returns and clears the buffered payload.
// If the buffer is empty, it blocks until at least one packet is present
func (buffer *Buffer) Flush() Payload {
	packet, ok := <-buffer.packets

	if !ok {
		return nil
	}

	payload := Payload{packet}

	buffer.Lock()
	defer buffer.Unlock()

	for {
		select {
		case packet, ok = <-buffer.packets:
			if !ok {
				return payload
			}

			payload = append(payload, packet)
		default:
			return payload
		}
	}
}
