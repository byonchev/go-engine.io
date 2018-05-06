package packet

import (
	"sync"
)

// Buffer is a synchronized buffer of packets
type Buffer struct {
	sync.Mutex

	flushCondition *sync.Cond
	flushLimit     int

	payload Payload

	closed bool
}

// NewBuffer returns new packet buffer with fixed flush limit
func NewBuffer(flushLimit int) *Buffer {
	buffer := &Buffer{flushLimit: flushLimit, closed: false}
	buffer.flushCondition = sync.NewCond(buffer)

	return buffer
}

// Add adds new packet to the payload buffer
func (buffer *Buffer) Add(packet Packet) {
	buffer.Lock()
	defer buffer.Unlock()

	if buffer.closed {
		return
	}

	buffer.payload = append(buffer.payload, packet)

	buffer.flushCondition.Broadcast()
}

// Close stops the buffering of packets
func (buffer *Buffer) Close() {
	buffer.Lock()
	defer buffer.Unlock()

	if buffer.closed {
		return
	}

	buffer.closed = true

	buffer.flushCondition.Broadcast()
}

// Flush returns and clears the buffered payload.
// If the buffer is empty, it blocks until at least one packet is present
func (buffer *Buffer) Flush() Payload {
	buffer.Lock()
	defer buffer.Unlock()

	for len(buffer.payload) == 0 && !buffer.closed {
		buffer.flushCondition.Wait()
	}

	length := len(buffer.payload)
	limit := length

	if buffer.flushLimit > 0 && buffer.flushLimit < length && !buffer.closed {
		limit = buffer.flushLimit
	}

	payload := buffer.payload[:limit]

	buffer.payload = buffer.payload[limit:]

	return payload
}
