package Parse

import (
	"bytes"
	"io"
	"net"
	"time"
)

/*
0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
*/

// Frame represents a WebSocket frame.
type Frame struct {
	Fin           bool   // 1 bit - 1 if this is the final fragment in the message (it could be the only fragment), otherwise 0
	Rsv1          bool   // 1 bit - reserved
	Rsv2          bool   // 1 bit - reserved
	Rsv3          bool   // 1 bit - reserved
	Op            opcode // 4 bits - defines the interpretation of the "Payload data". If an unknown opcode is received, the receiving endpoint MUST _Fail the WebSocket Connection_. The following values are defined.
	Mask          bool   // 1 bit - defines whether the "Payload data" is masked. If set to 1, a masking key is present in masking-key, and this is used to unmask the "Payload data".
	PayloadLength int    // 7 bits, 7+16 bits, or 7+64 bits
	Payload       []byte // x*8 bits of "Payload data"
}

type opcode byte

const FrameMaskSize = 4

const (
	ContinuationFrame opcode = 0x0
	TextFrame         opcode = 0x1
	BinaryFrame       opcode = 0x2
	CloseFrame        opcode = 0x8
	PingFrame         opcode = 0x9
	PongFrame         opcode = 0xA
)

type FrameErrorCode int
type FrameError struct {
	Code    FrameErrorCode
	Message string
}

const (
	ErrBadMask = iota
	ErrUnsupportedFrame
	ErrBuffer
	ErrSend
)

func ReadFromWebsocket(socket net.Conn) (*Frame, *FrameError) {
	frame := &Frame{}

	buffer := make([]byte, 2)
	_, err := io.ReadFull(socket, buffer)

	if err != nil {
		return nil, &FrameError{Code: ErrBuffer, Message: err.Error()}
	}

	/* check this is the final fragment in the message is make and operation with 0x80 to check if the first bit is 1
	example: 0x80 (10000000) & 0x01(00000001) = 00000000 */
	frame.Fin = buffer[0]&0x80 == 0x80

	if !frame.Fin {
		return frame, &FrameError{Code: ErrUnsupportedFrame, Message: "unsupported frame"}
	}

	/* check the opcode is make and operation with 0x0F to get the last 4 bits
	example: 0x0F (00001111) & 0x01(00000001) = 00000001 */
	frame.Op = opcode(buffer[0] & 0x0F)

	/* check if the payload is masked is make and operation with 0x80 to check if the first bit is 1
	example: 0x80(10000000) & 0x01(00000001) = 00000000 */
	frame.Mask = buffer[1]&0x80 == 0x80

	if !frame.Mask {
		return frame, &FrameError{Code: ErrBadMask, Message: "unmasked frame"}
	}

	if frame.Op == CloseFrame {
		return frame, nil
	} else if frame.Op == PingFrame {
		return frame, nil
	} else if frame.Op != TextFrame {
		return frame, &FrameError{Code: ErrUnsupportedFrame, Message: "unsupported frame"}
	}

	/* check the payload length is make and operation with 0x7F to get the last 7 bits
	example: 0x7F (01111111) & 0x01(00000001) = 00000001 */
	frame.PayloadLength = int(buffer[1] & 0x7F)

	if frame.PayloadLength == 126 || frame.PayloadLength == 127 {
		var additionalBytes = 0

		if frame.PayloadLength == 126 {
			additionalBytes = 2
		} else {
			additionalBytes = 8
		}

		buffer = make([]byte, additionalBytes)
		_, err := io.ReadFull(socket, buffer)

		if err != nil {
			return nil, &FrameError{Code: ErrBuffer, Message: err.Error()}
		}

		if frame.PayloadLength == 126 {
			frame.PayloadLength = int(buffer[0])<<8 | int(buffer[1])
		} else {
			for i := 0; i < 8; i++ {
				frame.PayloadLength |= int(buffer[i]) << (8 * (7 - i))
			}
		}
	}

	frame.Payload = make([]byte, frame.PayloadLength+FrameMaskSize)
	_, err = io.ReadFull(socket, frame.Payload)

	if err != nil {
		return nil, &FrameError{Code: ErrBuffer, Message: err.Error()}
	}

	/* unmask the payload */
	for i := 0; i < frame.PayloadLength; i++ {
		frame.Payload[i+FrameMaskSize] ^= frame.Payload[i%FrameMaskSize]
	}

	frame.Payload = frame.Payload[FrameMaskSize:]
	return frame, nil
}

func WriteToWebsocket(socket net.Conn, data string, c opcode) *FrameError {
	var buffer bytes.Buffer
	header := byte(0x80 | byte(c))
	buffer.Write([]byte{header})

	if len(data) < 126 {
		buffer.Write([]byte{byte(len(data))})
	} else if len(data) < 65536 {
		buffer.Write([]byte{126})
		buffer.Write([]byte{byte(len(data) >> 8)})
		buffer.Write([]byte{byte(len(data) & 0xFF)})
	} else {
		buffer.Write([]byte{127})

		for i := 0; i < 8; i++ {
			buffer.Write([]byte{byte(len(data) >> (8 * (7 - i)) & 0xFF)})
		}
	}

	buffer.Write([]byte(data))

	err := socket.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		return &FrameError{Code: ErrSend, Message: err.Error()}
	}

	_, err = socket.Write(buffer.Bytes())

	if err != nil {
		return &FrameError{Code: ErrSend, Message: err.Error()}
	}

	return nil
}
