package main

import "net/url"

var Main UrlDecoder

type UrlDecoder struct {
}

func (ud *UrlDecoder) Process(sessionPayload [][]byte) ([][]byte, error) {
	result := make([][]byte, 0, len(sessionPayload))

	for _, packetPayload := range sessionPayload {
		res, err := url.QueryUnescape(string(packetPayload))
		if err != nil {
			result = append(result, packetPayload)
		} else {
			result = append(result, []byte(res))
		}
	}
	return result, nil
}
