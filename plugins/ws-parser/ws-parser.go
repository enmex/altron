package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

var Main WebsocketParser

const (
	WebSocketConnectionHeader = "connection: upgrade"
	WebSocketUpgradeHeader    = "upgrade: websocket"
	SocketIoPath              = "/socket.io/"
	finalBit                  = 1 << 7
	rsv1Bit                   = 1 << 6
	rsv2Bit                   = 1 << 5
	rsv3Bit                   = 1 << 4
	maskBit                   = 1 << 7
)

type WebsocketParser struct {
	maskingKey      []byte
	reader          *bufio.Reader
	bufferedMessage []byte
}

func (ws *WebsocketParser) Process(sessionPayload [][]byte) ([][]byte, error) {
	if !ws.isWebSocketSession(sessionPayload) {
		return sessionPayload, nil
	}
	serverHandshake := sessionPayload[1]
	endPos := strings.Index(string(serverHandshake), "\r\n\r\n") + len("\r\n\r\n")
	encodedBody := serverHandshake[endPos:]

	parsedPayload := make([]byte, 0)
	if len(encodedBody) > 0 {
		var err error
		parsedPayload, err = ws.parseFrames([]byte(serverHandshake[endPos:]))
		if err != nil {
			return nil, err
		}
	}

	parsedSession := [][]byte{sessionPayload[0], append([]byte(serverHandshake[:endPos]), parsedPayload...)}
	for i := 2; i < len(sessionPayload)-2; i++ {
		parsed, err := ws.parseFrames(sessionPayload[i])
		if err != nil {
			return nil, err
		}
		parsedSession = append(parsedSession, parsed)
	}
	return parsedSession, nil
}

func (ws *WebsocketParser) parseFrames(frames []byte) ([]byte, error) {
	ws.reader = bufio.NewReaderSize(bytes.NewReader(frames), len(frames))
	isFinal := false
	var payload []byte

	for !isFinal {
		data, err := ws.read(2)
		if err != nil {
			return nil, err
		}

		isFinal = data[0]&finalBit != 0
		rsv1 := data[0]&rsv1Bit != 0
		rsv2 := data[0]&rsv2Bit != 0
		rsv3 := data[0]&rsv3Bit != 0
		mask := data[1]&maskBit != 0
		//frameType := int(data[0] & 0xf)

		if rsv1 || rsv2 || rsv3 {
			return nil, fmt.Errorf("RSV flag must be 0")
		}

		payloadLen := int64(data[1] & 0x7f)

		switch payloadLen {
		case 126:
			data, err := ws.read(2)
			if err != nil {
				return nil, err
			}
			payloadLen = int64(binary.BigEndian.Uint16(data))
		case 127:
			data, err := ws.read(8)
			if err != nil {
				return nil, err
			}
			payloadLen = int64(binary.BigEndian.Uint16(data))
		}

		if mask {
			data, err := ws.read(4)
			if err != nil {
				return nil, err
			}
			ws.maskingKey = data
		}
		data, err = ws.read(int(payloadLen))
		if err != nil {
			return nil, err
		}
		payload = make([]byte, 0, len(data))
		if mask {
			for i := 0; i < int(payloadLen); i++ {
				payload = append(payload, data[i]^ws.maskingKey[i%4])
			}
		} else {
			payload = data
		}
		if isFinal && len(ws.bufferedMessage) != 0 {
			ws.bufferedMessage = append(ws.bufferedMessage, payload...)
			payload = ws.bufferedMessage
		} else if !isFinal {
			if len(ws.bufferedMessage) == 0 {
				ws.bufferedMessage = make([]byte, 0)
			}
			ws.bufferedMessage = append(ws.bufferedMessage, payload...)
		}
	}

	return payload, nil
}

func (ws *WebsocketParser) read(n int) ([]byte, error) {
	data, err := ws.reader.Peek(n)
	if err != nil {
		return nil, err
	}
	ws.reader.Discard(len(data))
	return data, nil
}

func (ws *WebsocketParser) isWebSocketSession(sessionPayload [][]byte) bool {
	firstRequest := string(sessionPayload[0])
	if strings.Contains(strings.ToLower(firstRequest), SocketIoPath) && strings.Contains(strings.ToLower(firstRequest), "transport=websocket") {
		return true
	}
	return strings.Contains(strings.ToLower(firstRequest), WebSocketConnectionHeader) && strings.Contains(strings.ToLower(firstRequest), WebSocketUpgradeHeader)
}
