package random

import "math/rand/v2"

const allowedChars = "aAbBcCdDeEfFgGhHjJkKmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"

func NewRandomString(length int) string {
	chars := []rune(allowedChars)
	res := make([]rune, length)

	for i := range res {
		res[i] = chars[rand.IntN(len(chars))]
	}

	return string(res)
}
