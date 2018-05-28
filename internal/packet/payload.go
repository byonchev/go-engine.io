package packet

// Payload is a collection of packets
type Payload []Packet

// ContainsBinary returns true if payload contains at least one binary packet
func (payload Payload) ContainsBinary() bool {
	for _, packet := range payload {
		if packet.Binary {
			return true
		}
	}

	return false
}
