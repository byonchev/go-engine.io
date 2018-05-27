package codec

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/byonchev/go-engine.io/internal/packet"
)

// JSONP is a codec for encoding messages for cross-domain polling
type JSONP struct {
	Index    string
	delegate XHR
}

const hexCharacters = "0123456789abcdef"

// Encode encodes payload of packets for single poll
func (codec JSONP) Encode(payload packet.Payload, writer io.Writer) error {
	var buffer bytes.Buffer

	codec.delegate.ForceBase64 = true
	codec.delegate.Encode(payload, &buffer)

	bytes := []byte("___eio[" + codec.Index + "](\"")
	bytes = append(bytes, codec.escape(buffer.String())...)
	bytes = append(bytes, []byte("\");")...)

	_, err := writer.Write(bytes)

	return err
}

// Decode decodes payload of packets
func (codec JSONP) Decode(reader io.Reader) (packet.Payload, error) {
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	if len(data) < 2 {
		return nil, errors.New("invalid form data")
	}

	query, err := url.QueryUnescape(string(data[2:]))

	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBufferString(codec.unescape(query))

	return codec.delegate.Decode(buffer)
}

func (codec JSONP) escape(data string) []byte {
	var buffer bytes.Buffer

	for _, char := range data {
		switch char {
		case '\\', '"', '/':
			buffer.WriteRune('\\')
			buffer.WriteRune(char)
		case '\b':
			buffer.WriteRune('\\')
			buffer.WriteRune('b')
		case '\f':
			buffer.WriteRune('\\')
			buffer.WriteRune('f')
		case '\n':
			buffer.WriteRune('\\')
			buffer.WriteRune('n')
		case '\r':
			buffer.WriteRune('\\')
			buffer.WriteRune('r')
		case '\t':
			buffer.WriteRune('\\')
			buffer.WriteRune('t')
		case '\u2028', '\u2029':
			buffer.WriteString("\\u202")
			buffer.WriteByte(hexCharacters[char&0xF])
		default:
			if char < 0x20 {
				buffer.WriteString("\\u00")
				buffer.WriteByte(hexCharacters[char>>4])
				buffer.WriteByte(hexCharacters[char&0xF])
			} else {
				buffer.WriteRune(char)
			}
		}
	}

	return buffer.Bytes()
}

func (codec JSONP) unescape(data string) string {
	data = strings.Replace(data, `\\\\n`, "\\n", -1)
	data = strings.Replace(data, `\\n`, "\n", -1)

	return data
}
