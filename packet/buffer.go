package packet

import (
	"sync"
)

// Buffer is a synchronized buffer of packets
type Buffer struct {
	writeLock sync.Mutex
	flushLock sync.Mutex

	payload Payload

	flushable bool
	closed    bool
}

// NewBuffer returns new packet buffer with fixed size
func NewBuffer() *Buffer {
	buffer := &Buffer{closed: false, flushable: false}
	buffer.flushLock.Lock()

	return buffer
}

// Add adds new packet to the payload buffer
func (buffer *Buffer) Add(packet Packet) {
	buffer.writeLock.Lock()
	defer buffer.writeLock.Unlock()

	if buffer.closed {
		return
	}

	buffer.payload = append(buffer.payload, packet)

	if !buffer.flushable {
		buffer.flushLock.Unlock()
		buffer.flushable = true
	}
}

// Pop blocks untils a packet is buffered and then returns it
func (buffer *Buffer) Pop() Packet {
	buffer.writeLock.Lock()
	defer buffer.writeLock.Unlock()

	buffer.flushLock.Lock()

	length := len(buffer.payload)

	if length > 1 {
		defer buffer.flushLock.Unlock()
	} else {
		buffer.flushable = false
	}

	packet := buffer.payload[length-1]
	buffer.payload = buffer.payload[:length-1]

	return packet
}

// Close stops the buffering of packets
func (buffer *Buffer) Close() {
	buffer.writeLock.Lock()
	defer buffer.writeLock.Unlock()

	buffer.closed = true

	if !buffer.flushable {
		buffer.flushLock.Unlock()
	}
}

// Flush returns and clears the buffered payload.
// If the buffer is empty, it blocks until at least one packet is present
func (buffer *Buffer) Flush() Payload {
	if !buffer.closed {
		buffer.flushable = false
		buffer.flushLock.Lock()
	}

	buffer.writeLock.Lock()
	defer buffer.writeLock.Unlock()

	payload := buffer.payload

	buffer.payload = nil

	return payload
}
