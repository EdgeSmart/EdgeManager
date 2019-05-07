package library

import (
	"bytes"
	"io"
	"math/rand"
	"time"
)

const (
	randomString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bufLen       = 1024
)

// GetRandomString random string
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

// ReadAll read all data from io.Reader
func ReadAll(c io.Reader) ([]byte, error) {
	var (
		buffer      bytes.Buffer
		err         error
		requestData []byte
	)
	for {
		buf := make([]byte, bufLen)
		n, err := c.Read(buf)
		if err != nil {
			return nil, err
		}
		buffer.Write(buf[:n])
		if n < bufLen {
			break
		}
	}
	requestData = buffer.Bytes()
	return requestData, err
}
