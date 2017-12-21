package common

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

const (
	DefaultExpriy = 20000000
	DefaultTries  = 4
	DefaultDelay  = 8000000
)

type Node interface {
	//input:key,value,expriy
	//ouput:result
	Set(string, string, time.Duration) bool

	//input:key,value
	//key不存在的情况下也返回true
	Del(string, string) bool
}

type Mutex struct {
	name   string
	expriy time.Duration

	tries int
	delay time.Duration

	value string

	nodes []Node
	nodem sync.Mutex
}

func NewMutex(name string, nodes []Node) *Mutex {
	return &Mutex{
		name:   name,
		nodes:  nodes,
		tries:  DefaultTries,
		delay:  DefaultDelay,
		expriy: DefaultExpriy,
	}
}

func (m *Mutex) Lock() bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return false
	}
	m.value = base64.StdEncoding.EncodeToString(b)

	for i := 0; i < m.tries; i++ {
		ncnt := 0
		lockcnt := 0
		for _, node := range m.nodes {
			if node == nil {
				continue
			}

			ncnt++

			if node.Set(m.name, m.value, m.expriy) {
				lockcnt++
			}
		}

		fmt.Println("lockcont:", lockcnt, ",ncnt:", ncnt)

		if lockcnt >= ncnt/2+1 {
			return true
		} else {
			for _, node := range m.nodes {
				if node == nil {
					continue
				}
				node.Del(m.name, m.value)
			}
		}

		if i == m.tries-1 {
			continue
		}
		fmt.Println("wait", m.delay)
		time.Sleep(m.delay)
	}
	return false
}

func (m *Mutex) Unlock() bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	for _, node := range m.nodes {
		if node == nil {
			continue
		} else {
			node.Del(m.name, m.value)
		}
	}

	m.value = ""

	return true
}
