package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
)

var Main Ungzip

type Ungzip struct {
}

func (u *Ungzip) Process(sessionPayload [][]byte) ([][]byte, error) {
	if !u.isGzipSession(sessionPayload) {
		return sessionPayload, nil
	}

	result := make([][]byte, 0, len(sessionPayload))
	for idx, packetPayload := range sessionPayload {
		if idx%2 == 1 {
			packetPayload = u.decode(packetPayload)
		}
		result = append(result, packetPayload)
	}

	return result, nil
}

func (u *Ungzip) decode(packetPayload []byte) []byte {
	endPos := strings.Index(string(packetPayload), "\r\n\r\n") + len("\r\n\r\n")
	body := packetPayload[endPos:]
	reader, err := gzip.NewReader(bufio.NewReaderSize(bytes.NewBuffer(body), len(body)))
	if err != nil {
		return packetPayload
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return packetPayload
	}
	return append(packetPayload[:endPos], decompressedData...)
}

func (u *Ungzip) isGzipSession(sessionPayload [][]byte) bool {
	if len(sessionPayload) <= 1 {
		return false
	}
	return strings.Contains(string(sessionPayload[1]), "Content-Encoding: gzip")
}