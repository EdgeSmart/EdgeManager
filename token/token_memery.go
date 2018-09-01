package token

import (
	"errors"
	"strings"
	"time"
)

const expire = 600

type memeryitem struct {
	flag string
	data string
	ts   int64
}

type memery struct {
	tokenPool map[string]memeryitem
	flagPool  map[string][]string
}

// var

func (m *memery) init(config Config) error {
	m.tokenPool = map[string]memeryitem{}
	m.flagPool = map[string][]string{}
	return nil
}

func (m *memery) Generate(flag string) string {
	return randomGenerate(64)
}

func (m *memery) Set(flag string, token string, data string) error {
	ins, exists := m.tokenPool[token]
	if exists {
		if ins.flag != flag {
			return errors.New("error")
		}
		ins.data = data
	} else {
		item := memeryitem{
			data: data,
			ts:   time.Now().Unix(),
		}
		if strings.Trim(flag, " ") != "" {
			item.flag = flag
			if _, exists := m.flagPool[flag]; exists {
				m.flagPool[flag] = append(m.flagPool[flag], token)
			} else {
				flags := []string{flag}
				m.flagPool[flag] = flags
			}
		}
		m.tokenPool[token] = item
	}
	return nil
}

func (m *memery) Get(token string) (string, error) {
	ins, exists := m.tokenPool[token]
	if exists {
		return ins.data, nil
	} else if ins.ts+expire < time.Now().Unix() {
		m.Destroy(token)
		return "", errors.New("error")
	}
	return "", errors.New("error")
}

func (m *memery) GetToken(flag string) ([]string, error) {
	tok, exists := m.flagPool[flag]
	if exists {
		return tok, nil
	}
	return []string{}, nil
}

func (m *memery) Active(token string) error {
	ins, exists := m.tokenPool[token]
	if exists {
		ins.ts = time.Now().Unix()
		return nil
	}
	return errors.New("error")
}

func (m *memery) Destroy(token string) error {
	_, exists := m.tokenPool[token]
	if exists {
		delete(m.tokenPool, token)
	}
	return nil
}

func (m *memery) GC() error {
	return nil
}
