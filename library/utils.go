package library

import (
	"math/rand"
	"time"
)

const (
	randomString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GetRandomString(l int) string {
	bytes := []byte(randomString)
	candidate := len(bytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(candidate)])
	}
	return string(result)
}
