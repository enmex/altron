package utils

import "math/rand"

var runes []rune = []rune("abcdefghijklmnopqrstuvwxyz")

func RandomString(size int) string {
	res := make([]rune, 0)
	for i := 0; i < size; i++ {
		res = append(res, runes[rand.Intn(len(runes)-1)])
	}
	return string(res)
}
