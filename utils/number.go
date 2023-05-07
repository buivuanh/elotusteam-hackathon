package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomNumberString generates a random number string of the specified length
func GenerateRandomNumberString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
