package packet_test

import (
	"testing"
	"time"

	"github.com/byonchev/go-engine.io/packet"
	"github.com/stretchr/testify/assert"
)

func TestBufferFlush(t *testing.T) {
	buffer := packet.NewBuffer(0)

	p1 := packet.NewPong(nil)
	p2 := packet.NewStringMessage("hello")

	expected := packet.Payload{p1, p2}

	buffer.Add(p1)
	buffer.Add(p2)

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "flush doesn't return buffered packets")
}

func TestBufferFlushLimit(t *testing.T) {
	buffer := packet.NewBuffer(1)

	p1 := packet.NewPong(nil)
	p2 := packet.NewStringMessage("hello")

	buffer.Add(p1)
	buffer.Add(p2)

	expected := packet.Payload{p1}
	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "flush doesn't limit buffered packets")

	expected = packet.Payload{p2}
	actual = buffer.Flush()

	assert.Equal(t, expected, actual, "invalid flush limit offset")
}

func TestBufferFlushAdd(t *testing.T) {
	buffer := packet.NewBuffer(0)

	p1 := packet.NewPong(nil)

	buffer.Add(p1)
	buffer.Flush()

	p2 := packet.NewStringMessage("hello")

	expected := packet.Payload{p2}

	buffer.Add(p2)

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "flush doesn't return buffered packets")
}

func TestBufferFlushWait(t *testing.T) {
	buffer := packet.NewBuffer(0)

	flushed := false

	go func() {
		buffer.Flush()
		flushed = true
	}()

	time.Sleep(100 * time.Millisecond)

	assert.False(t, flushed, "buffer flush doesn't wait for at least one packet to be added")
}

func TestBufferCloseSinglePacket(t *testing.T) {
	buffer := packet.NewBuffer(0)

	p1 := packet.NewPong(nil)

	buffer.Close()
	buffer.Add(p1)

	expected := packet.Payload(nil)

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "invalid packets in closed buffer")
}

func TestBufferCloseMultiplePackets(t *testing.T) {
	buffer := packet.NewBuffer(0)

	p1 := packet.NewPong(nil)
	p2 := packet.NewStringMessage("hello")
	p3 := packet.NewStringMessage("world")

	buffer.Add(p1)
	buffer.Add(p2)
	buffer.Close()
	buffer.Add(p3)
	buffer.Close()

	expected := packet.Payload{p1, p2}

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "invalid packets in closed buffer")
}

func TestBufferCloseLimit(t *testing.T) {
	buffer := packet.NewBuffer(1)

	p1 := packet.NewStringMessage("hello")
	p2 := packet.NewStringMessage("world")

	buffer.Add(p1)
	buffer.Add(p2)
	buffer.Close()

	expected := packet.Payload{p1, p2}

	actual := buffer.Flush()

	assert.Equal(t, expected, actual, "limited flush in closed buffer")
}

func BenchmarkBufferAdd(b *testing.B) {
	buffer := packet.NewBuffer(10)

	for i := 0; i < b.N; i++ {
		buffer.Add(packet.NewNOOP())
	}
}

func BenchmarkBufferFlush(b *testing.B) {
	buffer := packet.NewBuffer(1)

	packet := packet.NewNOOP()

	go func() {
		for {
			buffer.Add(packet)
		}
	}()

	for i := 0; i < b.N; i++ {
		buffer.Flush()
	}
}
