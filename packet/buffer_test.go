package packet_test

import (
	"testing"
	"time"

	"github.com/byonchev/go-engine.io/packet"
	"github.com/stretchr/testify/assert"
)

func TestBufferPop(t *testing.T) {
	buffer := packet.NewBuffer()

	expected := packet.NewPong(nil)

	buffer.Add(expected)

	actual := buffer.Pop()

	assert.Equal(t, expected, actual, "packet was not added to the buffer")
}

func TestBufferFlush(t *testing.T) {
	buffer := packet.NewBuffer()

	p1 := packet.NewPong(nil)
	p2 := packet.NewStringMessage("hello")

	expected := packet.Payload{p1, p2}

	buffer.Add(p1)
	buffer.Add(p2)

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "flush doesn't return buffered packets")
}

func TestBufferFlushAdd(t *testing.T) {
	buffer := packet.NewBuffer()

	p1 := packet.NewPong(nil)

	buffer.Add(p1)
	buffer.Flush()

	p2 := packet.NewStringMessage("hello")

	expected := packet.Payload{p2}

	buffer.Add(p2)

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "flush doesn't return buffered packets")
}

func TestBufferCloseSinglePacket(t *testing.T) {
	buffer := packet.NewBuffer()

	p1 := packet.NewPong(nil)

	buffer.Close()
	buffer.Add(p1)

	expected := packet.Payload(nil)

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "invalid packets in closed buffer")
}

func TestBufferCloseMultiplePackets(t *testing.T) {
	buffer := packet.NewBuffer()

	p1 := packet.NewPong(nil)
	p2 := packet.NewStringMessage("hello")
	p3 := packet.NewStringMessage("world")

	buffer.Add(p1)
	buffer.Add(p2)
	buffer.Close()
	buffer.Add(p3)

	expected := packet.Payload{p1, p2}

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "invalid packets in closed buffer")
}

func TestBufferWait(t *testing.T) {
	buffer := packet.NewBuffer()

	flushed := false

	go func() {
		buffer.Flush()
		flushed = true
	}()

	time.Sleep(100 * time.Millisecond)

	assert.True(t, !flushed, "buffer flush doesn't wait for at least one packet to be added")
}
