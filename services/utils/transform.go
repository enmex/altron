package utils

import (
	"fmt"
	"math"
	"sync"
)

func TransformToAsciiBytes(enc []byte) string {
	wg := sync.WaitGroup{}
	lineLength := int(math.Min(float64(len(enc)), 10000))
	routinesNumber := len(enc) / lineLength

	wg.Add(routinesNumber + 1)
	parts := make(map[int]string)

	for i := 0; i <= routinesNumber; i++ {
		go func(idx int) {
			defer wg.Done()
			to := (idx + 1) * lineLength
			if to > len(enc) {
				to = len(enc)
			}
			payload := enc[idx*lineLength : to]

			resultPayload := ""
			for _, b := range payload {
				if b == '\\' {
					resultPayload += "\\\\"
				} else if 32 <= b && b <= 126 {
					resultPayload += string(b)
				} else {
					resultPayload += fmt.Sprintf("\\x%02x", b)
				}
			}
			parts[idx] = resultPayload
		}(i)
	}

	wg.Wait()

	transformed := ""
	for i := 0; i <= routinesNumber; i++ {
		transformed += parts[i]
	}

	return transformed
}
