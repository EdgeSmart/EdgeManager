package token

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type tokenInterface interface {
	init(config Config) error
	Generate(flag string) string
	Set(flag string, token string, data string) error
	Get(token string) (string, error)
	GetToken(flag string) ([]string, error)
	Active(token string) error
	Destroy(token string) error
	GC() error
}

type Config struct {
	Flag   string
	Expire uint
}

type instance struct {
}

var tokenPool map[string]tokenInterface = map[string]tokenInterface{}

func GetInstance(flag string) (tokenInterface, error) {
	if ins, exists := tokenPool[flag]; exists {
		fmt.Println("get")
		fmt.Println(ins)
		return ins, nil
	}
	return nil, errors.New("")
}

func NewInstance(flag string, driver string, config Config) (tokenInterface, error) {
	if flag == "" {
		return nil, errors.New("")
	}
	if _, err := GetInstance(flag); err == nil {
		return nil, errors.New("")
	}
	switch driver {
	case "memery":
		ins := &memery{}
		err := ins.init(config)
		fmt.Println("new")
		fmt.Println(ins)
		if err != nil {
			return nil, errors.New("")
		}
		tokenPool[flag] = ins
		return ins, nil
	}
	return nil, errors.New("")
}

func randomGenerate(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
