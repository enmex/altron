package main

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

const ChunkedEncodedHeader = "transfer-encoding: chunked"

var Main ChunkDecoder

type ChunkDecoder struct {
	bodyReader *bufio.Reader
}

func (cd *ChunkDecoder) Process(sessionPayload [][]byte) ([][]byte, error) {
	if !cd.isChunkedEncodedSession(sessionPayload) {
		return sessionPayload, nil
	}

	result := make([][]byte, 0, len(sessionPayload))
	for idx, packetPayload := range sessionPayload {
		if idx%2 == 1 {
			endPos := strings.Index(string(packetPayload), "\r\n\r\n") + len("\r\n\r\n")

			body := packetPayload[endPos:]

			decodedData := make([]byte, 0)
			cd.bodyReader = bufio.NewReaderSize(bytes.NewBuffer(body), len(body))
			readBytes := 0

			for {
				chunkSizeBytes, err := cd.readUntil('\r')
				if err != nil {
					return nil, err
				}
				//skip \n
				if _, err := cd.read(1); err != nil {
					return nil, err
				}

				chunkSize, err := strconv.ParseInt(string(chunkSizeBytes), 16, 64)
				if err != nil {
					return nil, err
				}
				chunk, err := cd.read(int(chunkSize))
				if err != nil {
					return nil, err
				}

				//skip \r\n
				if _, err := cd.read(2); err != nil {
					return nil, err
				}

				decodedData = append(decodedData, chunk...)

				readBytes += len(chunkSizeBytes) + 2 + int(chunkSize) + 2 // chuck size + \r\n + chunk + \r\n
				if readBytes >= len(body) {
					break
				}
			}
			packetPayload = append(packetPayload[:endPos], decodedData...)
		}
		result = append(result, packetPayload)
	}

	return result, nil
}

func (cd *ChunkDecoder) isChunkedEncodedSession(sessionPayload [][]byte) bool {
	if len(sessionPayload) <= 1 {
		return false
	}
	return strings.Contains(strings.ToLower(string(sessionPayload[1])), ChunkedEncodedHeader)
}

func (cd *ChunkDecoder) readUntil(delim byte) ([]byte, error) {
	var last byte = 0
	bytesRead := make([]byte, 0)
	for {
		b, err := cd.read(1)
		if err != nil {
			return nil, err
		}
		last = b[0]
		if last == delim {
			break
		}
		bytesRead = append(bytesRead, last)
	}
	return bytesRead, nil
}

func (cd *ChunkDecoder) read(n int) ([]byte, error) {
	data, err := cd.bodyReader.Peek(n)
	if err != nil {
		return nil, err
	}
	cd.bodyReader.Discard(len(data))
	return data, nil
}
